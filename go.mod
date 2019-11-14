module github.com/klim0v/minter-node-cli

go 1.13

require (
	github.com/MinterTeam/minter-go-node v1.0.4
	github.com/golang/protobuf v1.3.2
	github.com/tendermint/tendermint v0.32.6
	github.com/urfave/cli/v2 v2.0.0-alpha.2
	google.golang.org/grpc v1.25.0
)

replace github.com/MinterTeam/minter-go-node v1.0.4 => github.com/MinterTeam/minter-go-node v1.0.5-0.20191108104342-d263a60c747d
