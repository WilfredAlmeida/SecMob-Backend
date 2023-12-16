package crypto

import (
	"io/ioutil"
	"main/utils"
	"os"
)

func LoadAESKey(path string) {
	aesKeyFile, err := os.Open(path)
	if err != nil {
		utils.MLogger.ErrorLog(err.Error())
	}

	//Auto generated to close key file
	defer func(aesKeyFile *os.File) {
		err := aesKeyFile.Close()
		if err != nil {
			utils.MLogger.ErrorLog(err.Error())
		}
	}(aesKeyFile)

	aesKeyBytes, err := ioutil.ReadAll(aesKeyFile)
	if err != nil {
		utils.MLogger.ErrorLog(err.Error())
	}

	utils.AesKey = aesKeyBytes
}
