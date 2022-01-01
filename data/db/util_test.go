package db

import (
	"fmt"
	"math"
	"testing"

	"github.com/gonum/stat"
)

func TestDate(t *testing.T) {
	data := GetToday()
	fmt.Println(data)

	d, err := DaysFromToday("2020-10-11")
	fmt.Println(d, err)
}

func TestHV(t *testing.T) {
	//data := []float64{600, 470, 170, 430, 300}
	data := []float64{10.0, 10.1, 10.0, 10.0, 10.0}
	var m float64 = 0.0
	for i := range data {
		m += data[i]
	}
	m = m / 5.0
	fmt.Println("means is: ", m)

	var s float64 = 0.0
	for i := range data {
		s += (m - data[i]) * (m - data[i])
	}
	fmt.Printf("sum of diff power 2: %.4f\n", s)
	fmt.Printf("average of above is: %.5f\nsquare it: %.5f\n", s/4.0, math.Sqrt(s/4.0))

	means, sd := stat.MeanStdDev(data, nil)
	fmt.Println("StdDev ", means, sd)
}

func TestArray(t *testing.T) {
	org := make([]*int, 5)
	index := 2
	st := org[index]
	org[index] = nil
	i := 5
	st = &i
	fmt.Println(*st, org[index])
}

type A struct{}

func (a *A) p() {
	fmt.Println("This is A")
}

func (a *A) route2A() {
	fmt.Println("This is A")
}

type B struct {
	A
}

func (b *B) p() {
	fmt.Println("This is B")
	b.route2A()
}
func TestMutation(t *testing.T) {
	a := A{}
	a.p()
	b := B{}
	b.p()
}
