package dcel

import "fmt"

type FullEdge [2]Point

func (fe FullEdge) Lesser(d int) Point {
	if fe[0].Val(d) < fe[1].Val(d) {
		return fe[0]
	}
	return fe[1]
}

func (fe FullEdge) Greater(d int) Point {
	if fe[0].Val(d) > fe[1].Val(d) {
		return fe[0]
	}
	return fe[1]
}

func (fe FullEdge) Left() Point {
	return fe.Lesser(0)
}

func (fe FullEdge) Right() Point {
	return fe.Greater(0)
}

func (fe FullEdge) Top() Point {
	return fe.Greater(1)
}

func (fe FullEdge) Bottom() Point {
	return fe.Lesser(1)
}

func (fe FullEdge) Inner() Point {
	return fe.Greater(2)
}

func (fe FullEdge) Outer() Point {
	return fe.Lesser(2)
}

func (fe FullEdge) Slope() float64 {
	p1 := fe.Left()
	p2 := fe.Right()
	return (p2.Y() - p1.Y()) / (p2.X() - p2.X())
}

func (fe FullEdge) PointAt(d int, v float64) (Point, error) {
	e1 := fe.Lesser(d)
	e2 := fe.Greater(d)
	if v < e1[d] || v > e2[d] {
		fmt.Println(v, e1[d], e2[d])
		return Point{}, RangeError{}
	}
	v -= e1[d]
	span := e2[d] - e1[d]
	portion := v / span
	p := Point{}
	for i := 0; i < len(p); i++ {
		if i == d {
			p[i] = v
		} else {
			p[i] = e1[i] + (portion * (e2[i] - e1[i]))
		}
	}
	return p, nil
}

func (fe FullEdge) SubEdge(d int, v1, v2 float64) (FullEdge, error) {
	p1, err := fe.PointAt(d, v1)
	if err != nil {
		return FullEdge{}, err
	}
	p2, err := fe.PointAt(d, v2)
	if err != nil {
		return FullEdge{}, err
	}
	return FullEdge{p1, p2}, nil
}
