package main

import (
	"fmt"

	"github.com/rouzbehsbz/spenta"
)

func main() {
	arr := []int{3, 4, 5}

	parIter := spenta.SliceParMap(&arr, func(e int) int {
		return e * 2
	})

	parIter.Done()

	fmt.Println(arr)
	println("done ", parIter.Chunks)
}
