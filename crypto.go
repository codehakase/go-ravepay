package ravepay

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"errors"
)

type CryptoService struct {
	Client *Client
}

// TripleDesEncrypt uses the Tripple Digital Encryption Standard to encrypt the transaction data
// len(key) % 8 = 0
func (c *CryptoService) TripleDesEncrypt(origData, key string) (string, error) {
	// call main Encryption method
	return tripleDesEncrypt(origData, []byte(key), PKCS5Padding)
}

// TripleDesEncrypt uses the Tripple Digital Encryption Standard to decrypt encrypted transaction data
// len(key) % 8 = 0
func (c *CryptoService) TripleDesDecrypt(encrypted, key string) (string, error) {
	// call main Encryption method
	return tripleDesDecrypt(encrypted, []byte(key), PKCS5UnPadding)
}

func tripleDesEncrypt(origData string, key []byte, paddingFunc func([]byte, int) []byte) (string, error) {
	iv := []byte(key[:des.BlockSize])
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}
	orig := paddingFunc([]byte(origData), block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(orig))
	blockMode.CryptBlocks(crypted, orig)
	return hex.EncodeToString(crypted), nil
}

func tripleDesDecrypt(encrypted string, key []byte, unPaddingFunc func([]byte) []byte) (string, error) {
	e, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	iv := key[:des.BlockSize]

	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(e))
	blockMode.CryptBlocks(origData, e)
	origData = unPaddingFunc(origData)
	if string(origData) == "unpadding error" {
		return "", errors.New("unpadding error")
	}
	return string(origData), nil
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	if length < unpadding {
		return []byte("unpadding error")
	}
	return origData[:(length - unpadding)]
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	return PKCS5Padding(ciphertext, blockSize)
}

func PKCS7UnPadding(origData []byte) []byte {
	return PKCS5UnPadding(origData)
}
