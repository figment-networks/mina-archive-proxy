package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	version = "0.4.0"
)

var cmdOpts = struct {
	connStr       string
	codaBin       string
	ledgerEnabled bool
	corsEnabled   bool
	showVersion   bool
	debug         bool
}{}

func initFlags() {
	flag.StringVar(&cmdOpts.connStr, "db", "", "Database connection string")
	flag.BoolVar(&cmdOpts.showVersion, "version", false, "Show version")
	flag.StringVar(&cmdOpts.codaBin, "coda-bin", "coda", "Full path to Coda binary")
	flag.BoolVar(&cmdOpts.ledgerEnabled, "ledger-enabled", true, "Enable staking ledger dump endpoint")
	flag.BoolVar(&cmdOpts.corsEnabled, "cors-enabled", false, "Enable CORS on the server")
	flag.BoolVar(&cmdOpts.debug, "debug", false, "Enable debug mode")
	flag.Parse()

	if cmdOpts.connStr == "" {
		cmdOpts.connStr = os.Getenv("DATABASE_URL")
	}

	if !cmdOpts.debug {
		cmdOpts.debug = os.Getenv("DEBUG") == "1"
	}

	gin.SetMode(gin.ReleaseMode)
}

func main() {
	initFlags()

	if cmdOpts.showVersion {
		fmt.Println(version)
		return
	}

	if cmdOpts.debug {
		log.Println("using database connection string:", cmdOpts.connStr)
	}

	log.Println("connecting to database...")
	conn, err := gorm.Open("postgres", cmdOpts.connStr)
	if err != nil {
		log.Fatal("connection failed: ", err)
	}

	conn.DB().SetConnMaxIdleTime(time.Minute * 10)
	conn.DB().SetConnMaxLifetime(time.Minute * 60)

	defer conn.Close()

	// Enable full SQL queries in the logs
	conn.LogMode(cmdOpts.debug)

	router := gin.Default()

	if cmdOpts.corsEnabled {
		log.Println("CORS is enabled")

		router.Use(func(c *gin.Context) {
			c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			c.Header("Access-Control-Expose-Headers", "*")
			c.Header("Access-Control-Allow-Origin", "*")
		})
	}

	router.GET("/", handleInfo(conn))
	router.GET("/status", handleStatus(conn))
	router.GET("/chain", handleChain(conn))
	router.GET("/blocks", handleBlocks(conn))
	router.GET("/blocks/:hash", handleBlock(conn))
	router.GET("/blocks/:hash/user_commands", handleUserCommands(conn))
	router.GET("/blocks/:hash/internal_commands", handleInternalCommands(conn))
	router.GET("/block_producers", handleBlockProducers(conn))
	router.GET("/public_keys", handlePublicKeys(conn))
	router.GET("/public_keys/:id", handlePublicKey(conn))

	if cmdOpts.ledgerEnabled {
		log.Println("staking ledger endpoint is enabled")

		router.GET("/staking_ledger", handleStakingLedger(cmdOpts.codaBin))
	}

	listenAddr := os.Getenv("PORT")
	if listenAddr == "" {
		listenAddr = "3088"
	}
	listenAddr = "0.0.0.0:" + listenAddr

	log.Println("starting server on", listenAddr)
	if err := router.Run(listenAddr); err != nil {
		log.Fatal(err)
	}
}
