.PHONY: build queries dist setup swagger

build: queries
	go build

test:
	go test -cover -race

queries:
	sqlembed -path=./queries -package=queries > ./queries/queries.go
	go fmt ./queries/queries.go

setup:
	go get -u github.com/sosedoff/sqlembed

dist: queries
	GOOS=linux GOARCH=amd64 go build -o ./dist/mina-archive-proxy

swagger:
	swagger serve ./swagger.yml