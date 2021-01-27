package main

import (
	"database/sql"

	"github.com/figment-networks/indexing-engine/store/jsonquery"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func renderQuery(ctx *gin.Context, conn *gorm.DB, mode string, query string, args ...interface{}) {
	var result []byte
	var err error

	if mode == "object" {
		result, err = jsonquery.Object(conn, jsonquery.Prepare(query), args)
		if result == nil && err == nil {
			ctx.AbortWithStatusJSON(404, gin.H{"error": "record not found"})
			return
		}
	} else {
		result, err = jsonquery.MustArray(conn, jsonquery.Prepare(query), args)
	}

	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.Data(200, "application/json", result)
}

func renderObjectQuery(ctx *gin.Context, conn *gorm.DB, query string, args ...interface{}) {
	renderQuery(ctx, conn, "object", query, args...)
}

func renderArrayQuery(ctx *gin.Context, conn *gorm.DB, query string, args ...interface{}) {
	renderQuery(ctx, conn, "array", query, args...)
}

func renderRawQuery(ctx *gin.Context, conn *gorm.DB, query string, args ...interface{}) {
	result, err := scanBytes(conn.Raw(query, args...).Rows())
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.Data(200, "application/json", result)
}

func scanBytes(rows *sql.Rows, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []byte
	for rows.Next() {
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}
	}

	return data, nil
}
