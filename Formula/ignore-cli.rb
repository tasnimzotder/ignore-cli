class IgnoreCli < Formula
    desc "A simple command-line tool for managing .gitignore files in your Git repositories Topics"
    homepage "https://github.com/tasnimzotder/ignore-cli"
    url "https://github.com/tasnimzotder/ignore-cli/releases/download/v0.0.1/ignore-darwin-arm64-v0.0.1"
    sha256 "84ed217248ba9a49ee9a966e67f0beb8d41264770a08ea40fd953c523e4acfb0"
    version "v0.0.1"

    def install
        bin.install "ignore-darwin-arm64-v0.0.1" => "ignore"
    end

    test do
        system "#{bin}/ignore", "--help"
    end
end
