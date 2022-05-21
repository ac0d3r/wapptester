package wapptester

import (
	"context"
	"testing"
)

func TestMatch(t *testing.T) {
	t.Log(Match(context.TODO(), "https://zznq.imipy.com/",
		`resp.status ==200 && contains(resp.body, "hexo")`))
}
