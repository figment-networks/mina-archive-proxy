package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestStakingLedger(t *testing.T) {
	router := gin.Default()
	router.GET("/staking_ledger", handleStakingLedger("echo"))

	t.Run("check CLI presense", func(t *testing.T) {
		handler, err := initStakingLedgerHandler("foobar")

		assert.Equal(t, `exec: "foobar": executable file not found in $PATH`, err.Error())
		assert.Nil(t, handler)
	})

	t.Run("returns current ledger", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/staking_ledger", nil)
		writer := httptest.NewRecorder()
		router.ServeHTTP(writer, req)

		assert.Equal(t, 200, writer.Code)
		assert.Equal(t, "ledger export staking-epoch-ledger", strings.TrimSpace(writer.Body.String()))
	})

	t.Run("checks ledger type", func(t *testing.T) {
		examples := []struct {
			Type   string
			Status int
			Body   string
		}{
			{Type: "current", Status: 200, Body: "ledger export staking-epoch-ledger"},
			{Type: "next", Status: 200, Body: "ledger export next-epoch-ledger"},
			{Type: "staged", Status: 200, Body: "ledger export staged-ledger"},
			{Type: "blah", Status: 400, Body: `{"error":"invalid ledger type"}`},
		}

		for _, ex := range examples {
			req, _ := http.NewRequest("GET", "/staking_ledger", nil)

			q := req.URL.Query()
			q.Add("type", ex.Type)
			req.URL.RawQuery = q.Encode()

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)

			assert.Equal(t, ex.Status, writer.Code)
			assert.Equal(t, ex.Body, strings.TrimSpace(writer.Body.String()))
		}
	})
}
