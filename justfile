# Build and execute
build-run:
    just build
    just run

# Build binary
build:
    go build -o merlion ./cmd/merlion

# Run binary
run:
    ./merlion

# Remove the Exectutable
clean:
    rm -f ./merlion

# Display the live logs
log:
    tail -f ~/.cache/merlion/merlion.log
