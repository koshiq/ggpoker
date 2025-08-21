build:
	@go build -o bin/ggpoker

run: build
	@./bin/ggpoker

test:
	go test -v ./...

# Web frontend tasks
web-install:
	cd web && npm install

web-build:
	cd web && npm run build

web-dev:
	cd web && npm run dev
