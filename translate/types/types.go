package types

import "strings"

type TranslateResult struct {
	Translation string
	Details map[string][]string
}

func (r TranslateResult) String() string {
	var result []string

	result = append(result, r.Translation)
	result = append(result, "")

	for k, v := range r.Details {
		result = append(result, k)
		for _, vv := range v {
			result = append(result, "  " + vv)
		}
		result = append(result, "")
	}

	return strings.Join(result, "\n")
}
