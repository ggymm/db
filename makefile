buildSample:
	@echo "Building Sample"
	@cd sample && go build -ldflags="-s -w" -o ../dist