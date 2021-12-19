package main

import (
	"flag"
	. "github.com/krezefal/cox-stegosystem/stego"
	"github.com/sirupsen/logrus"
)

const alpha = 0.1

func embeddingProcedure(src, tg *string, r *int, m *string) {
	logrus.Info("Embedding message procedure starts")

	logrus.Debug("Starting read image")
	sourceImage, err := ReadImage(*src)
	if err != nil {
		logrus.Fatal(err)
	}

	if *r == 0 && *m == "" {
		logrus.Fatal("Empty message")
	}

	var message []byte
	if *r != 0 {
		logrus.Debug("Generating random message")
		message = GenerateRandomMessage(*r)
		logrus.Info("Random message to embed:\n", message)
	} else {
		logrus.Debug("Converting given message to bit sequence")
		message = GenerateMessage(*m)
	}

	logrus.Debug("Embedding message into the image")
	modifiedImage, messageLen := EmbedMessage(sourceImage, message, alpha)
	if messageLen == len(message) {
		logrus.Info("All ", messageLen, " bits of message has been placed into image-container")
	} else {
		logrus.Info("First ", messageLen, " bits of the message has been placed into image-container")
	}

	logrus.Debug("Saving image with embedded message")
	if err := WriteImage(*tg, modifiedImage); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Message embedding finished")
}

func extractingProcedure(src, tg *string) {
	logrus.Info("Extracting message procedure starts")

	logrus.Debug("Starting read images")
	srcImg, err := ReadImage(*src)
	if err != nil {
		logrus.Fatal(err)
	}
	tgImg, err := ReadImage(*tg)
	if err != nil {
		logrus.Fatal(err)
	}

	var message []byte

	logrus.Debug("Extracting message from the image")
	message, err = ExtractMessage(srcImg, tgImg)
	if err != nil {
		logrus.Fatal(err)
	}

	if len(message) == 0 {
		logrus.Info("Spectral coefficients of the images completely coincide. Unable to extract the message")
	} else {
		logrus.Info("Extracted message according to given source:\n", message)
		logrus.Info("Length of this message is ", len(message), " bits")
	}

	logrus.Info("Message extracting finished")
}

func main() {

	src := flag.String("src", "", "Path to source bmp image")
	tg := flag.String("tg", "", "Path to target bmp image")
	ext := flag.Bool("ext", false, "Flag to perform extracting message procedure")
	r := flag.Int("r", 0, "Flag for randomly generating message. Need to specify length")
	m := flag.String("m", "", "Message to embed into image-container")
	dbg := flag.Bool("dbg", false, "Debug mode")

	flag.Parse()

	if *dbg {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Debug mode: ON")
	}

	if *src == "" {
		logrus.Fatal("Empty file path to the source image")
	}

	if *tg == "" {
		logrus.Fatal("Empty file path to the target (container) image")
	}

	if !*ext {
		embeddingProcedure(src, tg, r, m)
	} else {
		extractingProcedure(src, tg)
	}
}
