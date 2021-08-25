package main

import (
	"fmt"
	"os"
	fullTranslate "translate/googletranslate"
	simpleTranslate "translate/googletranslatefree"
)

func printTranslate(result interface{}, err error) {
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}
}

func main() {
	// curl -L -o ca-bundle.crt https://curl.haxx.se/ca/cacert.pem
	os.Setenv("SSL_CERT_FILE", "/mnt/onboard/.adds/certs/ca-bundle.crt")

	text := os.Args[1]
	isSimple := os.Args[2] == "S"

	if isSimple {
		printTranslate(simpleTranslate.Translate(text))
	} else {
		printTranslate(fullTranslate.Translate(text))
	}
}
