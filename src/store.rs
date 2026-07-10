//! Filesystem access for the config file (the imperative shell).
//!
//! Reads and writes the line-oriented config. Writes go through a temp file and
//! an atomic rename so a crash mid-write never leaves a truncated config.

use std::fs;
use std::io::Write;
use std::os::unix::fs::DirBuilderExt;
use std::path::Path;
use std::sync::atomic::{AtomicU32, Ordering};

use anyhow::{Context, Result};

use crate::config;

/// Read the config file into trimmed-of-newline lines.
///
/// A missing file is not an error: it yields an empty list, matching the
/// "no known files yet" state.
pub fn read_lines(path: &Path) -> Result<Vec<String>> {
    match fs::read_to_string(path) {
        Ok(contents) => Ok(contents.lines().map(|l| l.to_string()).collect()),
        Err(e) if e.kind() == std::io::ErrorKind::NotFound => Ok(Vec::new()),
        Err(e) => Err(e).with_context(|| format!("reading config file {}", path.display())),
    }
}

/// Write `lines` to the config file after sanitising them (trim, dedupe, sort).
///
/// The parent directory is created with `0700` if absent, and the write is
/// atomic via a temp file in the same directory followed by a rename.
pub fn write_lines(path: &Path, lines: &[String]) -> Result<()> {
    let sanitized = config::sanitize(lines);

    let dir = path
        .parent()
        .filter(|p| !p.as_os_str().is_empty())
        .unwrap_or_else(|| Path::new("."));
    fs::DirBuilder::new()
        .recursive(true)
        .mode(0o700)
        .create(dir)
        .with_context(|| format!("creating config directory {}", dir.display()))?;

    let tmp = dir.join(temp_name());
    let write_result = (|| -> Result<()> {
        let mut file = fs::File::create(&tmp)
            .with_context(|| format!("creating temp file {}", tmp.display()))?;
        for line in &sanitized {
            writeln!(file, "{line}").with_context(|| format!("writing to {}", tmp.display()))?;
        }
        file.sync_all()
            .with_context(|| format!("flushing {}", tmp.display()))?;
        Ok(())
    })();

    if let Err(e) = write_result {
        let _ = fs::remove_file(&tmp);
        return Err(e);
    }

    fs::rename(&tmp, path).with_context(|| {
        let _ = fs::remove_file(&tmp);
        format!("replacing config file {}", path.display())
    })
}

fn temp_name() -> String {
    static COUNTER: AtomicU32 = AtomicU32::new(0);
    let seq = COUNTER.fetch_add(1, Ordering::Relaxed);
    format!(".kf-{}-{}.tmp", std::process::id(), seq)
}
