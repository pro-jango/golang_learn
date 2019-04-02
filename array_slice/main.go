package main

import "fmt"

func main() {
	var arr = [6]int{1, 2, 3, 4, 5, 6}
	var slice1 = arr[2:5] //从第2+1个元素开始取到第5个结束，
	fmt.Print("原始数据：", "\t")
	fmt.Println(arr, "\t", slice1, "\t")

	slice1[0] = 100
	fmt.Print("切片赋值：", "\t")
	fmt.Print(arr, "\t", slice1, "\t")
	fmt.Println("<--- 给切片赋值，数组和切片同时修改")

	slice1 = append(slice1, 103)

	fmt.Print("增加元素：", "\t")
	fmt.Print(arr, "\t", slice1, "\t")
	fmt.Println("<--- 切片增加元素103，数组被修改，因为切片指向的数组后面还有一个元素，此时切片新增元素不需要起新的数组")

	for key, value := range slice1 {
		slice1[key] = value + 200
	}
	fmt.Print("循环赋值1：", "\t")
	fmt.Print(arr, "\t", slice1, "\t")
	fmt.Println("<--- 循环赋值后，数组和切片的值都发生改变")

	fmt.Print("函数赋值1：", "\t")
	slicetest(slice1)
	fmt.Print(arr, "\t", slice1, "\t")
	fmt.Println("<--- 函数赋值后，数组和切片的值都发生改变")

	slice1 = append(slice1, 104)
	fmt.Print("增加元素：", "\t")
	fmt.Print(arr, "\t", slice1, "\t")
	fmt.Println("<--- 切片增加元素104，数组未新增元素，切片新增一个元素，此时切片指向了一个新的数组指针")

	for key, value := range slice1 {
		slice1[key] = value + 200
	}
	fmt.Print("循环赋值2：", "\t")
	fmt.Print(arr, "\t", slice1, "\t")
	fmt.Println("<--- 循环赋值后，切片的值发生改变，数组未发生改变，因为切片现在指向的数组已经是新数组了")

	fmt.Print("函数赋值2：", "\t")
	slicetest(slice1)
	fmt.Print(arr, "\t", slice1, "\t")
	fmt.Println("<--- 函数赋值后，切片的值发生改变，数组未发生改变，因为切片现在指向的数组已经是新数组了")

	fmt.Println("总结：")
	fmt.Println("1、不要从数组新建切片，太坑了！")
	fmt.Println("2、切片从数组中新建时，如果元素下表没有超过数组的上限，则不会发生数组copy，切片赋值会影响到数组的值！")
}

func slicetest(slice2 []int) {
	for key, value := range slice2 {
		slice2[key] = value + 200
	}
}
