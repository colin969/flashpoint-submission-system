// @title           Flashpoint Submission API
// @version         1.0
// @description     Yup, it's an API

// @license.name  MIT

// @host      fpfss.unstable.life
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

// Extremely barebones server to demonstrate OAuth 2.0 flow with Discord
// Uses native net/http to be dependency-less and easy to run.
// No sessions logic implemented, re-login needed each visit.
// Edit the config lines a little bit then go build/run it as normal.
package main

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"github.com/FlashpointProject/flashpoint-submission-system/authbot"
	"github.com/FlashpointProject/flashpoint-submission-system/config"
	"github.com/FlashpointProject/flashpoint-submission-system/database"
	"github.com/FlashpointProject/flashpoint-submission-system/logging"
	"github.com/FlashpointProject/flashpoint-submission-system/notificationbot"
	"github.com/FlashpointProject/flashpoint-submission-system/resumableuploadservice"
	"github.com/FlashpointProject/flashpoint-submission-system/transport"
	"github.com/FlashpointProject/flashpoint-submission-system/utils"
)

// testHandler is a sample handler
// @Summary Test API Endpoint
// @Description A simple test endpoint
// @Tags test
// @Accept json
// @Produce json
// @Success 200 {string} string "OK"
// @Router /api/test [get]
func main() {
	if os.Getenv("IN_KUBERNETES") != "" {
		_, err := os.Stat(".env")
		if !os.IsNotExist(err) {
			panic(err)
		}
	} else {
		err := godotenv.Load() // @title           Swagger Example API
		// @version         1.0
		// @description     This is a sample server celler server.
		// @termsOfService  http://swagger.io/terms/

		// @contact.name   API Support
		// @contact.url    http://www.swagger.io/support
		// @contact.email  support@swagger.io

		// @license.name  Apache 2.0
		// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

		// @host      localhost:8080
		// @BasePath  /api/v1

		// @securityDefinitions.basic  BasicAuth

		// @externalDocs.description  OpenAPI
		// @externalDocs.url          https://swagger.io/resources/open-api/
		if err != nil {
			panic(err)
		}
	}
	log := logging.InitLogger()
	l := log.WithField("commit", config.EnvString("GIT_COMMIT")).WithField("runID", utils.NewRealRandomStringProvider().RandomString(8))
	l.Infoln("hi")

	l.Infoln("loading config...")
	conf := config.GetConfig(l)
	l.Infoln("config loaded")

	db := database.OpenDB(l, conf)
	defer db.Close()

	pgdb := database.OpenPostgresDB(l, conf)
	defer pgdb.Close()

	var authBot *discordgo.Session
	var notificationBot *discordgo.Session
	var rsu *resumableuploadservice.ResumableUploadService
	// Skip some extra services when in FP Source Only mode
	if !conf.FlashpointSourceOnlyMode {
		authBot = authbot.ConnectBot(l, conf.AuthBotToken)
		notificationBot = notificationbot.ConnectBot(l, conf.NotificationBotToken)

		l.Infoln("connecting to the resumable upload service")
		var err error
		rsu, err = resumableuploadservice.New(conf.ResumableUploadDirFullPath)
		if err != nil {
			l.Fatal(err)
		}
		defer rsu.Close()
		l.Infoln("resumable upload service connected")
	}

	transport.InitApp(l, conf, db, pgdb, authBot, notificationBot, rsu)
}
