module github.com/e-money/em-ledger

go 1.12

require (
	github.com/cosmos/cosmos-sdk v0.39.1
	github.com/gorilla/mux v1.7.4
	github.com/onsi/ginkgo v1.7.0
	github.com/onsi/gomega v1.4.3
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.33.7
	github.com/tendermint/tm-db v0.5.1
	github.com/tidwall/gjson v1.3.2
	github.com/tidwall/sjson v1.0.4
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
)

// replace github.com/cosmos/cosmos-sdk => ./tmpvendor/cosmos-sdk

// replace github.com/tendermint/tendermint => ./tmpvendor/tendermint
