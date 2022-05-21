package wapptester

import (
	"testing"
)

func TestSelectWrapper(t *testing.T) {
	NewSelectWrapper(&Sample{
		Title: "xxx",
	}, "wapptester")
}
