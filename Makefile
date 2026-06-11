tests:
	go test -v ./...
	go vet ./...

swag-init:
	swag init -g cmd/server/main.go

serve:
	make swag-init
	air

start-server:
	go run cmd/server/main.go

codegen-client:
	./codegen-client.sh

start-client-alice:
	cd cmd/client && go run main.go -user alice -password 123456

start-client-bob:
	cd cmd/client && go run main.go -user bob -password qwerty

start-server-docker:
	docker run -p 8080:8080 tic-tac-toe-server

build-server:
	go build -o dist/server cmd/server/main.go

build-server-docker:
	docker build -t tic-tac-toe-server -f build/Dockerfile .

build-client:
	cd cmd/client && go build -o ../../dist/client main.go