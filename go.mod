module github.com/MinterTeam/minter-node-cli

go 1.13

require (
	github.com/MinterTeam/minter-go-node v1.0.5-0.20191113110340-a46b8ef88084
	github.com/golang/protobuf v1.3.2
	github.com/tendermint/tendermint v0.32.6
	github.com/urfave/cli/v2 v2.0.0-alpha.2
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	google.golang.org/grpc v1.25.0
)

replace github.com/MinterTeam/minter-go-node v1.0.4 => github.com/MinterTeam/minter-go-node v1.0.5-0.20191113165918-fa18116d6a26
