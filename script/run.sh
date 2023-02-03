NOW_DIR=$(pwd)
PROJ_ROOT=$(cd "$(dirname "$0")/.." || exit 1; pwd)
SERVICE_ROOT="$PROJ_ROOT/service"
GATEWAY_ROOT="$PROJ_ROOT/gateway"

function gateway() {
    cd "$GATEWAY_ROOT" || exit 1
    go run gateway.go -f "etc/gateway.yaml"
}

function run_etcd() {
    cd "$PROJ_ROOT" || exit 1
    etcd
}

function rpc() {
    cd "$NOW_DIR" || exit 1
    goctl rpc protoc *.proto --go_out=./types --go-grpc_out=./types --zrpc_out=. --style=goZero
}

function api() {
    cd "$GATEWAY_ROOT" || exit 1
    goctl api go -api gateway.api -dir . -style goZero
}

cmd=$1
shift

if [ "$cmd" = "gateway" ]; then
    gateway
elif [ "$cmd" = "etcd" ]; then
    run_etcd
elif [ "$cmd" = "rpc" ]; then
    rpc
elif [ "$cmd" = "api" ]; then
    api
else
    cd "$SERVICE_ROOT/$cmd" || exit 1
    MYSQL_DSN="$TIKTOK_MYSQL_DSN" go run "$cmd.go" -f "etc/$cmd.yaml"
fi