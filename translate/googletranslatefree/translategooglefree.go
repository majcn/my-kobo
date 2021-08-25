// Package translategooglefree:  https://github.com/bas24/googletranslatefree

package translategooglefree

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	. "translate/types"
)

func Translate(source string) (result TranslateResult, err error) {
	sourceLang := "en"
	targetLang := "sl"

	reqUrl := fmt.Sprintf(`https://translate.googleapis.com/translate_a/single?client=gtx&dt=t&sl=%s&tl=%s&q=%s`, sourceLang, targetLang, url.QueryEscape(source))

	r, err := http.Get(reqUrl)
	if err != nil {
		time.Sleep(time.Second)
		r, err = http.Get(reqUrl)

		if err != nil {
			err = errors.New("error getting translate.googleapis.com")
			return
		}
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.New("error reading response body")
		return
	}

	resultAsText := gjson.GetBytes(body, "0.0.0").String()
	if len(resultAsText) > 0 {
		result.Translation = resultAsText
		return
	} else {
		err = errors.New("no translated data in response")
		return
	}
}
