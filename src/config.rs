//! Pure helpers for the config file's contents and location.
//!
//! No filesystem access happens here: functions take the raw lines (or the
//! relevant environment values) and return transformed data, so precedence and
//! normalisation rules are tested without touching disk.

use std::collections::HashSet;
use std::path::PathBuf;

/// Normalise config lines for persistence: trim each line, drop blanks, remove
/// duplicates, and sort so the written file is deterministic.
pub fn sanitize(lines: &[String]) -> Vec<String> {
    let mut seen = HashSet::new();
    let mut result: Vec<String> = lines
        .iter()
        .map(|l| l.trim().to_string())
        .filter(|l| !l.is_empty())
        .filter(|l| seen.insert(l.clone()))
        .collect();
    result.sort();
    result
}

/// Split `lines` into the entries to keep and the entries removed, preserving
/// the original order of the kept entries. Matching is exact against the raw
/// stored line.
pub fn partition_remove(lines: &[String], to_remove: &[String]) -> (Vec<String>, Vec<String>) {
    let remove: HashSet<&String> = to_remove.iter().collect();
    let mut kept = Vec::new();
    let mut removed = Vec::new();
    for line in lines {
        if remove.contains(line) {
            removed.push(line.clone());
        } else {
            kept.push(line.clone());
        }
    }
    (kept, removed)
}

/// Resolve the config file path.
///
/// Precedence: explicit `flag` value > `XDG_CONFIG_HOME/kf/config` > `home/.config/kf/config`.
pub fn resolve_path(flag: Option<&str>, xdg_config_home: &str, home: &str) -> PathBuf {
    if let Some(flag) = flag {
        return PathBuf::from(flag);
    }
    let base = if xdg_config_home.is_empty() {
        PathBuf::from(home).join(".config")
    } else {
        PathBuf::from(xdg_config_home)
    };
    base.join("kf").join("config")
}

#[cfg(test)]
mod tests {
    use super::*;

    fn owned(items: &[&str]) -> Vec<String> {
        items.iter().map(|s| s.to_string()).collect()
    }

    #[test]
    fn sanitize_sorts_alphabetically() {
        assert_eq!(
            sanitize(&owned(&["~/z/file.txt", "~/a/file.txt", "~/m/file.txt"])),
            owned(&["~/a/file.txt", "~/m/file.txt", "~/z/file.txt"])
        );
    }

    #[test]
    fn sanitize_removes_duplicates() {
        assert_eq!(
            sanitize(&owned(&["~/b.txt", "~/a.txt", "~/b.txt"])),
            owned(&["~/a.txt", "~/b.txt"])
        );
    }

    #[test]
    fn sanitize_trims_whitespace() {
        assert_eq!(
            sanitize(&owned(&[" ~/b.txt ", " ~/a.txt "])),
            owned(&["~/a.txt", "~/b.txt"])
        );
    }

    #[test]
    fn sanitize_drops_blank_lines() {
        assert_eq!(
            sanitize(&owned(&["~/a.txt", "   ", "", "~/b.txt"])),
            owned(&["~/a.txt", "~/b.txt"])
        );
    }

    #[test]
    fn partition_splits_kept_and_removed() {
        let lines = owned(&["a", "b", "c", "d"]);
        let (kept, removed) = partition_remove(&lines, &owned(&["b", "d"]));
        assert_eq!(kept, owned(&["a", "c"]));
        assert_eq!(removed, owned(&["b", "d"]));
    }

    #[test]
    fn partition_ignores_unmatched_targets() {
        let lines = owned(&["a", "b"]);
        let (kept, removed) = partition_remove(&lines, &owned(&["x"]));
        assert_eq!(kept, owned(&["a", "b"]));
        assert!(removed.is_empty());
    }

    #[test]
    fn resolve_path_prefers_flag() {
        assert_eq!(
            resolve_path(Some("/custom/kfrc"), "/xdg", "/home/kevin"),
            PathBuf::from("/custom/kfrc")
        );
    }

    #[test]
    fn resolve_path_uses_xdg_when_set() {
        assert_eq!(
            resolve_path(None, "/home/kevin/.xdg", "/home/kevin"),
            PathBuf::from("/home/kevin/.xdg/kf/config")
        );
    }

    #[test]
    fn resolve_path_falls_back_to_home() {
        assert_eq!(
            resolve_path(None, "", "/home/kevin"),
            PathBuf::from("/home/kevin/.config/kf/config")
        );
    }
}
