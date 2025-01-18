
# Build and execute
build-run:
    just build
    just run

# Build binary
build:
    go build -o merlion ./cmd/merlion

# Run binary
run:
    ./merlion 2> app.log

clean:
    rm -f ./merlion

log:
	tail -f ~/.cache/merlion/merlion.log
