package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"net/url"
	"sort"
)

//GetRSASign RSA
//pirkey base64
//return hexstring
func GetRSASign(data map[string]string, prikey string) (string, error) {
	content := BuildSignStr(data)
	val, err := RsaSignPKCS1v15WithSHA256(prikey, []byte(content))
	if err != nil {
		return "", err
	}
	signStr := HexEncodeStr(string(val))
	return signStr, nil
}

//CheckRSASign RSA
//return error
func CheckRSASign(strData, signStr, pubKey string) error {
	var err error
	signData, err := hex.DecodeString(signStr)
	if err != nil {
		return err
	}
	originalData := []byte(strData)
	err = RsaVerifySignPKCS1v15WithSHA256(originalData, signData, pubKey)
	return err
}

//验签：对采用sha256算法
//RsaVerifySignPKCS1v15WithSHA256
func RsaVerifySignPKCS1v15WithSHA256(originalData, signData []byte, pubKey string) error {

	pb, _ := Base64Decode(pubKey)
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey([]byte(pb))
	if err != nil {
		return err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	if err != nil {
		return err
	}
	hash := sha256.New()
	hash.Write(originalData)
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash.Sum(nil), signData)
}

// 签名 对采用sha256算法
//RsaSignPKCS1v15WithSHA256
func RsaSignPKCS1v15WithSHA256(privateKey string, data []byte) ([]byte, error) {

	pb, _ := Base64Decode(privateKey)
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey([]byte(pb))
	if err != nil {
		return nil, err
	}
	hash := sha256.New()
	hash.Write(data)
	return rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hash.Sum(nil))
}

//BuildSignStr BuildSignStr
func BuildSignStr(data map[string]string) string {
	andSeparator := ""
	content := ""
	var keys []string

	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {

		value := data[key]
		content += andSeparator + key + "=" + value
		andSeparator = "&"
	}
	return content
}

func Base64Encode(str string) string {
	src := []byte(str)
	return string([]byte(base64.StdEncoding.EncodeToString(src)))
}

func Base64Decode(str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(str)
}

func HexEncodeStr(s string) string {
	return hex.EncodeToString([]byte(s))
}

//QueryStringToMap convert url query string to a map structure
func QueryStringToMap(query string) (ret map[string]string, err error) {

	ret = make(map[string]string)
	m, err := url.ParseQuery(query)

	if err != nil {
		return ret, err
	}

	for k, v := range m {
		if len(v) > 0 {
			ret[k] = v[0]
		}
	}

	return ret, nil
}
