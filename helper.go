package wapptester

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"regexp"
	"regexp/syntax"
	"strings"

	"github.com/gokitx/pkgs/bytesconv"
)

func Md5(strList ...string) string {
	str := strings.Join(strList, ":")
	h := md5.New()
	h.Write(bytesconv.StringToBytes(str))
	return hex.EncodeToString(h.Sum(nil))
}

func MMH3(strList ...string) string {
	str := strings.Join(strList, ":")
	h := md5.New()
	h.Write(bytesconv.StringToBytes(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Base64(s string) string {
	return base64.StdEncoding.EncodeToString(bytesconv.StringToBytes(s))
}

func Contains(src string, substr string) bool {
	return strings.Contains(strings.ToLower(src), strings.ToLower(substr))
}

func Equals(s, t string) bool {
	return strings.ToLower(strings.TrimSpace(s)) == strings.ToLower(strings.TrimSpace(t))
}

func Regex(s, pattern string) bool {
	if s == "" {
		return false
	}
	if pattern == "" {
		return true
	}
	regex, err := syntax.Parse(pattern, syntax.DotNL|syntax.Perl|syntax.FoldCase|syntax.WasDollar|syntax.NonGreedy)
	if err != nil {
		return false
	}
	rex, err := regexp.Compile(regex.String())
	if err != nil {
		return false
	}
	return rex.MatchString(s)
}

func Find(s, pattern string, posList ...interface{}) string {
	regex, err := syntax.Parse(pattern, syntax.DotNL|syntax.Perl|syntax.FoldCase|syntax.WasDollar)
	if err != nil {
		return ""
	}
	rex, err := regexp.Compile(regex.String())
	if err != nil {
		return ""
	}

	ret := rex.FindStringSubmatch(s)
	if len(ret) < 2 {
		return ""
	}
	pos := getPos(posList)
	if pos > len(ret) {
		return ""
	}
	if pos < 1 {
		pos = 1
	}
	return ret[pos]
}

func getPos(posList []interface{}) int {
	pos := 1
	if len(posList) <= 0 {
		return pos
	}
	fir := posList[0]
	switch val := fir.(type) {
	case int:
		pos = val
	case int8:
		pos = int(val)
	case int16:
		pos = int(val)
	case int32:
		pos = int(val)
	case int64:
		pos = int(val)
	case float32:
		pos = int(val)
	case float64:
		pos = int(val)
	}
	return pos
}
