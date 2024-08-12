.PHONY: tailwind-build
tailwind-build:
	tailwindcss -i ./assets/css/input.css -o ./assets/css/tailwind.css
	tailwindcss -i ./assets/css/input.css -o ./assets/css/tailwind.min.css --minify

.PHONY: tailwind-watch
tailwind-watch: 
	tailwindcss -i ./assets/css/input.css -o ./assets/css/tailwind.css --watch

.PHONY: templ-generate
templ-generate:
	templ generate

.PHONY: templ-watch
templ-watch:
	templ generate --watch
	
.PHONY: dev
dev:
	go run ./cmd/seed/seed.go
	make templ-watch &
	make tailwind-watch &	
	go build -o ./tmp/app ./cmd/app/main.go && air

.PHONY: build
build:
	make tailwind-build
	make templ-generate
	go build -ldflags "-X main.Environment=production" -o ./bin/app ./cmd/app/main.go

.PHONY: seed
seed:
	go run ./cmd/seed/seed.go

.PHONY: test
test:
	go test -v  ./tests/...
