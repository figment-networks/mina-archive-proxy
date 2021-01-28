package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	version = "0.2.0"
)

var cmdOpts = struct {
	codaBin       string
	ledgerEnabled bool
	showVersion   bool
}{}

func init() {
	flag.BoolVar(&cmdOpts.showVersion, "version", false, "Show version")
	flag.StringVar(&cmdOpts.codaBin, "coda-bin", "", "Full path to Coda binary")
	flag.BoolVar(&cmdOpts.ledgerEnabled, "ledger-enabled", true, "Enable staking ledger dump endpoint")
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)
}

func main() {
	if cmdOpts.showVersion {
		fmt.Println(version)
		return
	}

	log.Println("connecting to database...")
	conn, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Enable full SQL queries in the logs
	conn.LogMode(os.Getenv("TRACE_SQL") == "1")

	router := gin.Default()
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
