package stringo

import (
	"strings"
	"testing"
)

func TestStrExtract(t *testing.T) {
	pat := []string{"mTOR", "AMPK", "dsf", "RNA-Seq", "RNA-seq"}
	patStr := strings.Join(pat, "|")
	res := StrExtract("mTOR, AMPK, DSF, RNA-Seq", patStr, 4)
	if len(res) != 3 {
		log.Fatalf("StrExtract faild: length of res should equals 3")
	}
}
