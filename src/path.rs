//! Pure, I/O-free path manipulation.
//!
//! Every function takes the home directory (and, where relevant, the working
//! directory) as explicit arguments rather than reading the environment, so the
//! logic is deterministic and trivially testable.

/// Expand a leading `~` to `home`.
///
/// A path that does not begin with `~` is returned unchanged. `~` alone expands
/// to `home`; `~/rest` expands to `home` followed by `/rest`.
pub fn expand_tilde(path: &str, home: &str) -> String {
    if !path.starts_with('~') {
        return path.to_string();
    }
    format!("{}{}", home, &path[1..])
}

/// Normalise a user-supplied path into the canonical form stored in the config.
///
/// Precedence of the input shape:
/// - Absolute paths under `home` are shortened to `~/…`.
/// - Absolute paths outside `home` are kept as-is.
/// - Tilde paths are kept as-is.
/// - Everything else is treated as relative to `pwd`.
///
/// The result is always lexically cleaned (`.`/`..`/redundant separators removed).
pub fn sanitize_file_path(input: &str, home: &str, pwd: &str) -> String {
    let cur_dir = match pwd.strip_prefix(home) {
        Some(rest) => format!("~{rest}"),
        None => pwd.to_string(),
    };

    let path = if input.starts_with('/') {
        match input.strip_prefix(home) {
            Some(rest) => format!("~{rest}"),
            None => input.to_string(),
        }
    } else if input.starts_with('~') {
        input.to_string()
    } else {
        format!("{}/{}", cur_dir, input)
    };

    clean(&path)
}

/// Lexically clean a path, mirroring Go's `filepath.Clean`.
///
/// Collapses redundant separators, drops `.` elements, and resolves inner `..`
/// elements without touching the filesystem. A leading `~` is treated as an
/// ordinary path element, not a root.
pub fn clean(path: &str) -> String {
    if path.is_empty() {
        return ".".to_string();
    }

    let bytes = path.as_bytes();
    let rooted = bytes[0] == b'/';
    let n = bytes.len();
    let mut buf = vec![0u8; n];
    let mut w = 0usize;
    let mut r = 0usize;
    let mut dotdot = 0usize;

    if rooted {
        buf[w] = b'/';
        w += 1;
        r = 1;
        dotdot = 1;
    }

    while r < n {
        if bytes[r] == b'/' || (bytes[r] == b'.' && (r + 1 == n || bytes[r + 1] == b'/')) {
            r += 1;
        } else if bytes[r] == b'.'
            && r + 1 < n
            && bytes[r + 1] == b'.'
            && (r + 2 == n || bytes[r + 2] == b'/')
        {
            r += 2;
            if w > dotdot {
                w -= 1;
                while w > dotdot && buf[w] != b'/' {
                    w -= 1;
                }
            } else if !rooted {
                if w > 0 {
                    buf[w] = b'/';
                    w += 1;
                }
                buf[w] = b'.';
                w += 1;
                buf[w] = b'.';
                w += 1;
                dotdot = w;
            }
        } else {
            if (rooted && w != 1) || (!rooted && w != 0) {
                buf[w] = b'/';
                w += 1;
            }
            while r < n && bytes[r] != b'/' {
                buf[w] = bytes[r];
                w += 1;
                r += 1;
            }
        }
    }

    if w == 0 {
        return ".".to_string();
    }
    String::from_utf8_lossy(&buf[..w]).into_owned()
}

#[cfg(test)]
mod tests {
    use super::*;

    const HOME: &str = "/Users/kevin";
    const PWD: &str = "/Users/kevin/fake/dir";

    #[test]
    fn expand_tilde_expands_to_absolute() {
        assert_eq!(
            expand_tilde("~/mydir/file.txt", HOME),
            "/Users/kevin/mydir/file.txt"
        );
    }

    #[test]
    fn expand_tilde_leaves_absolute_alone() {
        assert_eq!(expand_tilde("/etc/passwd", HOME), "/etc/passwd");
    }

    #[test]
    fn expand_tilde_bare_expands_to_home() {
        assert_eq!(expand_tilde("~", HOME), HOME);
    }

    #[test]
    fn sanitize_leaves_abspath_outside_home_alone() {
        assert_eq!(sanitize_file_path("/etc/passwd", HOME, PWD), "/etc/passwd");
    }

    #[test]
    fn sanitize_cleans_abspath() {
        assert_eq!(
            sanitize_file_path("/relative/path/../path/file.txt", HOME, PWD),
            "/relative/path/file.txt"
        );
    }

    #[test]
    fn sanitize_shrinks_home_path_with_tilde() {
        assert_eq!(
            sanitize_file_path("/Users/kevin/file.txt", HOME, PWD),
            "~/file.txt"
        );
    }

    #[test]
    fn sanitize_leaves_tilde_alone() {
        assert_eq!(
            sanitize_file_path("~/mydir/file.txt", HOME, PWD),
            "~/mydir/file.txt"
        );
    }

    #[test]
    fn sanitize_joins_relative_to_pwd() {
        assert_eq!(
            sanitize_file_path("file.txt", HOME, PWD),
            "~/fake/dir/file.txt"
        );
    }

    #[test]
    fn sanitize_relative_with_parent_segments() {
        assert_eq!(
            sanitize_file_path("../file.txt", HOME, PWD),
            "~/fake/file.txt"
        );
    }

    #[test]
    fn clean_collapses_redundant_separators() {
        assert_eq!(clean("/a//b/./c"), "/a/b/c");
    }

    #[test]
    fn clean_empty_is_dot() {
        assert_eq!(clean(""), ".");
    }
}
