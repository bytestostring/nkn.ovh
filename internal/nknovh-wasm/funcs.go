package nknovh_wasm
import (
	"strconv"
	"regexp"
)
func NumSeparate(n int) string {
	s := strconv.Itoa(n)
	re := regexp.MustCompile("(\\d+)(\\d{3})")
	var i string
	for ; i != s; {
        i = s
        s = re.ReplaceAllString(s, "$1,$2")
    }
    return s
}
