package main

import (
	"os"
	"strings"
	detailedGoogleTranslate "translate/googletranslate"
	simpleGoogleTranslate "translate/googletranslatefree"
	"translate/termania"
	. "translate/types"
)

func WrapText(s string, limit int) string {
	var sb = strings.Builder{}

	start := 0
	for start+limit < len(s) {
		part := s[start:start+limit]
		partLastSpaceIndex := strings.LastIndex(part, " ")
		if partLastSpaceIndex != -1 {
			part = part[:partLastSpaceIndex]
			start += 1
		}

		start += len(part)

		sb.WriteString(part)
		sb.WriteByte('\n')
	}

	sb.WriteString(s[start:])

	return sb.String()
}

func printTranslate(result TranslateResult, err error) {
	if err != nil {
		println(err.Error())
	} else {
		println(WrapText(result.String(), 40))
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
