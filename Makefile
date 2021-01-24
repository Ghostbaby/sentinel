run:
	go run ./cmd/main.go --config config/config.json

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cmd/sentinel cmd/main.go
	upx cmd/sentinel