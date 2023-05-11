# Config

Update config in config/config.go:

1. GeminiSymbols: Array of symbols for which we need to create a market data subscription
2. GUIEnabled: Opens a GUI terminal if set to true.

# Build & Run
### Pre-requisite
Install Go 1.20.3 from https://go.dev/doc/install

### Commands
    go get -v ./...
    go mod tidy
    go test -v ./...
    go build -o app main.go
    ./app
