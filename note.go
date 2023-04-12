package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func notePath(baseDir string, r dateRange) string {
	return filepath.Join(
		baseDir,
		strconv.Itoa(r.startYear),
		fmt.Sprintf("%02d.md", r.startWeek),
	)
}

func newNote(baseDir string, curRange, prevRange dateRange) *note {
	return &note{
		curRange:  curRange,
		prevRange: prevRange,
		baseDir:   baseDir,
		path:      notePath(baseDir, curRange),
	}
}

type note struct {
	curRange  dateRange
	prevRange dateRange
	baseDir   string
	path      string
}

func (n *note) print(w io.Writer) error {
	fd, err := os.Open(n.path)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = io.Copy(w, fd)
	return err
}

func (n *note) create() error {
	if err := os.MkdirAll(filepath.Dir(n.path), 0o700); err != nil {
		return err
	}

	var (
		exists        bool
		writeTemplate bool
	)

	stat, err := os.Stat(n.path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		// Note file does not exist.
		writeTemplate = true
	} else {
		// File exists.
		// Only update empty files with template.
		writeTemplate = stat.Size() == 0
	}

	if !writeTemplate {
		return nil
	}

	flags := os.O_RDWR | os.O_CREATE
	if !exists {
		flags |= os.O_EXCL
	}
	fd, err := os.OpenFile(n.path, flags, 0o700)
	if err != nil {
		return err
	}
	defer fd.Close()

	// Read previous note's TODO section, if any.
	var prevTODO string
	prevFd, err := os.Open(notePath(n.baseDir, n.prevRange))
	if err == nil {
		s := bufio.NewScanner(prevFd)
		var buf strings.Builder
		var postTODO bool
		for s.Scan() {
			if strings.HasPrefix(s.Text(), "TODO:") {
				postTODO = true
				continue
			}
			if !postTODO {
				continue
			}
			buf.Write(s.Bytes())
			buf.WriteRune('\n')
		}
		prevTODO = buf.String()
	}

	return n.writeTemplate(fd, prevTODO)
}

func (n *note) rangeCursorLine() (int, error) {
	fd, err := os.Open(n.path)
	if err != nil {
		return 0, err
	}
	defer fd.Close()

	return n.cursorLine(fd, n.curRange.date)
}

func (n *note) cursorLine(r io.Reader, date time.Time) (int, error) {
	prefix := n.heading(date)

	var (
		line  int
		found bool
	)
	s := bufio.NewScanner(r)
	for s.Scan() {
		line++
		t := s.Text()
		if !found {
			if strings.HasPrefix(t, prefix) {
				found = true
			}
			continue
		}
		// Find next blank line or heading
		if t == "" {
			return line, nil
		}
		if strings.HasPrefix(t, "#") {
			return line - 1, nil
		}
	}

	return 0, s.Err()
}

func (n *note) writeTemplate(w io.Writer, prevTODO string) error {
	_, err := fmt.Fprintf(w, "# Week %02d, %d\n\n---\n\n", n.curRange.startWeek, n.curRange.startYear)
	if err != nil {
		return err
	}

	for _, day := range n.curRange.days {
		if _, err := fmt.Fprintf(w, "%s\n\n", n.heading(day)); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(w, "---\n\nTODO:\n%s\n", prevTODO); err != nil {
		return err
	}

	return nil
}

func (n *note) heading(date time.Time) string {
	return date.Format("## Monday, 02-January-2006")
}
