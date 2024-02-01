package http

import (
	"net/url"

	gcast "github.com/sidgwick/gutil/cast"
	"github.com/spf13/cast"
)

func BuildGetQueryString(input interface{}) string {
	data := url.Values{}

	xInput := gcast.ToXStringMap(input)
	for k, v := range xInput {
		data.Set(k, cast.ToString(v))
	}

	return data.Encode()
}

func MergeQueryString(a string, b string) string {
	qa, _ := url.ParseQuery(a)
	qb, _ := url.ParseQuery(b)

	for kb, _ := range qb {
		if qa.Get(kb) != "" {
			qa.Add(kb, qb.Get(kb))
		}

		qa.Set(kb, qb.Get(kb))
	}

	return qa.Encode()
}
