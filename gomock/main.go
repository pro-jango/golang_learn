package main

import (
	"fmt"
)

func main() {
	fmt.Println("gomock test")
	fmt.Println(addVal(1, 2))
}

type AddInterface interface {
	AddVal(a, b int) (c int)
}

func addVal(a, b int) (c int) {
	return a + b
}
