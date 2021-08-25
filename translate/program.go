package main

import (
	"os"
	detailedGoogleTranslate "translate/googletranslate"
	simpleGoogleTranslate "translate/googletranslatefree"
	"translate/termania"
	. "translate/types"
)

func printTranslate(result TranslateResult, err error) {
	if err != nil {
		println(err.Error())
	} else {
		println(result.String())
	}
}

func main() {
	// curl -L -o ca-bundle.crt https://curl.haxx.se/ca/cacert.pem
	os.Setenv("SSL_CERT_FILE", "/mnt/onboard/.adds/certs/ca-bundle.crt")

	text := os.Args[1]
	switch os.Args[2] {
	case "GS":
		printTranslate(simpleGoogleTranslate.Translate(text))
	case "GD":
		printTranslate(detailedGoogleTranslate.Translate(text))
	case "T":
		printTranslate(termania.Translate(text))
	}
}
