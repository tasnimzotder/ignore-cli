class IgnoreCli < Formula
    desc "A simple command-line tool for managing .gitignore files in your Git repositories Topics"
    homepage "https://github.com/tasnimzotder/ignore-cli"
    url "https://github.com/tasnimzotder/ignore-cli/releases/download/v0.0.5/ignore-darwin-arm64-v0.0.5"
    sha256 "ac8874ffdca9a4ac93574242fd8505eede3759fc8d500cf3145e9332bd59b9a7"
    version "v0.0.5"

    def install
        bin.install "ignore-darwin-arm64-v0.0.5" => "ignore"
    end

    test do
        system "#{bin}/ignore", "--help"
    end
end
