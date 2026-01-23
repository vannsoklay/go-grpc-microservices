#!/usr/bin/env bash
set -e

ROOT_DIR=$(cd "$(dirname "$0")/.." && pwd)
PROTO_DIR="$ROOT_DIR/proto"

# -----------------------------------------
# Ensure Go bin is in PATH (current session)
# -----------------------------------------
export PATH="$HOME/go/bin:$PATH"

# -----------------------------------------
# Detect OS
# -----------------------------------------
OS="$(uname -s)"

install_protoc_hint() {
  case "$OS" in
    Linux)
      if [ -f /etc/arch-release ]; then
        echo "➡ Install protobuf with:"
        echo "   sudo pacman -S protobuf"
      elif [ -f /etc/debian_version ]; then
        echo "➡ Install protobuf with:"
        echo "   sudo apt update && sudo apt install -y protobuf-compiler"
      else
        echo "➡ Install protobuf using your distro package manager"
      fi
      ;;
    Darwin)
      echo "➡ Install protobuf with:"
      echo "   brew install protobuf"
      ;;
    *)
      echo "➡ Please install protobuf manually for your OS"
      ;;
  esac
}

# -----------------------------------------
# Check protoc (system dependency)
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
# Generate protos
# -----------------------------------------
echo "🔧 Generating Auth proto..."
protoc \
  -I="$PROTO_DIR" \
  --go_out="$ROOT_DIR/services/auth-service" \
  --go-grpc_out="$ROOT_DIR/services/auth-service" \
  "$PROTO_DIR/auth/auth.proto"

echo "🔧 Generating User proto..."
protoc \
  -I="$PROTO_DIR" \
  --go_out="$ROOT_DIR/services/user-service" \
  --go-grpc_out="$ROOT_DIR/services/user-service" \
  "$PROTO_DIR/user/user.proto"

echo "🔧 Generating Product proto..."
protoc \
  -I="$PROTO_DIR" \
  --go_out="$ROOT_DIR/services/product-service" \
  --go-grpc_out="$ROOT_DIR/services/product-service" \
  "$PROTO_DIR/product/product.proto"

echo "🔧 Generating Order proto..."
protoc \
  -I="$PROTO_DIR" \
  --go_out="$ROOT_DIR/services/order-service" \
  --go-grpc_out="$ROOT_DIR/services/order-service" \
  "$PROTO_DIR/order/order.proto"

echo "🔧 Generating Payment proto..."
protoc \
  -I="$PROTO_DIR" \
  --go_out="$ROOT_DIR/services/payment-service" \
  --go-grpc_out="$ROOT_DIR/services/payment-service" \
  "$PROTO_DIR/payment/payment.proto"

echo "🔧 Generating Shop proto..."
protoc \
  -I="$PROTO_DIR" \
  --go_out="$ROOT_DIR/services/shop-service" \
  --go-grpc_out="$ROOT_DIR/services/shop-service" \
  "$PROTO_DIR/shop/shop.proto"



echo "✅ Proto generation complete"