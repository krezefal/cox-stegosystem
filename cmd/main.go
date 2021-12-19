package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"os"
)

const keyLen = 8
const maxInterval = 8

func embeddingProcedure(src, tg *string, r *int, m *string) {
	logrus.Info("Embedding message procedure starts")

	if *tg == "" {
		logrus.Fatal("Empty file path of the target image")
	}

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
		logrus.Debug("Converting given message to byte sequence")
		message = GenerateMessage(*m)
	}

	logrus.Debug("Generating random key for intervals")
	key := GenerateKey(keyLen, maxInterval)

	logrus.Debug("Embedding message into image according to random intervals")
	modifiedImage, messageLen := EmbedMessage(sourceImage, message, key)
	if messageLen == len(message) {
		logrus.Info("All ", messageLen, " bits of message has been placed into image-container")
	} else {
		logrus.Info("First ", messageLen, " bits of the message has been placed into image-container")
	}

	logrus.Debug("Saving image with embedded message only")
	if err := WriteImage(*tg, modifiedImage); err != nil {
		logrus.Fatal(err)
	}

	logrus.Debug("Embedding message metadata into image")
	if file, err := os.ReadFile(*tg); err != nil {
		logrus.Fatal(err)
	} else {
		file = EmbedMetadata(file, uint32(messageLen), key)
		logrus.Debug("Saving image with embedded message and its metadata")
		if err := os.WriteFile(*tg, file, 0644); err != nil {
			logrus.Fatal(err)
		}
	}

	logrus.Info("Message embedding finished")
}

func extractingProcedure(src *string) {
	logrus.Info("Extracting message procedure starts")

	logrus.Debug("Starting read image")
	file, err := os.ReadFile(*src)
	if err != nil {
		logrus.Fatal(err)
	}

	var message []byte

	logrus.Debug("Detecting message embedding")
	if DetectEmbedding(file, keyLen) {
		logrus.Debug("Message embedding confirmed. Begin extracting")
		img, err := ReadImage(*src)
		if err != nil {
			logrus.Fatal(err)
		}

		messageLen, key := ExtractMetadata(file, keyLen)
		message = ExtractMessage(img, messageLen, key)

		logrus.Info("Extracted message:\n", message)
		logrus.Info("Length of this message is ", len(message), " bits")
	} else {
		logrus.Debug("Detecting failed")
		logrus.Info("Unable to extract message. " +
			"Specified image either does not contain hidden message or file has been corrupted")
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
		logrus.Fatal("Empty file path of the source image")
	}

	if !*ext {
		embeddingProcedure(src, tg, r, m)
	} else {
		extractingProcedure(src)
	}
}
