// Package translategooglefree:  https://github.com/bas24/googletranslatefree

package translategooglefree

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Translate(source string) (string, error) {
	sourceLang := "en"
	targetLang := "sl"

	reqUrl := fmt.Sprintf(`https://translate.googleapis.com/translate_a/single?client=gtx&dt=t&sl=%s&tl=%s&q=%s`, sourceLang, targetLang, url.QueryEscape(source))

	r, err := http.Get(reqUrl)
	if err != nil {
		return "err", errors.New("error getting translate.googleapis.com")
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "err", errors.New("error reading response body")
	}

	result := gjson.GetBytes(body, "0.0.0").String()
	if len(result) > 0 {
		return result, nil
	} else {
		return "err", errors.New("no translated data in response")
	}
}
