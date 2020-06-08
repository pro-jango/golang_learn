package main

import "fmt"

func main() {
	str1 := "iaodiudfk@i0skw.com"
	str2 := "撒sss去问"
	str3 := "撒收到"

	str4 := "撒收到s"
	str5 := "撒3"
	fmt.Println(formatNickName(str1))
	fmt.Println(formatNickName(str2))
	fmt.Println(formatNickName(str3))
	fmt.Println(formatNickName(str4))
	fmt.Println(formatNickName(str5))
}

func formatNickName(str string) string {
	newStr := ""
	str2 := []rune(str)
	lenStr2 := len(str2)
	if lenStr2 >= 10 {
		newStr = string(str2[:4]) + "**" + string(str2[len(str2)-3:])
	} else if lenStr2 >= 6 {
		newStr = string(str2[:3]) + "**" + string(str2[len(str2)-2:])
	} else if lenStr2 >= 4 {
		newStr = string(str2[:2]) + "**" + string(str2[len(str2)-1:])
	} else if lenStr2 > 2 {
		newStr = string(str2[:2]) + "**"
	} else {
		newStr = string(str2)
	}
	return newStr
}
