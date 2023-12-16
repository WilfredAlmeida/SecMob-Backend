package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"main/utils"
)

func EncryptData(data []byte) ([]byte, error) {
	return aesGcmEncrypt(utils.AesKey, data)
}

func aesGcmEncrypt(key []byte, plaintext []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		utils.MLogger.ErrorLog(err.Error())
		return nil, err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		utils.MLogger.ErrorLog(err.Error())
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		utils.MLogger.ErrorLog(err.Error())
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	a := base64.StdEncoding.EncodeToString(nonce)
	b := base64.StdEncoding.EncodeToString(ciphertext)

	var final string
	final = a + "." + b

	return []byte(final), nil
}
