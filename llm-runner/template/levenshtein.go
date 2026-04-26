package template

func stringEditDistance(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	m, n := len(ra), len(rb)
	if m == 0 {
		return n
	}
	if n == 0 {
		return m
	}

	row := make([]int, n+1)
	for j := 0; j <= n; j++ {
		row[j] = j
	}

	for i := 1; i <= m; i++ {
		prev := row[0]
		row[0] = i
		rai := ra[i-1]
		for j := 1; j <= n; j++ {
			cost := 1
			if rai == rb[j-1] {
				cost = 0
			}
			tmp := row[j]
			row[j] = min(row[j]+1, row[j-1]+1, prev+cost)
			prev = tmp
		}
	}

	return row[n]
}
