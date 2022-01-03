package main

import (
	"flag"
	. "github.com/krezefal/cox-stegosystem/bmp"
	. "github.com/krezefal/cox-stegosystem/stego"
	"github.com/sirupsen/logrus"
	"math"
	"os"
)

const alpha = 0.5

func embeddingProcedure(src, tg *string, psnrFlag *bool, r *int, m *string) {
	logrus.Info("Embedding message procedure starts")

	logrus.Debug("Starting read image")
	srcImg, err := ReadImage(*src)
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

	logrus.Debug("Embedding message into image")
	modifiedImg, messageLen := EmbedMessage(srcImg, message, alpha)
	if messageLen == len(message) {
		logrus.Info("All ", messageLen, " bits of message has been placed into image-container")
	} else {
		logrus.Info("First ", messageLen, " bits of the message has been placed into image-container")
	}

	logrus.Debug("Saving image with embedded message only")
	if err := WriteImage(*tg, modifiedImg); err != nil {
		logrus.Fatal(err)
	}

	logrus.Debug("Embedding message length into image")
	if file, err := os.ReadFile(*tg); err != nil {
		logrus.Fatal(err)
	} else {
		file = EmbedMessageLen(file, messageLen)
		logrus.Debug("Saving image with embedded message and its metadata")
		if err := os.WriteFile(*tg, file, 0644); err != nil {
			logrus.Fatal(err)
		}
	}

	logrus.Info("Message embedding finished")

	if *psnrFlag {
		logrus.Debug("Calculating PSNR")
		psnr, err := PSNRimg(srcImg, modifiedImg)
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Info("PSNR(src, tg): \033[31m", psnr, "\033[0m")
	}
}

func extractingProcedure(src, tg *string, psnrFlag, forceFlag *bool) {
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
	tgFile, err := os.ReadFile(*tg)
	if err != nil {
		logrus.Fatal(err)
	}

	var message []byte

	logrus.Debug("Detecting message embedding in target")
	detectionFlag := DetectEmbedding(tgFile)
	if detectionFlag || *forceFlag {

		var messageLen int

		if detectionFlag {
			logrus.Debug("Message embedding confirmed. Begin extracting")
			messageLen = ExtractMessageLen(tgFile)
		} else {
			logrus.Debug("Detecting failed. Force flag provided. Continue...")
			messageLen = math.MaxInt32
		}

		message, err = ExtractMessage(srcImg, tgImg, messageLen)
		if err != nil {
			logrus.Fatal(err)
		}

		if len(message) == 0 {
			logrus.Info("Spectral coefficients of the images completely coincide. Unable to extract the message")
		} else {
			logrus.Info("Extracted message according to given source:\n", message)
			logrus.Info("Length of this message is ", len(message), " bits")
		}
	} else {
		logrus.Debug("Detecting failed")
		logrus.Info("Unable to extract message. " +
			"Specified image either does not contain hidden message or file header has been corrupted")
	}

	logrus.Info("Message extracting finished")

	if *psnrFlag {
		logrus.Debug("Calculating PSNR")
		psnr, err := PSNRimg(srcImg, tgImg)
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Info("PSNR(src, tg): \033[31m", psnr, "\033[0m")
	}
}

func main() {

	src := flag.String("src", "", "Path to source bmp image")
	tg := flag.String("tg", "", "Path to target bmp image")
	ext := flag.Bool("ext", false, "Flag to perform extracting message procedure")
	psnr := flag.Bool("psnr", false, "Flag to calculate PSNR between source and target images")
	r := flag.Int("r", 0, "Flag for randomly generating message. Need to specify length")
	m := flag.String("m", "", "Message to embed into image-container")
	dbg := flag.Bool("dbg", false, "Debug mode")
	force := flag.Bool("force", false, "Forcefully skip embedding detection step while extracting "+
		"procedure.\nSee the 'Skipping embedding detection' README subsection for details")

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
		embeddingProcedure(src, tg, psnr, r, m)
	} else {
		extractingProcedure(src, tg, psnr, force)
	}
}
