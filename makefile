build:
	@echo "Building..."
	@cd sample && go build -ldflags="-s -w" -o ../dist