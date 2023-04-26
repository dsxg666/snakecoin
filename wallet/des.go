package wallet

import (
	"bytes"
	"crypto/des"
	"encoding/hex"
	"errors"
)

func DesEncrypt(text, key []byte) (string, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	bs := block.BlockSize()
	text = ZeroPadding(text, bs)
	if len(text)%bs != 0 {
		return "", errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(text))
	dst := out
	for len(text) > 0 {
		block.Encrypt(dst, text[:bs])
		text = text[bs:]
		dst = dst[bs:]
	}
	return hex.EncodeToString(out), nil
}

func DesDecrypt(decrypted string, key []byte) ([]byte, error) {
	src, err := hex.DecodeString(decrypted)
	if err != nil {
		return nil, err
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(src))
	dst := out
	bs := block.BlockSize()
	if len(src)%bs != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	for len(src) > 0 {
		block.Decrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	out = ZeroUnPadding(out)
	return out, nil
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

func TrimKey(key []byte) []byte {
	l := len(key)
	if len(key) < 8 {
		for i := 0; i < 8-l; i++ {
			key = append(key, byte(';'))
		}
		return key
	} else if len(key) > 8 {
		return key[:8]
	}
	return key
}
