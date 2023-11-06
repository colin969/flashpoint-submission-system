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
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if os.Getenv("IN_KUBERNETES") != "" {
		_, err := os.Stat(".env")
		if !os.IsNotExist(err) {
			panic(err)
		}
	} else {
		err := godotenv.Load()
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
		rsu, err := resumableuploadservice.New(conf.ResumableUploadDirFullPath)
		if err != nil {
			l.Fatal(err)
		}
		defer rsu.Close()
		l.Infoln("resumable upload service connected")
	}

	transport.InitApp(l, conf, db, pgdb, authBot, notificationBot, rsu)
}
