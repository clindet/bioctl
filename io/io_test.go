package io

import (
	"fmt"
	"testing"
)

func TestSetQueryFromEnd(t *testing.T) {
	err := CheckInputFiles([]string{"ls"})
	if err != nil {
		fmt.Println(err)
	}
	err = CheckInputFiles([]string{"sadf", "KM034562v1-tr.fa"})
	if err != nil {
		fmt.Println(err)
	}
}
