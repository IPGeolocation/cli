#!/usr/bin/env bash
set -e

APP_NAME="ipgeolocation"
VERSION="1.0.0"   # fixed version

rm -rf dist
mkdir -p dist

# build targets
targets=(
  "linux amd64"
  "linux arm64"
  "darwin amd64"
  "darwin arm64"
  "windows amd64"
)

echo "🚀 Building $APP_NAME $VERSION ..."

for target in "${targets[@]}"; do
  os=$(echo $target | cut -d' ' -f1)
  arch=$(echo $target | cut -d' ' -f2)

  output="${APP_NAME}-${VERSION}-${os}-${arch}"
  if [ "$os" == "windows" ]; then
    output="$output.exe"
  fi

  echo "👉 Building for $os/$arch ..."
  GOOS=$os GOARCH=$arch go build -o dist/$output ./cmd/$APP_NAME

  # package
  if [ "$os" == "windows" ]; then
    zip -j dist/${output%.exe}.zip dist/$output
    rm dist/$output
  else
    tar -czf dist/$output.tar.gz -C dist $(basename $output)
    rm dist/$output
  fi
done

echo "Done! Binaries are in ./dist"

