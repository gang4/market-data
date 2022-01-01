package algorithms

import (
	"fmt"
	"testing"
)

func TestGetDta(t *testing.T) {
	var prices []float64 = []float64{7, 1, 5, 3, 6, 4}
	fmt.Println(MaxProfit(prices))
}

func TestGetDta1(t *testing.T) {
	var prices []float64 = []float64{7, 6, 4, 3, 1}
	fmt.Println(MaxProfit(prices))
}

func TestGetDta2(t *testing.T) {
	var prices []float64 = []float64{1, 2}
	fmt.Println(MaxProfit(prices))
}
