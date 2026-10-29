#!/usr/bin/env bash
set -euo pipefail

cd -- "$(dirname -- "${BASH_SOURCE[0]}")"
APP_NAME="ipgeolocation"

compiler=$(go env GOVERSION)
if [[ ! "$compiler" =~ ^go1\.([0-9]+)\.([0-9]+)$ ]] ||
   (( BASH_REMATCH[1] < 26 || (BASH_REMATCH[1] == 26 && BASH_REMATCH[2] < 10) || (BASH_REMATCH[1] == 27 && BASH_REMATCH[2] < 3) )); then
  echo "Use Go 1.26.10, Go 1.27.3, or a newer patched stable version." >&2
  exit 1
fi
command -v govulncheck >/dev/null || { echo "Install govulncheck: go install golang.org/x/vuln/cmd/govulncheck@v1.9.0" >&2; exit 1; }

version_output=$(go run ./cmd/ipgeolocation --version)
VERSION="${version_output##* }"
[[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]] || { echo "Invalid CLI version: $VERSION" >&2; exit 1; }
if [[ "${GITHUB_REF_TYPE:-}" == "tag" && "${GITHUB_REF_NAME:-}" != "v$VERSION" ]]; then
  echo "Release tag must match the CLI version v$VERSION." >&2
  exit 1
fi

go mod verify
go test -race ./...
go vet ./...
govulncheck ./...

rm -rf dist
mkdir -p dist
targets=(
  "linux/amd64"
  "linux/arm64"
  "darwin/amd64"
  "darwin/arm64"
  "windows/amd64"
)

echo "Building $APP_NAME $VERSION with $compiler ..."

for target in "${targets[@]}"; do
  IFS=/ read -r os arch <<< "$target"
  output="${APP_NAME}-${VERSION}-${os}-${arch}"
  if [[ "$os" == "windows" ]]; then
    output="$output.exe"
  fi

  echo "Building and scanning $target ..."
  CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" go build -trimpath -o "dist/$output" "./cmd/$APP_NAME"
  govulncheck -mode=binary "dist/$output"

  if [[ "$os" == "windows" ]]; then
    zip -j "dist/${output%.exe}.zip" "dist/$output"
  else
    COPYFILE_DISABLE=1 tar -czf "dist/$output.tar.gz" -C dist "$output"
  fi
  rm "dist/$output"
done

(
  cd dist
  if command -v sha256sum >/dev/null; then
    sha256sum ./*.tar.gz ./*.zip > SHA256SUMS
    sha256sum -c SHA256SUMS
  else
    shasum -a 256 ./*.tar.gz ./*.zip > SHA256SUMS
    shasum -a 256 -c SHA256SUMS
  fi
)
echo "Archives and SHA256SUMS are in ./dist"

