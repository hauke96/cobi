package main

import (
	"cobi/encoding"
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
		// Determine correct reader for the input image
		var reader image.Reader
		switch inputFileExt {
		case ".png":
			reader = &png.Reader{}
		default:
			sigolo.Fatal("Unsupported file extension %s for compression", inputFileExt)
		}

		// Compress the image
		encodedAreas, err := compress(cli.Input, reader)
		sigolo.FatalCheck(err)
		err = encoding.Write(cli.Output, encodedAreas)
		sigolo.FatalCheck(err)

		if cli.Debug {
			decodedImage, err := encoding.Decode(encodedAreas)
			sigolo.FatalCheck(err)

			pngWriter := png.Writer{}
			err = pngWriter.Write(inputFileName+"_decoded.png", *decodedImage)
			sigolo.FatalCheck(err)
			err = pngWriter.Write(inputFileName+"_decoded_debug.png", *encoding.GetDebugImage(decodedImage.Width, decodedImage.Height, encodedAreas))
			sigolo.FatalCheck(err)
		}
	case ModeDecompress:
		sigolo.Fatal("Decompression is not implemented yet")
	default:
		sigolo.Fatal("Invalid compression mode %d", mode)
	}
}

func compress(filePath string, reader image.Reader) ([4][]encoding.EncodedArea, error) {
	img, err := reader.Read(filePath)
	if err != nil {
		return [4][]encoding.EncodedArea{}, err
	}

	//if sigolo.LogLevel == sigolo.LOG_DEBUG {
	//	img.Print()
	//}

	encodedAreas, err := encoding.Encode(*img)
	if err != nil {
		return [4][]encoding.EncodedArea{}, err
	}

	return encodedAreas, nil
}
