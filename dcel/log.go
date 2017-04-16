package dcel

import "math"

func logStarN(n int) int {
	v := float64(n)
	i := 0
	for v >= 1 {
		v = math.Log2(v)
		i++
	}
	return i - 1
}

func logN(n, limit int) int {
	v := float64(n)
	for i := 0; i < limit; i++ {
		v = math.Log2(v)
	}
	return int(math.Ceil(float64(n) / v))
}
