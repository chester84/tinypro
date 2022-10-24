package weixin

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/chester84/libtools"
	"encoding/base64"
	"fmt"
)

func Decrypt(rawData, key, iv string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(rawData)
	keyB, err1 := base64.StdEncoding.DecodeString(key)
	ivB, _ := base64.StdEncoding.DecodeString(iv)

	if err != nil {
		return "", err
	}
	if err1 != nil {
		return "", err1
	}

	dnData, err := aesCBCDecrypt(data, keyB, ivB)
	if err != nil {
		return "", err
	}
	return string(dnData), nil
}

func aesCBCDecrypt(encryptData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	if len(encryptData) < blockSize {
		err = fmt.Errorf(`ciphertext too short`)
		return nil, err
	}
	if len(encryptData)%blockSize != 0 {
		err = fmt.Errorf(`ciphered text is not a multiple of the block size`)
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encryptData, encryptData)

	// 解填充
	encryptData = libtools.PKCS7UnPadding(encryptData)

	return encryptData, nil
}
