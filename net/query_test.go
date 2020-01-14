package net

import (
	"testing"
)

func TestSetQueryFromEnd(t *testing.T) {
	from, end := SetQueryFromEnd(-1, -1, 1)
	if from != 0 || end != 1 {
		log.Printf("SetQueryFromEnd(-1, -1, 1)  returns %d, %d", from, end)
		log.Fatalln("SetQueryFromEnd(-1, -1, 1) should returns 0, 1")
	}
	from, end = SetQueryFromEnd(0, 1, 1)
	if from != 0 || end != 1 {
		log.Printf("SetQueryFromEnd(0, 1, 1)  returns %d, %d", from, end)
		log.Fatalln("SetQueryFromEnd(0, 1, 1) should returns 0, 1")
	}
	from, end = SetQueryFromEnd(0, 3, 3)
	if from != 0 || end != 3 {
		log.Printf("SetQueryFromEnd(0, 3, 3)  returns %d, %d", from, end)
		log.Fatalln("SetQueryFromEnd(0, 3, 3) should returns 0, 3")
	}
	from, end = SetQueryFromEnd(0, 3, 4)
	if from != 0 || end != 3 {
		log.Printf("SetQueryFromEnd(0, 3, 4)  returns %d, %d", from, end)
		log.Fatalln("SetQueryFromEnd(0, 3, 4) should returns 0, 3")
	}
}
