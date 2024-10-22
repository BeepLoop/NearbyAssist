package main

import (
	"log"
	"nearbyassist/internal/config"
	"nearbyassist/internal/db"
	message_repo "nearbyassist/internal/repository/message"
	"nearbyassist/internal/server"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/service/email"
	"nearbyassist/internal/service/fs"
	"nearbyassist/internal/service/mailer"
	"nearbyassist/internal/service/route_engine"
	"nearbyassist/internal/service/suggestion_engine"
	"nearbyassist/internal/service/websocket"

	"github.com/go-sql-driver/mysql"
)

func main() {
	// Load configuration file
	config := config.LoadConfig()

	// Load encryption algorithm
	encrypt := auth.NewAES([]byte(config.ENCRYPTION_KEY))

	// Load hashing algorithm
	hash := auth.NewSha256()

	// Load JWT authenticator
	jwt := auth.NewJWTAuthenticator(config.JWT_SECRET, config.JWT_DURATION)

	// Load file disk
	directories := map[fs.Category]string{
		fs.ID_BACK:               config.ID_BACK_DIR,
		fs.ID_FRONT:              config.ID_FRONT_DIR,
		fs.FACE:                  config.FACE_IMG_DIR,
		fs.APPLICATION_PROOF_DIR: config.APPLICATION_PROOF_DIR,
		fs.SERVICE_PHOTO_DIR:     config.SERVICE_PHOTO_DIR,
		fs.SYS_COMPLAINT_DIR:     config.SYS_COMPLAINT_DIR,
		fs.POLICE_CLEARANCE_DIR:  config.POLICE_CLEARANCE_DIR,
	}
	storage := fs.NewDiskStorage(directories, hash)

	// Load database configuration
	mysql, err := db.NewMysql(mysql.Config{
		User:                 config.DB_USER,
		Passwd:               config.DB_PWD,
		Net:                  config.DB_NET,
		Addr:                 config.DB_HOST + ":" + config.DB_PORT,
		DBName:               config.DB_NAME,
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer mysql.Close()

	chatStore := message_repo.NewMysqlChatRepository(mysql)
	ws := websocket.NewWebsocket(chatStore)

	mailService, err := email.NewGoMail(config.EMAIL_FROM, config.EMAIL_PASS)
	if err != nil {
		log.Fatal(err)
	}
	mailer := mailer.New(mailService)

	serverConfig := server.ServerConfig{
		Config: config,

		WS: ws,

		Mailer: mailer,

		DB: mysql,
		FS: storage,

		JWT:     jwt,
		Encrypt: encrypt,
		Hash:    hash,

		RouteEngine:      route_engine.NewOSRM(config),
		SuggestionEngine: suggestion_engine.NewWeightedScoring(),
	}

	// Create and start the server
	server, err := server.NewServer(serverConfig)
	if err != nil {
		log.Fatal(err)
	}

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
