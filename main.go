package main

import (
	"os"
	"regexp"

	"github.com/eliukblau/pixterm/pkg/ansimage"
	"github.com/jessevdk/go-flags"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/crypto/ssh/terminal"
)

// Option represents application options
type Option struct {
}

func terminalSize() (width, height int, err error) {
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		return terminal.GetSize(int(os.Stdout.Fd()))
	}
	return 80, 24, nil
}

func run(args []string) int {
	var (
		opt Option
		pix *ansimage.ANSImage
	)

	args, err := flags.ParseArgs(&opt, args)
	if err != nil || len(args) != 1 {
		return 2
	}

	path := args[0]

	tx, ty, err := terminalSize()
	if err != nil {
		return 1
	}

	sfx, sfy := 1, 2

	bg, err := colorful.Hex("#000000")
	if err != nil {
		return 1
	}

	sm := ansimage.ScaleModeFit

	dm := ansimage.NoDithering

	if matched, _ := regexp.MatchString("^https?://", path); matched {
		pix, err = ansimage.NewScaledFromURL(path, ty*sfy, tx*sfx, bg, sm, dm)
		if err != nil {
			return 1
		}
	} else {
		pix, err = ansimage.NewScaledFromFile(path, ty*sfy, tx*sfx, bg, sm, dm)
		if err != nil {
			return 1
		}
	}

	pix.Draw()

	return 0
}

func main() {
	os.Exit(run(os.Args[1:]))
}
