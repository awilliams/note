package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	day  = 24 * time.Hour
	week = 7 * day
)

//go:embed VERSION
var version string

func main() {
	var cfg config
	if err := cfg.parse(os.Args); err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}

	if cfg.version {
		v := version
		if v == "" {
			v = "DEV"
		} else {
			v = "v" + v
		}
		fmt.Print(v)
		return
	}

	a := app{
		date:       time.Now().Add(time.Duration(cfg.weekOffset) * week),
		weekOffset: cfg.weekOffset,
		names:      cfg.names,
		baseDir:    cfg.baseDir,
		editor:     cfg.editor,
	}
	if err := a.run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
}

type app struct {
	date       time.Time
	weekOffset int
	names      bool
	baseDir    string
	editor     string
}

func (a *app) run() error {
	year, week := a.date.ISOWeek()
	notePath := filepath.Join(
		a.baseDir,
		strconv.Itoa(year),
		fmt.Sprintf("%02d.md", week),
	)
	if a.names {
		fmt.Println(notePath)
		return nil
	}

	if err := a.createNote(notePath); err != nil {
		return nil
	}

	return a.editNote(notePath)
}

func (a *app) editNote(notePath string) error {
	args := []string{notePath}

	switch filepath.Base(a.editor) {
	case "vim", "nvim":
		// Instruct (n)vim to open file at given line.
		if line := a.curLine(notePath); line > 0 {
			args = append(args, fmt.Sprintf("+%d", line))
		}
	default:
	}

	cmd := exec.Command(a.editor, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (a *app) curLine(notePath string) int {
	if a.weekOffset != 0 {
		// When opening a note from other than current week,
		// don't open to a particular line.
		return 0
	}

	fd, err := os.Open(notePath)
	if err != nil {
		return 0
	}
	defer fd.Close()

	nextHeader := a.heading(a.date.Add(day))
	var line int

	s := bufio.NewScanner(fd)
	for s.Scan() {
		line++
		if strings.HasPrefix(s.Text(), nextHeader) {
			return line - 1
		}
	}

	return line
}

func (a *app) createNote(notePath string) error {
	if err := os.MkdirAll(filepath.Dir(notePath), 0700); err != nil {
		return err
	}

	var (
		writeTemplate bool
		exists        bool
	)
	s, err := os.Stat(notePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		// File does not exist.
		writeTemplate = true
	} else {
		exists = true
		// Only update empty files with template.
		writeTemplate = s.Size() == 0
	}

	if !writeTemplate {
		return nil
	}

	flags := os.O_RDWR | os.O_CREATE
	if !exists {
		flags |= os.O_EXCL
	}
	fd, err := os.OpenFile(notePath, flags, 0700)
	if err != nil {
		return err
	}
	defer fd.Close()

	return a.writeTemplate(fd)
}

func (a *app) writeTemplate(w io.Writer) error {
	monday := a.date.Add(time.Duration(time.Monday-a.date.Weekday()) * day)
	for d := time.Duration(0); d < 5; d++ {
		day := monday.Add(d * day)
		if _, err := fmt.Fprintf(w, "%s\n\n", a.heading(day)); err != nil {
			return err
		}
	}

	return nil
}

func (a *app) heading(t time.Time) string {
	return t.Local().Format("### Monday, 02-January-2006")
}
