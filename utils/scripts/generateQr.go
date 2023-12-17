package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/skip2/go-qrcode"
	"image/png"
	_ "image/png"
	"os"
)

func main() {
	var uid, aesKey string
	flag.StringVar(&uid, "uid", "", "User ID")
	flag.StringVar(&aesKey, "aesKey", "", "AES Key")
	flag.Parse()

	if uid == "" || aesKey == "" {
		fmt.Println("Please provide both 'uid' and 'aesKey' as command-line arguments.")
		os.Exit(1)
	}

	combinedString := fmt.Sprintf("%s.%s", uid, aesKey)

	// Base64 encode the combined string
	base64Encoded := base64.StdEncoding.EncodeToString([]byte(combinedString))

	qrCode, err := qrcode.New(base64Encoded, qrcode.Low)
	if err != nil {
		fmt.Println("Error generating QR code:", err)
		os.Exit(1)
	}

	err = saveQRCodeAsPNG(qrCode, "qrcode.png")
	if err != nil {
		fmt.Println("Error saving QR code as PNG:", err)
		os.Exit(1)
	}

	fmt.Println("QR code generated and saved as 'qrcode.png'.")
}

func saveQRCodeAsPNG(qrCode *qrcode.QRCode, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	qrImage := qrCode.Image(256)

	err = png.Encode(file, qrImage)
	if err != nil {
		return err
	}

	return nil
}
