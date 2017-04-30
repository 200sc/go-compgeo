package off

import (
	"bufio"
	"strconv"
	"strings"

	compgeo "github.com/200sc/go-compgeo"
)

func readIntLine(s *bufio.Scanner, l int) ([]int, error) {
	var err error
	out := make([]int, l)

	if !s.Scan() {
		return out, compgeo.TypeError{}
	}

	ints := strings.Split(s.Text(), " ")
	if len(ints) < l {
		return nil, compgeo.TypeError{}
	}

	for i := 0; i < l; i++ {
		out[i], err = strconv.Atoi(ints[i])
		if err != nil {
			return nil, compgeo.TypeError{}
		}
	}

	return out, nil
}

func readFloat64Line(s *bufio.Scanner, l int) ([]float64, error) {
	var err error
	out := make([]float64, l)

	if !s.Scan() {
		return out, compgeo.TypeError{}
	}

	ints := strings.Split(s.Text(), " ")
	if len(ints) < l {
		return nil, compgeo.TypeError{}
	}

	for i := 0; i < l; i++ {
		out[i], err = strconv.ParseFloat(ints[i], 64)
		if err != nil {
			return nil, compgeo.TypeError{}
		}
	}

	return out, nil
}

// The number of elements in this line is defined by the first value.
func readIntsLineNoLength(s *bufio.Scanner) (int, []int, error) {
	var err error
	if !s.Scan() {
		return 0, make([]int, 0), compgeo.TypeError{}
	}

	ints := strings.Split(s.Text(), " ")

	length, err := strconv.Atoi(ints[0])
	if err != nil {
		return 0, nil, compgeo.TypeError{}
	}

	out := make([]int, length)

	if len(ints) < (length + 1) {
		return 0, nil, compgeo.TypeError{}
	}

	for i := 0; i < length; i++ {
		out[i], err = strconv.Atoi(ints[i+1])
		if err != nil {
			return 0, nil, compgeo.TypeError{}
		}
	}

	return length, out, nil
}
