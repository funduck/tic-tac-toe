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
	go run cmd/client/main.go --user alice

start-client-bob:
	go run cmd/client/main.go --user bob

build-client:
	go build -o dist/client cmd/client/main.go