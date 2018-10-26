DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
GOPATH=$(go env GOPATH):${DIR}
go run main.go
