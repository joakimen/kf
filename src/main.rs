//! `kf` — manage "known files", the files you visit regularly.
//!
//! This module is the imperative shell: it parses arguments, reads the
//! environment, performs filesystem and process I/O, and drives the interactive
//! prompts. All parsing, normalisation, and decision logic lives in the
//! I/O-free core modules ([`path`], [`config`]).

mod config;
mod path;
mod store;

use std::io::IsTerminal;
use std::path::{Path, PathBuf};
use std::process::Command;

use anyhow::{Context, Result, bail};
use clap::{Parser, Subcommand};
use inquire::{InquireError, MultiSelect, Select};

#[derive(Parser)]
#[command(
    name = "kf",
    version,
    about = "Manage known files — files you visit regularly"
)]
struct Cli {
    /// Path to the config file (overrides $XDG_CONFIG_HOME and the default)
    #[arg(short, long, global = true, value_name = "PATH")]
    config: Option<String>,

    #[command(subcommand)]
    command: Cmd,
}

#[derive(Subcommand)]
enum Cmd {
    /// Add a file to the list of known files
    Add {
        /// File path to add
        file: String,
    },
    /// List known files
    List {
        /// Show file existence indicators
        #[arg(long)]
        status: bool,
        /// Only show files that do not exist locally
        #[arg(long, conflicts_with = "exists")]
        missing: bool,
        /// Only show files that exist locally
        #[arg(long)]
        exists: bool,
    },
    /// Remove a file from the list of known files
    Forget {
        /// File path to remove; omit to choose interactively
        file: Option<String>,
    },
    /// Select a known file with a fuzzy finder and print its path
    Pick,
    /// Remove all known files that no longer exist locally
    Prune,
    /// Open the configuration file in $EDITOR
    Edit,
    /// Print the configuration file path
    Config,
}

fn main() {
    if let Err(err) = run() {
        eprintln!("error: {err:#}");
        std::process::exit(1);
    }
}

fn run() -> Result<()> {
    let cli = Cli::parse();
    let config_path = resolve_config_path(cli.config.as_deref())?;

    match cli.command {
        Cmd::Add { file } => cmd_add(&file, &config_path),
        Cmd::List {
            status,
            missing,
            exists,
        } => cmd_list(status, missing, exists, &config_path),
        Cmd::Forget { file } => cmd_forget(file.as_deref(), &config_path),
        Cmd::Pick => cmd_pick(&config_path),
        Cmd::Prune => cmd_prune(&config_path),
        Cmd::Edit => cmd_edit(&config_path),
        Cmd::Config => {
            println!("{}", config_path.display());
            Ok(())
        }
    }
}

fn getenv(key: &str) -> String {
    std::env::var(key).unwrap_or_default()
}

fn require_env(key: &str) -> Result<String> {
    let value = getenv(key);
    if value.is_empty() {
        bail!("${key} is not set");
    }
    Ok(value)
}

fn resolve_config_path(flag: Option<&str>) -> Result<PathBuf> {
    let xdg = getenv("XDG_CONFIG_HOME");
    let home = getenv("HOME");
    if flag.is_none() && xdg.is_empty() && home.is_empty() {
        bail!("cannot determine config path: set $HOME, $XDG_CONFIG_HOME, or pass --config");
    }
    Ok(config::resolve_path(flag, &xdg, &home))
}

fn cmd_add(file: &str, config_path: &Path) -> Result<()> {
    let home = require_env("HOME")?;
    let pwd = require_env("PWD")?;

    let sanitized = path::sanitize_file_path(file, &home, &pwd);
    let expanded = path::expand_tilde(&sanitized, &home);
    if !Path::new(&expanded).exists() {
        eprintln!("warning: {expanded} does not exist");
    }

    let mut lines = store::read_lines(config_path)?;
    if lines.iter().any(|line| line == &sanitized) {
        bail!("entry already exists in configuration file");
    }

    lines.push(sanitized.clone());
    store::write_lines(config_path, &lines)?;
    println!("Added {sanitized}");
    Ok(())
}

