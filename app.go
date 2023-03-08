package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/glamour"
)

type app struct {
	baseDir    string
	weekOffset int
	names      bool
	print      bool
	printMD    bool
	editor     string
}

func (a *app) run() error {
	curRange := relWorkWeek(a.weekOffset)
	prevRange := relWorkWeek(a.weekOffset - 1)

	n := newNote(a.baseDir, curRange, prevRange)

	switch {
	case a.names:
		fmt.Println(n.path)
		return nil

	case a.print:
		return n.print(os.Stdout)

	case a.printMD:
		r, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			// Limit line length.
			glamour.WithWordWrap(100),
			// Allow relative links to work properly.
			glamour.WithBaseURL("https://github.com/awilliams/note/blob/main/"),
		)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		if err = n.print(&buf); err != nil {
			return err
		}

		md, err := r.RenderBytes(buf.Bytes())
		if err != nil {
			return err
		}
		_, err = io.Copy(os.Stdout, bytes.NewReader(md))
		return err
	}

	// Default is to edit the note.

	if err := n.create(); err != nil {
		return nil
	}

	return a.editNote(n)
}

func (a *app) editNote(n *note) error {
	args := []string{n.path}

	switch filepath.Base(a.editor) {
	case "vim", "nvim":
		// FYI: To start in 'insert' mode, add '+startinsert' arg.
		// Instruct (n)vim to open file and position cursor at given line.
		if a.weekOffset == 0 {
			line, err := n.rangeCursorLine()
			if err != nil {
				return err
			}
			if line > 0 {
				args = append(args, fmt.Sprintf("+%d", line))
			}
		}
	default:
	}

	cmd := exec.Command(a.editor, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
