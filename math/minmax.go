package math

// Min returns the minimum between two integer
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// Max returns the maximum between two integer
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