fn cmd_list(status: bool, missing: bool, exists: bool, config_path: &Path) -> Result<()> {
    let home = require_env("HOME")?;
    let lines = store::read_lines(config_path)?;

    if !status && !missing && !exists {
        for line in &lines {
            println!("{}", path::expand_tilde(line, &home));
        }
        return Ok(());
    }

    let use_color = status && std::io::stdout().is_terminal() && !no_color();

    for line in &lines {
        let expanded = path::expand_tilde(line, &home);
        let present = Path::new(&expanded).exists();

        if missing && present {
            continue;
        }
        if exists && !present {
            continue;
        }

        if !status {
            println!("{expanded}");
            continue;
        }

        match (use_color, present) {
            (true, true) => println!("\x1b[32m✓ {expanded}\x1b[0m"),
            (true, false) => println!("\x1b[31m✗ {expanded}\x1b[0m"),
            (false, true) => println!("✓ {expanded}"),
            (false, false) => println!("✗ {expanded}"),
        }
    }
    Ok(())
}

fn cmd_forget(file: Option<&str>, config_path: &Path) -> Result<()> {
    match file {
        Some(file) => forget_by_arg(file, config_path),
        None => forget_interactive(config_path),
    }
}

fn forget_by_arg(file: &str, config_path: &Path) -> Result<()> {
    let home = require_env("HOME")?;
    let pwd = require_env("PWD")?;
    let sanitized = path::sanitize_file_path(file, &home, &pwd);

    let lines = store::read_lines(config_path)?;
    let (kept, removed) = config::partition_remove(&lines, std::slice::from_ref(&sanitized));
    if removed.is_empty() {
        println!("No matching entry found");
        return Ok(());
    }

    store::write_lines(config_path, &kept)?;
    println!("Removed {sanitized}");
    Ok(())
}

fn forget_interactive(config_path: &Path) -> Result<()> {
    let lines = store::read_lines(config_path)?;
    if lines.is_empty() {
        println!("No known files");
        return Ok(());
    }

    let selected = match MultiSelect::new("Select files to forget", lines.clone()).prompt() {
        Ok(selected) => selected,
        Err(InquireError::OperationCanceled | InquireError::OperationInterrupted) => return Ok(()),
        Err(e) => return Err(e).context("fuzzy finder failed"),
    };

    let (kept, removed) = config::partition_remove(&lines, &selected);
    if removed.is_empty() {
        return Ok(());
    }

    store::write_lines(config_path, &kept)?;
    for file in &removed {
        println!("Removed {file}");
    }
    Ok(())
}

fn cmd_pick(config_path: &Path) -> Result<()> {
    let home = require_env("HOME")?;
    let files: Vec<String> = store::read_lines(config_path)?
        .iter()
        .map(|line| path::expand_tilde(line, &home))
        .collect();
    if files.is_empty() {
        bail!("no known files");
    }

    match Select::new("Pick a file", files).prompt() {
        Ok(choice) => {
            println!("{choice}");
            Ok(())
        }
        Err(InquireError::OperationCanceled | InquireError::OperationInterrupted) => Ok(()),
        Err(e) => Err(e).context("fuzzy finder failed"),
    }
}

fn cmd_prune(config_path: &Path) -> Result<()> {
    let home = require_env("HOME")?;
    let lines = store::read_lines(config_path)?;

    let stale: Vec<String> = lines
        .iter()
        .filter(|line| !Path::new(&path::expand_tilde(line, &home)).exists())
        .cloned()
        .collect();

    if stale.is_empty() {
        println!("No stale entries");
        return Ok(());
    }

    let (kept, removed) = config::partition_remove(&lines, &stale);
    store::write_lines(config_path, &kept)?;
    for file in &removed {
        println!("Pruned {file}");
    }
    Ok(())
}

fn cmd_edit(config_path: &Path) -> Result<()> {
    let editor = require_env("EDITOR")?;
    let status = Command::new(&editor)
        .arg(config_path)
        .status()
        .with_context(|| format!("running editor {editor}"))?;
    if !status.success() {
        bail!("editor exited unsuccessfully");
    }
    Ok(())
}

/// Honour the `NO_COLOR` convention: colour is disabled when the variable is
/// present and non-empty.
fn no_color() -> bool {
    std::env::var_os("NO_COLOR").is_some_and(|v| !v.is_empty())
}
