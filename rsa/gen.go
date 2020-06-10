package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

//RSA公钥私钥产生 PKCS1   公钥私钥 进行base64返回
func GenRsaKeyWithPKCS1(bits int) (pubkey, prikey string, err error) {
	// 生成私钥文件

	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)

	prikey = Base64Encode(string(derStream))
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return
	}
	pubkey = Base64Encode(string(derPkix))
	return
}
