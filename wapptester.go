package wapptester

import (
	"context"
	"strings"

	"github.com/PaesslerAG/gval"
)

func Match(ctx context.Context, URL, expression string) (interface{}, error) {
	sample, err := MakeSample(ctx, URL)
	if err != nil {
		return nil, err
	}

	v, err := NewSelectWrapper(sample, "wapptester")
	if err != nil {
		return nil, err
	}

	return gval.Evaluate(expression, map[string]any{
		"resp": v,

		"md5":    Md5,
		"mmh3":   MMH3,
		"base64": Base64,

		"regex":    Regex,
		"find":     Find,
		"contains": Contains,
		"equals":   Equals,
		"starts":   strings.HasPrefix,
		"ends":     strings.HasSuffix,
		"lower":    strings.ToLower,
		"upper":    strings.ToUpper,
		"trim":     strings.TrimSpace,
	})
}
