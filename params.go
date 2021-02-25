package main

import (
	"github.com/gin-gonic/gin"
)

type blocksParams struct {
	Canonical   bool    `form:"canonical"`
	Creator     *string `form:"creator"`
	StartHeight uint    `form:"start_height"`
	Limit       uint    `form:"limit"`
}

func parseBlockParams(c *gin.Context) *blocksParams {
	params := &blocksParams{}
	if err := c.Bind(params); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return nil
	}

	if params.Limit == 0 {
		params.Limit = 100
	}
	if params.Limit > 1000 {
		params.Limit = 1000
	}

	return params
}
