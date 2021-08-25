package termania

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"regexp"
	"time"
	. "translate/types"
)

func Translate(source string) (result TranslateResult, err error) {
	reqUrl := fmt.Sprintf(`https://www.termania.net/iskanje?SearchIn=Linked&ld=70&sl=2&tl=61&query=%s`, url.QueryEscape(source))

	r, err := http.Get(reqUrl)
	if err != nil {
		time.Sleep(time.Second)
		r, err = http.Get(reqUrl)

		if err != nil {
			err = errors.New("error getting termania.net")
		}
	}
	defer r.Body.Close()

	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return result, errors.New("error loading html document")
	}

	result.Translation = doc.Find(".lang_sl").Next().Contents().First().Text()

	result.Details = make(map[string][]string)
	doc.Find(".lang_en").Each(func(_ int, selection *goquery.Selection) {
		if source != selection.Parent().Find("h4").Text() {
			return
		}

		contentSpan := selection.Parent().Find(".content span")

		translateType := regexp.MustCompile(`^\s*\(([a-z]+)\)`).FindStringSubmatch(contentSpan.Text())[1]

		for _, s := range contentSpan.Find(".lang_sl").First().NextAll().Nodes {
			if s.FirstChild == nil {
				break
			}

			result.Details[translateType] = append(result.Details[translateType], s.FirstChild.Data)
		}
	})

	return result, nil
}
