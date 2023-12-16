package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"main/utils"
	"strings"
)

func DecryptData(data []byte) ([]byte, error) {
	return aesGcmDecrypt(utils.AesKey, data)
}

func aesGcmDecrypt(key []byte, message []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		utils.MLogger.ErrorLog(err.Error())
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		utils.MLogger.ErrorLog(err.Error())
		return nil, err
	}

	chunks := strings.Split(string(message), ".")

	nonce, err := base64.StdEncoding.DecodeString(chunks[0])
	if err != nil {
		utils.MLogger.ErrorLog(err.Error())
		return nil, err
	}

	cipherText, err := base64.StdEncoding.DecodeString(chunks[1])
	if err != nil {
		utils.MLogger.ErrorLog(err.Error())
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		utils.MLogger.ErrorLog(err.Error())
		return nil, err
	}

	return plaintext, nil
}
