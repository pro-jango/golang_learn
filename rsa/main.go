package main

import "fmt"

func main() {
	pubkey, prikey, err := GenRsaKeyWithPKCS1(2048)
	fmt.Println(err)
	fmt.Println(pubkey)
	fmt.Println(prikey)
}
