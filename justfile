default:
    just --list

# Build and execute
dev:
    just build
    just run

# Build binary
build:
    go build -o merlion ./cmd/merlion

# Run binary
run:
    LOG_LEVEL=DEBUG ./merlion

# Remove the Exectutable
clean:
    rm -f ./merlion

# Display the live logs
log:
    tail -f ~/.cache/merlion/merlion.log

export:
    just build
    mv ./merlion ~/tools/
