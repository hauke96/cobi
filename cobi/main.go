package main

import (
	"github.com/alecthomas/kong"
	"github.com/hauke96/sigolo"
	"path/filepath"
	"strings"
)

var cli struct {
	Debug  bool   `help:"Enable debug mode." short:"d"`
	Input  string `help:"The input file" short:"i" required:"true"`
	Output string `help:"The output file" short:"o" optional:"true"`
}

type Mode int

const (
	ModeCompress   = iota // Turn image to .cobi file
	ModeDecompress        // Turn .cobi file into image
)

func main() {
	kong.Parse(&cli)

	if cli.Debug {
		sigolo.LogLevel = sigolo.LOG_DEBUG
	}

	// Determine what to do (compress or decompress)
	mode := ModeCompress
	if strings.HasSuffix(cli.Input, ".cobi") {
		mode = ModeDecompress
	}

	// Determine output filename and set it in the "cli" instance
	if cli.Output == "" {
		inputFilename := strings.TrimSuffix(cli.Input, filepath.Ext(cli.Input))
		outputFilename := inputFilename + ".cobi" // compression is default
		if mode == ModeDecompress {
			outputFilename = inputFilename + ".png"
		}
		cli.Output = outputFilename
	}
}
