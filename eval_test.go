package fval

import (
	"testing"
)

func Test_Evaluate(t *testing.T) {
	t.Log(
		Evaluate(`staus_code == 200 && regex(get(), "xxasd")`, map[string]interface{}{
			"a":          1.1,
			"b":          "1.1",
			"get":        get,
			"regex":      regex,
			"staus_code": 200,
		}),
	)
}

func Test_Evaluate2(t *testing.T) {
	t.Log(
		Evaluate(`a == 1.1 || b == "1.1"`, map[string]interface{}{
			"a": 1.1,
			"b": "1.1",
		}),
	)
}

func Test_Evaluate3(t *testing.T) {
	t.Log(
		Evaluate(`get() == "12"`, map[string]interface{}{
			"get": get,
		}),
	)
}

func get() string {
	return "12"
}

func regex(s1, s2 string) string {
	return s1 + s2
}
