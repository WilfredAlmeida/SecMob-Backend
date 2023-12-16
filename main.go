package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	routes "main/apiRoutes"
	"main/crypto"
	"main/utils"
	"net/http"
	"os"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		utils.MLogger.ErrorLog(err.Error())
	}
}

func main() {

	utils.MLogger = utils.AggregatedLogger{
		InfoLogger:  log.New(os.Stdout, "INFO ", utils.LoggerFlags),
		WarnLogger:  log.New(os.Stdout, "WARN ", utils.LoggerFlags),
		ErrorLogger: log.New(os.Stdout, "ERROR ", utils.LoggerFlags),
	}

	logFile, err := os.Create("secureaccess-logs.log")
	if err != nil {
		utils.MLogger.ErrorLog(err.Error())
	}

	if logFile != nil {
		utils.MLogger.InfoLogger.SetOutput(logFile)
		utils.MLogger.WarnLogger.SetOutput(logFile)
		utils.MLogger.ErrorLogger.SetOutput(logFile)
	}

	PORT := os.Getenv("PORT")
	GIN_MODE := os.Getenv("GIN_MODE")

	if GIN_MODE == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(CORSMiddleware())
	err = router.SetTrustedProxies(nil)
	if err != nil {
		utils.MLogger.ErrorLog(err.Error())
		return
	}

	router.GET("/", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "SecureAccess Server Running",
			"payload": nil,
		})
		utils.MLogger.InfoLog("/ served")
	})

	router.POST("/getCommands", routes.GetCommandsApiHandler)

	router.POST("/runCommand", routes.RunCommandApiHandler)

	routes.LoadCommandsFromFile()

	go utils.WatchFile("./commands/commands.json", routes.LoadCommandsFromFile)

	crypto.LoadAESKey("./keys/aesKey")

	err = router.Run(":" + PORT)
	if err != nil {
		utils.MLogger.ErrorLog(err.Error())
		return
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
