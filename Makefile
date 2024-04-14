.PHONY: tailwind-build
tailwind-build:
	nix-shell -p tailwindcss --run "tailwindcss -i ./assets/css/input.css -o ./assets/css/style.css"
	nix-shell -p tailwindcss --run "tailwindcss -i ./assets/css/input.css -o ./assets/css/style.min.css --minify"

.PHONY: templ-generate
templ-generate:
	nix run github:a-h/templ generate

.PHONY: templ-watch
templ-watch:
	nix run github:a-h/templ generate --watch
	
.PHONY: dev
dev:
	go build -o ./tmp/$(APP_NAME) ./$(APP_NAME)/main.go && air

.PHONY: build
build:
	make tailwind-build
	make templ-generate
	go build -ldflags "-X main.Environment=production" -o ./bin/$(APP_NAME) ./$(APP_NAME)/main.go

