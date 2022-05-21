package wapptester

import (
	"context"
	"testing"
)

func TestMakeSample(t *testing.T) {
	t.Log(MakeSample(context.TODO(), "https://zznq.imipy.com/"))
}
