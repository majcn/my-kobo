// Package googletranslate:  https://github.com/FreddieDeWitt/extended-google-translate-api/
package googletranslate

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
	. "translate/types"
)

func GetRawObjectGetParams(baseUrl string) (result url.Values, err error) {
	r, err := http.Get(baseUrl)
	if err != nil {
		time.Sleep(time.Second)
		r, err = http.Get(baseUrl)

		if err != nil {
			err = errors.New("error getting translate.google.com")
		}
	}
	defer r.Body.Close()

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.New("error reading response body")
		return
	}

	result = url.Values{
		"rpcids": []string{"MkEWBc"},
		"f.sid": []string{string(regexp.MustCompile(`"FdrFJe":"(.*?)"`).FindSubmatch(bodyBytes)[1])},
		"bl": []string{string(regexp.MustCompile(`"cfb2h":"(.*?)"`).FindSubmatch(bodyBytes)[1])},
		"hl": []string{"en-US"},
		"soc-app": []string{"1"},
		"soc-platform": []string{"1"},
		"soc-device": []string{"1"},
		"_reqid": []string{strconv.Itoa(1000 + rand.Intn(9000))},
		"rt": []string{"c"},
	}

	return
}

func GetRawObject(source, sourceLang, targetLang string) (result []byte, err error) {
	baseUrl := "https://translate.google.com"

	params, err := GetRawObjectGetParams(baseUrl)
	if err != nil {
		return
	}

	requestUrl := fmt.Sprintf(`%s/_/TranslateWebserverUi/data/batchexecute?%s`, baseUrl, params.Encode())
	requestBody := url.Values{
		"f.req": []string{fmt.Sprintf(`[[["MkEWBc","[[\"%s\",\"%s\",\"%s\",true],[null]]",null,"generic"]]]`, source, sourceLang, targetLang)},
	}

	resp, err := http.PostForm(requestUrl, requestBody)
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	lengthLength := bytes.IndexByte(bodyBytes[6:], '\n')
	length := bodyBytes[6:6+lengthLength]
	lengthAsInt, err := strconv.Atoi(string(length))

	result = bodyBytes[7+lengthLength:8+lengthLength+lengthAsInt]

	return
}

func Translate(source string) (result TranslateResult, err error) {
	sourceLang := "en"
	targetLang := "sl"

	rawObjectJson, err := GetRawObject(source, sourceLang, targetLang)
	rawObjectJsonData := gjson.GetBytes(rawObjectJson, "0.2").String()

	result.Translation = gjson.Get(rawObjectJsonData, "1.0.0.5.0.0").String()

	result.Details = make(map[string][]string)
	for _, v := range gjson.Get(rawObjectJsonData, "3.5.0").Array() {
		translateType := gjson.Get(v.Raw, "0").String()
		for _, x := range gjson.Get(v.Raw, "1.#.0").Array() {
			result.Details[translateType] = append(result.Details[translateType], x.String())
		}
	}

	return
}
