package cmd

import (
	"fmt"
	"testing"
)

func TestConvertIndex(t *testing.T) {
	fmt.Println(convertIndex("1:10;14"))
}
