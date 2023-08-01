package main

import (
	"cobi/image"
	"cobi/png"
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

	// Determine file name and extention
	inputFileExt := filepath.Ext(cli.Input)
	inputFileName := strings.TrimSuffix(cli.Input, inputFileExt)

	// Determine output filename and set it in the "cli" instance
	if cli.Output == "" {
		outputFilename := inputFileName + ".cobi" // compression is default
		if mode == ModeDecompress {
			outputFilename = inputFileName + ".png"
		}
		cli.Output = outputFilename
	}

	switch mode {
	case ModeCompress:
		var reader image.Reader
		switch inputFileExt {
		case ".png":
			reader = &png.Reader{}
		default:
			sigolo.Fatal("Unsupported file extension %s for compression", inputFileExt)
		}
		_, err := compress(cli.Input, reader)
		sigolo.FatalCheck(err)
	case ModeDecompress:
		sigolo.Fatal("Decompression is not implemented yet")
	default:
		sigolo.Fatal("Invalid compression mode %d", mode)
	}
}

func compress(filePath string, reader image.Reader) ([]byte, error) {
	img, err := reader.Read(filePath)
	if err != nil {
		return nil, err
	}

	if sigolo.LogLevel == sigolo.LOG_DEBUG {
		img.Print()
	}

	return make([]byte, 0), nil
}
