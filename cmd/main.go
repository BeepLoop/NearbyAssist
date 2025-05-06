package main

import (
	"log"
	"nearbyassist/internal/config"
	"nearbyassist/internal/db"
	"nearbyassist/internal/server"
	"nearbyassist/internal/service/cache"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/service/fs"
	"nearbyassist/internal/service/mailer"
	notification_service "nearbyassist/internal/service/notification"
	"nearbyassist/internal/service/route_engine"
	searchhistory "nearbyassist/internal/service/search_history"
	"nearbyassist/internal/service/sse"
	"nearbyassist/internal/service/suggestion_engine"
	"nearbyassist/internal/service/websocket"

	"github.com/go-sql-driver/mysql"
)

func main() {
	// Load configuration file
	cfg := config.GetConfig()

	// Init search history tracker
	searchhistory.New()

	// Init cache
	cache.NewGoCache()

	// Load encryption algorithm
	encrypt := core.NewAES([]byte(cfg.ENCRYPTION_KEY))

	// Load hashing algorithm
	hash := core.NewSha256()

	// Load JWT authenticator
	jwt := core.NewJWTAuthenticator(cfg.JWT_SECRET, cfg.JWT_DURATION)

	// Load file disk
	directories := map[fs.Category]string{
		fs.ID_BACK:               cfg.ID_BACK_DIR,
		fs.ID_FRONT:              cfg.ID_FRONT_DIR,
		fs.FACE:                  cfg.FACE_IMG_DIR,
		fs.APPLICATION_PROOF_DIR: cfg.APPLICATION_PROOF_DIR,
		fs.SERVICE_PHOTO_DIR:     cfg.SERVICE_PHOTO_DIR,
		fs.POLICE_CLEARANCE_DIR:  cfg.POLICE_CLEARANCE_DIR,
		fs.BUG_REPORT_DIR:        cfg.BUG_REPORT_DIR,
		fs.REPORT_USER_DIR:       cfg.REPORT_USER_DIR,
	}
	storage := fs.NewDiskStorage(directories, hash)

	// Load mailer
	// mailer := mailer.NewConsoleMailer()
	mailer := mailer.NewPostmarkMailer(
		config.MustGetEnv("POSTMARK_SERVER_TOKEN"),
		config.MustGetEnv("POSTMARK_ACCOUNT_TOKEN"),
	)

	// Load database configuration
	mysql, err := db.NewMysql(mysql.Config{
		User:                 cfg.DB_USER,
		Passwd:               cfg.DB_PWD,
		Net:                  cfg.DB_NET,
		Addr:                 cfg.DB_HOST + ":" + cfg.DB_PORT,
		DBName:               cfg.DB_NAME,
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer mysql.Close()

	// Initialize sse values and query db for initial data
	sse.New().SetValues(mysql)

	ws := websocket.NewWebsocket()

	notification_service.NewOneSignal(cfg.ONE_SIGNAL_APP_ID, cfg.ONE_SIGNAL_API_KEY)

	serverConfig := server.ServerConfig{
		Config: cfg,

		WS: ws,

		DB: mysql,
		FS: storage,

		Mailer: mailer,

		JWT:     jwt,
		Encrypt: encrypt,
		Hash:    hash,

		RouteEngine:      route_engine.NewOSRM(cfg),
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
