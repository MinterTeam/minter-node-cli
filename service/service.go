package service

import (
	"context"
	"encoding/json"
	"github.com/MinterTeam/minter-go-node/config"
	"github.com/MinterTeam/minter-go-node/core/minter"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/klim0v/minter-node-cli/pb"
	rpc "github.com/tendermint/tendermint/rpc/client"
)

type Manager struct {
	blockchain *minter.Blockchain
	tmRPC      *rpc.Local
	cfg        *config.Config
}

func (m *Manager) Status(context.Context, *empty.Empty) (*pb.StatusResponse, error) {
	response := new(pb.StatusResponse)
	resultStatus, err := m.tmRPC.Status()
	if err != nil {
		return response, err
	}

	bytes, err := json.Marshal(resultStatus)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(bytes, response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (m *Manager) NetInfo(context.Context, *empty.Empty) (*pb.NetInfoResponse, error) {
	response := new(pb.NetInfoResponse)
	resultNetInfo, err := m.tmRPC.NetInfo()
	if err != nil {
		return response, err
	}

	bytes, err := json.Marshal(resultNetInfo)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(bytes, response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (m *Manager) PruneBlocks(context.Context, *pb.PruneBlocksRequest) (*empty.Empty, error) {
	//m.blockchain.
	panic("PruneBlocks")
}

func (m *Manager) DealPeer(context.Context, *pb.DealPeerRequest) (*empty.Empty, error) {
	//m.blockchain.
	panic("DealPeer")
}

func NewManager(blockchain *minter.Blockchain, tmRPC *rpc.Local, cfg *config.Config) *Manager {
	return &Manager{blockchain: blockchain, tmRPC: tmRPC, cfg: cfg}
}
