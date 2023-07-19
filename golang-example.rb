# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class GolangExample < Formula
  desc ""
  homepage ""
  version "0.2.0"
  depends_on :macos

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/attachmentgenie/golang-example/releases/download/v0.2.0/golang-example_0.2.0_darwin_arm64.tar.gz"
      sha256 "f40d4fd9586987620fad118a9f1f2f5d438a160821998afc62d02fec9982dde8"

      def install
        bin.install "golang-example"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/attachmentgenie/golang-example/releases/download/v0.2.0/golang-example_0.2.0_darwin_amd64.tar.gz"
      sha256 "f2140e360e2e303ac9696d985eafb63e1b54f64bd793f84489aba14de7c0ec18"

      def install
        bin.install "golang-example"
      end
    end
  end
end
