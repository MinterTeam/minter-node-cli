package service

import (
	"context"
	"github.com/MinterTeam/minter-go-node/config"
	"github.com/MinterTeam/minter-go-node/core/minter"
	rpc "github.com/tendermint/tendermint/rpc/client"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"
)

func TestStartCLIServer(t *testing.T) {
	var (
		blockchain *minter.Blockchain
		tmRPC      *rpc.Local
		cfg        *config.Config
	)
	ctx, cancel := context.WithCancel(context.Background())
	socketPath, _ := filepath.Abs(filepath.Join(".", "file.sock"))
	err := ioutil.WriteFile(socketPath, []byte("address already in use"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		err := StartCLIServer(socketPath, NewManager(blockchain, tmRPC, cfg), ctx)
		if err != nil {
			t.Log(err)
		}
	}()
	time.Sleep(time.Millisecond)
	RunCli(socketPath, []string{"exec", "help"})
	cancel()
}
