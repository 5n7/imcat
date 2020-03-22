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
	Scale bool `short:"s" long:"scale" description:"Show with scaling"`
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
	if err != nil {
		return 2
	}

	paths := args

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

	xSize := tx * sfx
	ySize := ty * sfy

	if opt.Scale {
		xSize /= len(paths)
		ySize /= len(paths)
	}

	for _, path := range paths {
		if matched, _ := regexp.MatchString("^https?://", path); matched {
			pix, err = ansimage.NewScaledFromURL(path, ySize, xSize, bg, sm, dm)
			if err != nil {
				return 1
			}
		} else {
			pix, err = ansimage.NewScaledFromFile(path, ySize, xSize, bg, sm, dm)
			if err != nil {
				return 1
			}
		}

		pix.Draw()
	}

	return 0
}

func main() {
	os.Exit(run(os.Args[1:]))
}
