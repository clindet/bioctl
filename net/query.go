package net

// SetQueryFromEnd get from and end according to the from, size and total
func SetQueryFromEnd(from int, size int, total int) (int, int) {
	if size == -1 {
		size = total + 1
	}
	end := from + size
	if end == -1 || end > total {
		end = total + 1
	}
	if from < 0 {
		from = 0
	} else if from > total {
		from = total
	}
	if end <= from {
		end = from + 1
	}
	return from, end
}
