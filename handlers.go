package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/figment-networks/mina-archive-proxy/queries"
)

func handleInfo(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		renderObjectQuery(c, conn, queries.Info)
	}
}

func handleChain(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		params := parseBlockParams(c)
		if params == nil {
			return
		}
		renderArrayQuery(c, conn, queries.Chain, params.StartHeight, params.Limit)
	}
}

func handleBlocks(conn *gorm.DB) gin.HandlerFunc {
	tpl, err := template.New("blocks").Parse(queries.Blocks)
	if err != nil {
		log.Fatal("template parse error:", err)
	}

	return func(c *gin.Context) {
		params := parseBlockParams(c)
		if params == nil {
			return
		}

		buff := bytes.NewBuffer(nil)

		if err := tpl.Execute(buff, params); err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		renderArrayQuery(c, conn, buff.String(), params.StartHeight, params.Limit)
	}
}

func handleBlock(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		renderObjectQuery(c, conn, queries.Block, c.Param("hash"))
	}
}

func handleBlockProducers(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		renderRawQuery(c, conn, queries.BlockProducers)
	}
}

func handleUserCommands(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		renderArrayQuery(c, conn, queries.UserCommands, c.Param("hash"))
	}
}

func handleInternalCommands(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		renderArrayQuery(c, conn, queries.InternalCommands, c.Param("hash"))
	}
}

func handlePublicKeys(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		renderRawQuery(c, conn, queries.PublicKeys)
	}
}

func handlePublicKey(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		renderObjectQuery(c, conn, queries.PublicKey, c.Param("id"))
	}
}

func handleStatus(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := conn.Exec("SELECT 1").Error
		if err != nil {
			log.Println("db connection test error:", err)
			c.AbortWithStatusJSON(500, gin.H{"healthy": false})
			return
		}
		c.JSON(200, gin.H{"healthy": true})
	}
}

func handleStakingLedger(codaBinPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ledgerType := c.Query("type")
		if ledgerType == "" {
			ledgerType = "current"
		}
		if !(ledgerType == "current" || ledgerType == "next") {
			c.AbortWithStatusJSON(400, gin.H{"error": "invalid ledger type"})
			return
		}

		ledgerbuf := bytes.NewBuffer(nil)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		cmd := exec.CommandContext(ctx, codaBinPath, "advanced", "dump-staking-ledger", ledgerType, "-json")
		cmd.Stderr = os.Stderr
		cmd.Stdout = ledgerbuf

		if err := cmd.Run(); err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		c.Data(200, "application/json", ledgerbuf.Bytes())
	}
}
