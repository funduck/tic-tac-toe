install:
	go mod download
	cd cmd/client && go mod download
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/air-verse/air@latest

tests:
	go test -v ./...
	go vet ./...

test-race:
	go test -race -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

swag-init:
	swag init -g cmd/server/main.go

serve:
	make swag-init
	air

start-server:
	go run cmd/server/main.go

codegen-client:
	./codegen-client.sh

start-client:
	cd cmd/client && go run main.go -user ${USER} -password ${PASSWORD} ${ARGS}

start-server-docker:
	docker run -p 8080:8080 tic-tac-toe-server

build-server:
	go build -o dist/server cmd/server/main.go

build-server-docker:
	docker build -t tic-tac-toe-server -f build/Dockerfile .

build-client:
	cd cmd/client && go build -o ../../dist/client main.go