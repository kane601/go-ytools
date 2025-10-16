package tools

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
)

//GenMd5 生成MD5
func GenMd5(data []byte) (string, error) {
	buf := bytes.NewBuffer(data)
	md5g := md5.New()
	if _, err := io.Copy(md5g, buf); err != nil {
		return "", fmt.Errorf("Gen md5 fail,err:%v", err)
	}
	return fmt.Sprintf("%x", md5g.Sum(nil)), nil
}

//GenFileMd5 生成文件md5
func GenFileMd5(path string) (string, error) {
	file, err := OpenReadFile(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	md5g := md5.New()
	if _, err := io.Copy(md5g, file); err != nil {
		return "", fmt.Errorf("Gen md5 fail file:%s,err:%v", path, err)
	}
	return fmt.Sprintf("%x", md5g.Sum(nil)), err
}

//EncodeBase64 生成Base64
func EncodeBase64(data []byte) string {
	encoding := base64.StdEncoding.EncodeToString(data)
	return encoding
}

//DecodeBase64 解码Base64
func DecodeBase64(base string) []byte {
	dncoding, _ := base64.StdEncoding.DecodeString(base)
	return dncoding
}

//GenAesKey ...
func GenAesKey() []byte {
	b := make([]byte, 24)
	_, err := rand.Read(b)
	if err != nil {
		return nil
	}
	return b
}

//GenAesIV ...
func GenAesIV() []byte {
	b := make([]byte, aes.BlockSize)
	_, err := rand.Read(b)
	if err != nil {
		return nil
	}
	return b
}

//AesEncrypt 加密
func AesEncrypt(key, iv, src []byte) []byte {
	ciph, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil
	}

	paddinglen := aes.BlockSize - (len(src) % aes.BlockSize)
	for i := 0; i < paddinglen; i++ {
		src = append(src, byte(paddinglen))
	}
	enbuf := make([]byte, len(src))
	cbce := cipher.NewCBCEncrypter(ciph, []byte(iv))
	cbce.CryptBlocks(enbuf, src)
	return enbuf
}

//AesDecrypt 解密
func AesDecrypt(key, iv, src []byte) []byte {
	ciph, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil
	}

	debuf := make([]byte, len(src))
	cbcd := cipher.NewCBCDecrypter(ciph, []byte(iv))
	cbcd.CryptBlocks(debuf, src)
	paddinglen := int(debuf[len(src)-1])
	return debuf[:len(src)-paddinglen]
}

//GenRSAPair ...
func GenRSAPair() (priv, pub []byte) {
	private, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}
	priv = x509.MarshalPKCS1PrivateKey(private)
	pub, _ = x509.MarshalPKIXPublicKey(&private.PublicKey)
	return
}

//GenSignVerifyPair ...
func GenSignVerifyPair() (pass, iv, pri, pub string) {
	pribin, pubbin := GenRSAPair()
	keybin := GenAesKey()
	ivbin := GenAesIV()
	enpri := AesEncrypt(keybin, ivbin, pribin)
	enpub := AesEncrypt(keybin, ivbin, pubbin)
	pass = EncodeBase64(keybin)
	iv = EncodeBase64(ivbin)
	pri = EncodeBase64(enpri)
	pub = EncodeBase64(enpub)
	return
}

//Hash256 ...
func Hash256(data []byte) []byte {
	hashInstance := crypto.SHA256.New()
	hashInstance.Write(data)
	return hashInstance.Sum(nil)
}

//SignData 签名数据
func SignData(data []byte, pass, iv, pri string) (sig string) {
	passbin := DecodeBase64(pass)
	ivbin := DecodeBase64(iv)
	pribin := DecodeBase64(pri)
	depri := AesDecrypt(passbin, ivbin, pribin)
	priobj, err := x509.ParsePKCS1PrivateKey(depri)
	if err != nil {
		return ""
	}
	sigbin, err := rsa.SignPKCS1v15(rand.Reader, priobj, crypto.SHA256, Hash256(data))
	if err != nil {
		return ""
	}
	return EncodeBase64(sigbin)
}

//VerifyData 验证数据
func VerifyData(data []byte, sig, pass, iv, pub string) bool {
	passbin := DecodeBase64(pass)
	ivbin := DecodeBase64(iv)
	pubbin := DecodeBase64(pub)
	depub := AesDecrypt(passbin, ivbin, pubbin)
	sigbin := DecodeBase64(sig)
	pubinter, err := x509.ParsePKIXPublicKey(depub)
	pubobj := pubinter.(*rsa.PublicKey)
	if err != nil {
		return false
	}
	err = rsa.VerifyPKCS1v15(pubobj, crypto.SHA256, Hash256(data), sigbin)
	if err != nil {
		return false
	}
	return true
}

//VerifyData 验证数据
func IsValidVerifyPair(pass, iv, pri, pub string) bool {
	depri := AesDecrypt(DecodeBase64(pass), DecodeBase64(iv), DecodeBase64(pri))
	_, err := x509.ParsePKCS1PrivateKey(depri)
	if err != nil {
		return false
	}
	depub := AesDecrypt(DecodeBase64(pass), DecodeBase64(iv), DecodeBase64(pub))
	_, err = x509.ParsePKIXPublicKey(depub)
	return err == nil
}
