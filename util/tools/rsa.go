package tools

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

// GenRsaKey 生成 rsa 密钥对
func GenRsaKey(bits int) (pubKey string, priKey string, err error) {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	priBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	priKey = string(pem.EncodeToMemory(priBlock))
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	ret, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return
	}
	pubBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: ret,
	}
	pubKey = string(pem.EncodeToMemory(pubBlock))
	return
}

// LoadPrivateKeyFile 加载RSA私钥
func LoadPrivateKeyFile(priKey string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(priKey))
	if block == nil {
		return nil, errors.New("private key error")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// LoadPublicKeyFile 加载RSA公钥
func LoadPublicKeyFile(pubKey string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubKey))
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pubKeyInterface.(*rsa.PublicKey), nil
}

// PublicEncrypt 使用公钥加密
func PublicEncrypt(pubKey string, data string) (string, error) {
	publicKey, err := LoadPublicKeyFile(pubKey)
	if err != nil {
		return "", err
	}
	buffer := bytes.NewBufferString("")
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(data))
	if err != nil {
		return "", err
	}
	buffer.Write(encrypted)
	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

// PrivateDecrypt 使用私钥解密
func PrivateDecrypt(priKey string, encrypted string) (string, error) {
	privateKey, err := LoadPrivateKeyFile(priKey)
	if err != nil {
		return "", err
	}
	data, err := base64.StdEncoding.DecodeString(encrypted)
	buffer := bytes.NewBufferString("")
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
	if err != nil {
		return "", err
	}
	buffer.Write(decrypted)
	return buffer.String(), err
}
