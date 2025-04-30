build:
	go install github.com/gone-io/gonectl@latest
	go install go.uber.org/mock/mockgen@latest
	gonectl build -ldflags="-w -s" -tags musl -o bin/server ./cmd/server

run:
	go install github.com/gone-io/gonectl@latest
	go install go.uber.org/mock/mockgen@latest
	gonectl run ./cmd/server


build-docker:
	docker compose build