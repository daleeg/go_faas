package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

var (
	PackageName = "rsa"
)

// PluginRSAGenKey 生成rsa公私钥pem文件
func PluginRSAGenKey(bits int, privateKeyName string, publishKeyName string) error {
	/*
		生成私钥
	*/
	//1、使用RSA中的GenerateKey方法生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	fmt.Println("privateKey", privateKey)
	//2、通过X509标准将得到的RAS私钥序列化为：ASN.1 的DER编码字符串

	privateStream := x509.MarshalPKCS1PrivateKey(privateKey)
	fmt.Println("privateStream", privateStream)

	//3、将私钥字符串设置到pem格式块中
	privateBlock := pem.Block{
		Type:  "private key",
		Bytes: privateStream,
	}
	//4、通过pem将设置的数据进行编码"privateKey.pem"，并写入磁盘文件

	fPrivate, err := os.Create(privateKeyName)
	if err != nil {
		return err
	}
	defer fPrivate.Close()

	err = pem.Encode(fPrivate, &privateBlock)
	if err != nil {
		return err
	}
	/*
		生成公钥
		publicKey.pem
	*/

	publicKey := privateKey.PublicKey

	publicStream, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return err
	}
	publicBlock := pem.Block{
		Type:  "public key",
		Bytes: publicStream,
	}
	fPublic, err := os.Create(publishKeyName)
	if err != nil {
		return err
	}
	defer fPublic.Close()
	pem.Encode(fPublic, &publicBlock)

	return nil
}
