class IgnoreCli < Formula
    desc "A simple command-line tool for managing .gitignore files in your Git repositories Topics"
    homepage "https://github.com/tasnimzotder/ignore-cli"
    url "https://github.com/tasnimzotder/ignore-cli/releases/download/v0.0.2/ignore-darwin-arm64-v0.0.2"
    sha256 "f1711a535d618828f2707049ab4731fde1a788ebedb8ca6e659a100f39a3dd0e"
    version "0.0.2"

    def install
        bin.install "ignore"
    end

    test do
        system "#{bin}/ignore", "--help"
    end
end
