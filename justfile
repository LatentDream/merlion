default:
    just --list

version := "1.1.0"
commit := `git rev-parse HEAD`
build_time := `date '+%Y-%m-%d %H:%M:%S'`
module_path := `go list -m`

# Build and execute
dev:
    just build
    just run

# Build binary
build:
    go build \
        -ldflags "-X {{module_path}}/cmd/merlion/version.Version=\"{{version}}\" -X {{module_path}}/cmd/merlion/version.Commit=\"{{commit}}\"" \
        -o merlion \
        ./cmd/merlion


# Run binary
run:
    EDITOR=vim LOG_LEVEL=DEBUG APP_ENV=dev MERLION_DB_PATH=./dev.db MERLION_PATH="~/host/Documents/notes/test/Test/" ./merlion

# Remove the Exectutable
clean:
    rm -f ./merlion

# Display the live logs
log:
    tail -f $(go run os_check.go)

# Export the binary in the user local bin
export:
    just build
    mv ./merlion ~/.local/bin/
