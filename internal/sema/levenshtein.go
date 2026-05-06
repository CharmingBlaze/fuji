package sema

func levenshtein(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}
	row := make([]int, len(b)+1)
	for j := range row {
		row[j] = j
	}
	for i, ca := range a {
		prev := i + 1
		for j, cb := range b {
			next := row[j]
			if ca != cb {
				next++
			}
			if row[j+1]+1 < next {
				next = row[j+1] + 1
			}
			if prev+1 < next {
				next = prev + 1
			}
			row[j], prev = prev, next
		}
		row[len(b)] = prev
	}
	return row[len(b)]
}
