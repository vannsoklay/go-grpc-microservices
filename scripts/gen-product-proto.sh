#!/usr/bin/env bash
set -e

ROOT_DIR=$(cd "$(dirname "$0")/.." && pwd)
PROTO_DIR="$ROOT_DIR/proto"

# -----------------------------------------
# Ensure Go bin is in PATH
# -----------------------------------------
export PATH="$HOME/go/bin:$PATH"

# -----------------------------------------
# Detect OS for hint messaging
# -----------------------------------------
OS="$(uname -s)"

install_protoc_hint() {
  case "$OS" in
    Linux)
      if [ -f /etc/arch-release ]; then
        echo "➡ Install protobuf with: sudo pacman -S protobuf"
      elif [ -f /etc/debian_version ]; then
        echo "➡ Install protobuf with: sudo apt update && sudo apt install -y protobuf-compiler"
      fi
      ;;
    Darwin)
      echo "➡ Install protobuf with: brew install protobuf"
      ;;
    *)
      echo "➡ Please install protobuf manually for your OS"
      ;;
  esac
}

# -----------------------------------------
# Check system dependencies
# -----------------------------------------
if ! command -v protoc >/dev/null 2>&1; then
  echo "❌ protoc not found"
  install_protoc_hint
  exit 1
fi

# -----------------------------------------
# Install protoc-gen-go if missing
# -----------------------------------------
if ! command -v protoc-gen-go >/dev/null 2>&1; then
  echo "⬇ Installing protoc-gen-go..."
    
fi

# -----------------------------------------
# Install protoc-gen-go-grpc if missing
# -----------------------------------------
if ! command -v protoc-gen-go-grpc >/dev/null 2>&1; then
  echo "⬇ Installing protoc-gen-go-grpc..."
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# -----------------------------------------
# Generate Code
# -----------------------------------------

echo "🔧 Generating Product, Variant, and Tag protos..."
# We include all three files together because they are in the same package 
# and reference each other.
protoc -I="$PROTO_DIR" \
  --go_out="$ROOT_DIR/services/product-service" \
  --go-grpc_out="$ROOT_DIR/services/product-service" \
   "$PROTO_DIR/product/tag.dev.proto"
  # "$PROTO_DIR/product/product.dev.proto"
  # "$PROTO_DIR/product/variant.dev.proto" \
  # "$PROTO_DIR/product/category.dev.proto" \
  # "$PROTO_DIR/product/tag.dev.proto"

echo "✅ Proto product generation complete"