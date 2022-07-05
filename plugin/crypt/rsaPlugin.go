package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
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
	//fmt.Println("privateKey", privateKey)
	//2、通过X509标准将得到的RAS私钥序列化为：ASN.1 的DER编码字符串

	privateStream := x509.MarshalPKCS1PrivateKey(privateKey)
	//fmt.Println("privateStream", privateStream)

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

// PluginEncryptRSA 对数据进行加密操作
func PluginEncryptRSA(src []byte, path string) (res []byte, err error) {
	//1.获取秘钥（从本地磁盘读取）
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	fileInfo, _ := f.Stat()

	b := make([]byte, fileInfo.Size())
	f.Read(b)

	// 2、将得到的字符串解码
	block, _ := pem.Decode(b)
	// 使用X509将解码之后的数据 解析出来
	//x509.MarshalPKCS1PublicKey(block):解析之后无法用，所以采用以下方法：ParsePKIXPublicKey

	keyInit, err := x509.ParsePKIXPublicKey(block.Bytes) //对应于生成秘钥的x509.MarshalPKIXPublicKey(&publicKey)
	//keyInit1,err:=x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return
	}
	//4.使用公钥加密数据
	pubKey := keyInit.(*rsa.PublicKey)
	res, err = rsa.EncryptPKCS1v15(rand.Reader, pubKey, src)
	return
}

// PluginDecryptRSA 对数据进行解密操作
func PluginDecryptRSA(src []byte, path string) (res []byte, err error) {
	//1.获取秘钥（从本地磁盘读取）
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	fileInfo, _ := f.Stat()
	b := make([]byte, fileInfo.Size())
	f.Read(b)
	block, _ := pem.Decode(b)                                 //解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes) //还原数据
	res, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, src)
	return
}
