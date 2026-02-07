default:
    just --list

version := "1.2.1"
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

# Release the latest version to GitHub
release:
    goreleaser release

# Update Homebrew formula for the latest release
release-brew VERSION:
    #!/usr/bin/env bash
    set -euo pipefail
    
    echo "Updating Homebrew formula for version {{VERSION}}..."
    
    # Download the tarball and calculate SHA256
    URL="https://github.com/latentDream/merlion/archive/refs/tags/{{VERSION}}.tar.gz"
    echo "Downloading $URL..."
    SHA256=$(curl -sL "$URL" | shasum -a 256 | awk '{print $1}')
    echo "SHA256: $SHA256"
    
    # Update the formula
    cat > merlion.rb << EOF
    class Merlion < Formula
      desc "TUI for note-taking app with Obsidian Vault support"
      homepage "https://note.merlion.dev"
      url "https://github.com/latentDream/merlion/archive/refs/tags/{{VERSION}}.tar.gz"
      sha256 "$SHA256"
      license "MIT"
    
      depends_on "go" => :build
    
      def install
        ldflags = %W[
          -s -w
          -X merlion/cmd/merlion/version.Version=#{version}
          -X merlion/cmd/merlion/version.Commit=#{tap.user}
        ]
        system "go", "build", *std_go_args(ldflags:), "./cmd/merlion"
      end
    
      test do
        assert_match "version: #{version}", shell_output("#{bin}/merlion version 2>&1")
      end
    end
    EOF
    
    echo "Formula updated in merlion.rb"
    echo ""
    echo "To update your tap, run:"
    echo "  cp merlion.rb \"\$(brew --repository)/Library/Taps/latentdream/homebrew-merlion/Formula/merlion.rb\""
    echo "  cd \"\$(brew --repository)/Library/Taps/latentdream/homebrew-merlion\""
    echo "  git add Formula/merlion.rb"
    echo "  git commit -m \"Update merlion to {{VERSION}}\""
    echo "  git push"

# Export the binary in the user local bin
export:
    just build
    mv ./merlion ~/.local/bin/


