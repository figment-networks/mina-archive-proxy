.PHONY: build queries dist

build: queries
	go build

queries:
	sqlembed -path=./queries -package=queries > ./queries/queries.go
	go fmt ./queries/queries.go

dist: queries
	GOOS=linux GOARCH=amd64 go build -o ./dist/mina-archive-proxy