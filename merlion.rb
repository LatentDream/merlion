class Merlion < Formula
  desc "TUI for note-taking app with Obsidian Vault support"
  homepage "https://note.merlion.dev"
  url "https://github.com/latentDream/merlion/archive/refs/tags/1.2.0.tar.gz"
  sha256 "529d7b6ee2652550adc538bb8054dc593863c97f95531e3adb957e780df0048e"
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
