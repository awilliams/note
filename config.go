package main

import (
	"bytes"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/charmbracelet/glamour"
)

//go:embed README.md
var readme string

type config struct {
	weekOffset int
	baseDir    string
	editor     string
	version    bool
}

var errFullHelp = errors.New("fullHelp")

func (c *config) parse(args []string) error {
	if len(args) == 0 {
		return errors.New("arguments should include at least 1 value (exec name)")
	}
	exeName := filepath.Base(args[0])
	args = args[1:]

	c.baseDir = func() string {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		return filepath.Join(home, exeName)
	}()
	c.editor = os.Getenv("EDITOR")

	var help bool
	fs := flag.NewFlagSet("note", flag.ContinueOnError)
	fs.StringVar(&c.baseDir, "d", c.baseDir, "Root directory of note files")
	fs.StringVar(&c.editor, "e", c.editor, "Editor executable ($EDITOR)")
	fs.BoolVar(&c.version, "v", c.version, "Print version information")
	fs.BoolVar(&help, "help", help, "Print README in addition to standard help (-h) information")

	// Copy default usage string to buffer, then redirect
	// output to /dev/null. This is to ensure usage is output once,
	// since fs.Parse may or may not call PrintDefaults internally.
	var defaults bytes.Buffer
	fs.SetOutput(&defaults)
	fs.PrintDefaults()
	fs.SetOutput(io.Discard)

	err := func() error {
		if len(args) > 0 && len(args[0]) > 0 {
			farg := args[0]
			re := regexp.MustCompile(`^([-\+])?(\d+$)`)
			matches := re.FindStringSubmatch(farg)
			if len(matches) == 3 {
				var m int
				switch matches[1] {
				case "-":
					m = -1
				case "+":
					m = 1
				default:
					// TODO: Allow for absoulte weeks, e.g '23' = week 23,
					// in addition to relative offsets.
					return fmt.Errorf("unable to parse offset value %q", farg)
				}

				v, err := strconv.Atoi(matches[2])
				if err != nil {
					return fmt.Errorf("unable to parse offset value %q: %v", farg, err)
				}
				c.weekOffset = (m * v)
				args = args[1:]
			}
		}

		if err := fs.Parse(args); err != nil {
			return err
		}
		if help {
			return errFullHelp
		}
		if c.baseDir == "" {
			return errors.New("noteDir cannot be blank")
		}
		if c.editor == "" {
			return errors.New("editor cannot be blank")
		}
		return nil
	}()

	if err == nil {
		return nil
	}

	var eb errBuf
	appendUsage := func() {
		fmt.Fprintf(&eb, "Usage: %s [OFFSET] [OPTIONS]\n", exeName)
		fmt.Fprintln(&eb, `OFFSET:
  open note file OFFSET number of weeks relative to now. Example: -2 (two weeks ago); +1 (next week)

OPTIONS:`)
		defaults.WriteTo(&eb)
	}

	switch err {
	case flag.ErrHelp:
		appendUsage()
	case errFullHelp:
		appendUsage()

		r, _ := glamour.NewTermRenderer(
			// ASCII output (no escape sequences).
			glamour.WithStyles(glamour.ASCIIStyleConfig),
			// Wrap output to 80 chars.
			glamour.WithWordWrap(80),
		)
		md, _ := r.Render(readme)
		fmt.Fprintln(&eb, "\nREADME")
		fmt.Fprint(&eb, md)
	default:
		fmt.Fprintf(&eb, "Error: %s\n\n", err.Error())
		appendUsage()
	}

	return &eb
}

type errBuf struct {
	bytes.Buffer
}

func (e *errBuf) Error() string {
	return e.String()
}
