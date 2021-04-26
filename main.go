package main

import (
	_ "embed"
	"fmt"
	"os"
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
		baseDir:    cfg.baseDir,
		weekOffset: cfg.weekOffset,
		names:      cfg.names,
		print:      cfg.print,
		printMD:    cfg.printMD,
		editor:     cfg.editor,
	}
	if err := a.run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
}
