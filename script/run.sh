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

cmd=$1
shift

if [ "$cmd" = "gateway" ]; then
    gateway
elif [ "$cmd" = "etcd" ]; then
    run_etcd
else
    cd "$SERVICE_ROOT/$cmd" || exit 1
    MYSQL_DSN="$TIKTOK_MYSQL_DSN" go run "$cmd.go" -f "etc/$cmd.yaml"
fi