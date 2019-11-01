package main

import (
	"fmt"
	"testing"
)

func TestGenerateSectionIntSliceOfOrderly(t *testing.T) {
	fmt.Printf("Generate orderly slice:%+v\n", GenerateSectionIntSliceOfOrderly(1, 20, 3))
}

func TestGenerateSectionIntSliceOfDisorderly(t *testing.T) {
	fmt.Printf("Generate disorderly slice:%+v\n", GenerateSectionIntSliceOfDisorderly(15, 20))
}
