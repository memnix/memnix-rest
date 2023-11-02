export GOPATH="$(go env GOPATH)"
export PATH="${PATH}:${GOPATH}/bin"

# Install wire
go install github.com/google/wire/cmd/wire@latest

# Install gofumpt, fieldalignment, gci
go install mvdan.cc/gofumpt@latest
go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
go install github.com/daixiang0/gci@latest

# Install golangci-lint
wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.55.1
