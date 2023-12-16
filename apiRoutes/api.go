package routes

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"main/crypto"
	"main/utils"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type CommandsJsonObj struct {
	CommandsJsonArray []CommandSchema `json:"commands"`
}

type CommandSchema struct {
	Title   string `json:"title"`
	Command string `json:"command"`
	Id      string `json:"id"`
}

type RunCommandRequestSchema struct {
	CommandId string `json:"commandId" binding:"min=6,max=6" validate:"required"` //To change the length, adjust values here
}

type CommandResponseSchema struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type RunCommandResponseSchema struct {
	Op string `json:"op"`
}

type RunCommandResponseObj struct {
	CommandResult []RunCommandResponseSchema `json:"commandResult"`
}

func encryptResponse(jsonData gin.H) string {
	jsonStringBytes, _ := json.Marshal(jsonData)
	encryptedJson, _ := crypto.EncryptData(jsonStringBytes)
	encryptedJsonBase64 := base64.StdEncoding.EncodeToString(encryptedJson)
	return encryptedJsonBase64
}

func RunCMD(path string, args []string, debug bool) (out string, err error) {

	cmd := exec.Command(path, args...)

	var b []byte
	b, err = cmd.CombinedOutput()
	out = string(b)

	if debug {
		utils.MLogger.InfoLog(strings.Join(cmd.Args[:], " "))
		if err != nil {
			utils.MLogger.ErrorLog(err.Error())
		}
	}

	return
}

func RunCommandApiHandler(c *gin.Context) {

	rawData, _ := c.GetRawData()

	rawDataQuotesTrimmed := strings.Trim(string(rawData), "\"")

	receivedData, err := base64.StdEncoding.DecodeString(rawDataQuotesTrimmed)
	if err != nil {
		utils.MLogger.ErrorLog(err.Error())

		c.IndentedJSON(http.StatusBadRequest, encryptResponse(gin.H{
			"status":  0,
			"message": "Invalid Body",
			"payload": nil,
		}))

		utils.MLogger.InfoLog("/runCommand served with error 'Invalid Body'")
		return
	}

	decryptedBody, err := crypto.DecryptData(receivedData)
	if err != nil {
		utils.MLogger.ErrorLog(err.Error())
		//Sending response
		c.IndentedJSON(http.StatusBadRequest, encryptResponse(gin.H{
			"status":  0,
			"message": "Invalid Body",
			"payload": nil,
		}))

		utils.MLogger.InfoLog("/runCommand served with error 'Invalid Body'")
		return
	}

	var body RunCommandRequestSchema
	err = json.Unmarshal(decryptedBody, &body)
	if err != nil {
		utils.MLogger.ErrorLog(err.Error())

		c.IndentedJSON(http.StatusBadRequest, encryptResponse(gin.H{
			"status":  0,
			"message": "Invalid Body",
			"payload": nil,
		}))

		utils.MLogger.InfoLog("/runCommand served with error 'Invalid Body'")

		return
	}

	cmdObj := getCommandById(body.CommandId)
	if cmdObj == nil {
		utils.MLogger.ErrorLog("Invalid Command ID in request")
		return
	}

	var args []string

	cmdBreakdown := strings.Split(cmdObj.Command, " ")
	commandToRun := cmdBreakdown[0]
	args = cmdBreakdown[1:]

	output, mError := RunCMD(commandToRun, args, true)
	if mError != nil {
		utils.MLogger.ErrorLog(mError.Error())

		c.IndentedJSON(http.StatusInternalServerError, encryptResponse(gin.H{
			"status":  0,
			"message": mError.Error(),
			"payload": nil,
		}))
		utils.MLogger.InfoLog("/runCommand served with error")

		return
	}

	c.IndentedJSON(http.StatusOK, encryptResponse(gin.H{
		"status":  1,
		"message": "Command Executed Successfully",
		"payload": []gin.H{
			{
				"output": output,
			},
		},
	}))
	utils.MLogger.InfoLog("/runCommand served with output")

}

var commandsJsonParsedObj CommandsJsonObj

func LoadCommandsFromFile() {
	jsonFile, err := os.Open("./commands/commands.json")
	// if os.Open returns an error then handle it
	if err != nil {
		utils.MLogger.ErrorLog(err.Error())
	}

	utils.MLogger.InfoLog("Successfully Opened commandsJsonParsedObj.json")

	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {

		}
	}(jsonFile)

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &commandsJsonParsedObj)
	if err != nil {
		utils.MLogger.ErrorLog(err.Error())
		return
	}

	for i := 0; i < len(commandsJsonParsedObj.CommandsJsonArray); i++ {
		commandsJsonParsedObj.CommandsJsonArray[i].Id = utils.GenerateRandomString(6)
	}

}

func GetCommandsApiHandler(c *gin.Context) {

	var res []CommandResponseSchema

	for i := 0; i < len(commandsJsonParsedObj.CommandsJsonArray); i++ {
		cmd := commandsJsonParsedObj.CommandsJsonArray[i]
		res = append(res, CommandResponseSchema{
			Id:    cmd.Id,
			Title: cmd.Title,
		})
	}

	c.IndentedJSON(http.StatusOK, encryptResponse(gin.H{
		"status":  1,
		"message": "Data Fetched Successfully",
		"payload": res,
	}))

	utils.MLogger.InfoLog("/getCommands served with output")
}

func getCommandById(id string) *CommandSchema {
	for i := 0; i < len(commandsJsonParsedObj.CommandsJsonArray); i++ {
		if commandsJsonParsedObj.CommandsJsonArray[i].Id == id {
			return &commandsJsonParsedObj.CommandsJsonArray[i]
		}
	}
	return nil
}
