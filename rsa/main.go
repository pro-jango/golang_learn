package main

import (
	"encoding/hex"
	"flag"
	"fmt"
)

const prikey = ""
const pubkey = ""

func main() {
	run := flag.String("run", "", "run function")
	dataStr := flag.String("dataStr", "", "queryString 原始字符串")
	signStr := flag.String("signStr", "", "signString 加密字符串")
	flag.Parse()
	fmt.Printf("dataStr=%s\n", *dataStr)
	fmt.Printf("signStr=%s\n", *signStr)
	fmt.Println("-----------程序执行结果------------")
	if *run == "keygen" {
		keygen()
	} else if *run == "sign" {
		sign(*dataStr)
	} else if *run == "check" {
		check(*dataStr, *signStr)
	} else {
		flag.Usage()
	}
}

//keygen获得公钥和私钥
func keygen() {
	pubkey, prikey, _ := GenRsaKeyWithPKCS1(1024)
	fmt.Println("pubkey:", pubkey)
	fmt.Println("prikey:", prikey)
}

func sign(query string) {
	var data = make(map[string]string)

	data, err := QueryStringToMap(query)
	if err != nil {
		fmt.Println(fmt.Sprintf("QueryStringToMap-error:%v", err))
		return
	}
	fmt.Printf("data:%v\n", data)
	originSign, err := GetRSASign(data, prikey)
	if err != nil {
		fmt.Println(fmt.Sprintf("GetRSASign-error:%v", err))
		return
	}
	fmt.Printf("加密串:%s\n", originSign)
}

func check(queryString, signString string) {
	signData, err := hex.DecodeString(signString)
	if err != nil {
		fmt.Println(fmt.Sprintf("DecodeString error:%v", err))
		return
	}

	err = RsaVerifySignPKCS1v15WithSHA256([]byte(queryString), signData, pubkey)
	if err != nil {
		fmt.Println(fmt.Sprintf("验证失败 error:%v", err))
	} else {
		fmt.Println("rsa验证通过")
	}
}
