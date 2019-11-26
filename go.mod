module github.com/MinterTeam/minter-node-cli

go 1.13

require (
	github.com/MinterTeam/minter-go-node v1.0.5-0.20191113110340-a46b8ef88084
	github.com/c-bata/go-prompt v0.2.3
	github.com/golang/protobuf v1.3.2
	github.com/mattn/go-tty v0.0.3 // indirect
	github.com/pkg/term v0.0.0-20190109203006-aa71e9d9e942 // indirect
	github.com/tendermint/tendermint v0.32.6
	github.com/urfave/cli/v2 v2.0.0-alpha.2
	google.golang.org/grpc v1.25.0
)

replace github.com/MinterTeam/minter-go-node v1.0.4 => github.com/MinterTeam/minter-go-node v1.0.5-0.20191113165918-fa18116d6a26
