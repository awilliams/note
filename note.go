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

func newNote(baseDir string, r dateRange) *note {
	p := filepath.Join(
		baseDir,
		strconv.Itoa(r.startYear),
		fmt.Sprintf("%02d.md", r.startWeek),
	)
	return &note{
		r:    r,
		path: p,
	}
}

type note struct {
	r    dateRange
	path string
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
	if err := os.MkdirAll(filepath.Dir(n.path), 0700); err != nil {
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
	fd, err := os.OpenFile(n.path, flags, 0700)
	if err != nil {
		return err
	}
	defer fd.Close()

	return n.writeTemplate(fd)
}

func (n *note) rangeCursorLine() (int, error) {
	fd, err := os.Open(n.path)
	if err != nil {
		return 0, err
	}
	defer fd.Close()

	return n.cursorLine(fd, n.r.date)
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

func (n *note) writeTemplate(w io.Writer) error {
	_, err := fmt.Fprintf(w, "# Week %02d, %d\n---\n\n", n.r.startWeek, n.r.startYear)
	if err != nil {
		return err
	}

	for _, day := range n.r.days {
		if _, err := fmt.Fprintf(w, "%s\n\n", n.heading(day)); err != nil {
			return err
		}
	}

	return nil
}

func (n *note) heading(date time.Time) string {
	return date.Format("## Monday, 02-January-2006")
}
