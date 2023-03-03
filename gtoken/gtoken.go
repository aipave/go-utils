package gtoken

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	deviceType     = "1"
	rpcAccessToken = "1"
	fixedBits      = 16
)

var (
	key = []byte("\x2e\xe2\x53\xba\x33\x14\x59\x48\xa0\xa4\x4e\x3c\x56\x3c\xa7\xb6")
	iv  = []byte("\x38\x37\xf5\x98\x84\xf7\x41\x0c\x2f\x05\xa3\x15\x79\x86\xd1\x5a")
)

// CreateToken
func CreateToken(uusrid uint64) (string, error) {
	expireSeconds := time.Now().Add(time.Hour * 24 * 2).Unix()
	plainText := fmt.Sprintf("%v:%s:%s:%s", uusrid, deviceType, strconv.FormatInt(expireSeconds, 10), rpcAccessToken)
	return encrypt([]byte(plainText))
}

func encrypt(rawData []byte) (string, error) {
	data, err := aesCBCEncrypt(rawData)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(data), nil
}

func aesCBCEncrypt(rawData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	rawData = pkcs5Padding(rawData, 16)
	cipherText := make([]byte, len(rawData))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, rawData)
	return cipherText, nil
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

// VerifyToken
func VerifyToken(usrid uint64, ticket string) error {
	ticketUid, expiredTime, err := decryptTicket(ticket)
	if err != nil {
		return err
	}

	if usrid != ticketUid {
		err := fmt.Errorf("invalid usrid")
		return err
	}

	if expiredTime < time.Now().Unix() {
		err := fmt.Errorf("token expire")
		return err
	}

	return nil
}

// DecryptTicket
func decryptTicket(encryptedStr string) (usrid uint64, expire int64, err error) {
	plainText, err := decrypt(encryptedStr)
	if err != nil {
		return 0, -1, err
	}
	arrays := strings.Split(plainText, ":")
	if len(arrays) != 4 {
		return 0, -1, fmt.Errorf("invalid plainText(%v)", plainText)
	}

	usrid, err = strconv.ParseUint(arrays[0], 10, 64)
	if err != nil {
		return 0, -1, fmt.Errorf("invalid usrid from plainText(%v)", plainText)
	}

	expire, err = strconv.ParseInt(arrays[2], 10, 64)
	if err != nil {
		return 0, -1, fmt.Errorf("invalid expiredtime from plainText(%v)", plainText)
	}
	return usrid, expire, nil
}

func decrypt(encryptedStr string) (string, error) {
	data, err := hex.DecodeString(encryptedStr)
	if err != nil {
		return "", err
	}
	dnData, err := aesCBCDecrypt(data)
	if err != nil {
		return "", err
	}
	return string(dnData), nil
}

func aesCBCDecrypt(encryptData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(encryptData)%fixedBits != 0 {
		err := fmt.Errorf("ciphertext is not a multiple of the block size")
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encryptData, encryptData)
	encryptData = pKCS7UnPadding(encryptData)
	return encryptData, nil
}

func pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	if length <= 0 {
		return nil
	}

	unPadding := int(origData[length-1])
	if length < unPadding {
		return nil
	}

	return origData[:(length - unPadding)]
}
