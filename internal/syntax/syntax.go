package syntax

import (
	"bytes"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

const (
	style     = "dracula"
	formatter = "terminal16m"
)

// Colorize reads a file from the filesystem and returns its contents, colorized
func Colorize(filename string, contents string) (string, error) {
	lexer := determineLexer(filename, contents)
	formatter := formatters.Get(formatter)
	style := styles.Get(style)

	var nilOptions *chroma.TokeniseOptions
	iterator, err := lexer.Tokenise(nilOptions, contents)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// determineLexer determines the lexer to use for a given file based on its filename or, alternatively, file contents.
// If type could not be determined, it falls back to the generic lexer.
func determineLexer(filename string, contents string) chroma.Lexer {
	filenameLexer := lexers.Match(filename)
	if filenameLexer != nil {
		return filenameLexer
	}

	filetypeLexer := lexers.Analyse(contents)
	if filetypeLexer != nil {
		return filetypeLexer
	}

	return lexers.Fallback
}
