package service

import (
	"context"
	"encoding/json"
	"github.com/MinterTeam/minter-go-node/config"
	"github.com/MinterTeam/minter-go-node/core/minter"
	"github.com/MinterTeam/minter-node-cli/pb"
	"github.com/golang/protobuf/ptypes/empty"
	rpc "github.com/tendermint/tendermint/rpc/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type Manager struct {
	blockchain *minter.Blockchain
	tmRPC      *rpc.Local
	cfg        *config.Config
}

func NewManager(blockchain *minter.Blockchain, tmRPC *rpc.Local, cfg *config.Config) *Manager {
	return &Manager{blockchain: blockchain, tmRPC: tmRPC, cfg: cfg}
}

func (m *Manager) Status(context.Context, *empty.Empty) (*pb.StatusResponse, error) {
	resultStatus, err := m.tmRPC.Status()
	if err != nil {
		return new(pb.StatusResponse), status.Error(codes.Internal, err.Error())
	}

	response := &pb.StatusResponse{
		Version:           "",
		LatestBlockHash:   string(resultStatus.SyncInfo.LatestBlockHash),
		LatestAppHash:     string(resultStatus.SyncInfo.LatestBlockHash),
		LatestBlockHeight: resultStatus.SyncInfo.LatestBlockHeight,
		LatestBlockTime:   resultStatus.SyncInfo.LatestBlockTime.Format(time.RFC3339),
		StateHistory:      "",
		TmStatus: &pb.StatusResponse_TmStatus{
			NodeInfo: &pb.NodeInfo{
				ProtocolVersion: nil,
				Id:              string(resultStatus.NodeInfo.ID_),
				ListenAddr:      resultStatus.NodeInfo.ListenAddr,
				Network:         resultStatus.NodeInfo.Network,
				Version:         resultStatus.NodeInfo.Version,
				Channels:        string(resultStatus.NodeInfo.Channels),
				Moniker:         resultStatus.NodeInfo.Moniker,
				Other: &pb.NodeInfo_Other{
					TxIndex:    resultStatus.NodeInfo.Other.TxIndex,
					RpcAddress: resultStatus.NodeInfo.Other.RPCAddress,
				},
			},
			SyncInfo: &pb.StatusResponse_TmStatus_SyncInfo{
				LatestBlockHash:   string(resultStatus.SyncInfo.LatestBlockHash),
				LatestAppHash:     string(resultStatus.SyncInfo.LatestAppHash),
				LatestBlockHeight: resultStatus.SyncInfo.LatestBlockHeight,
				LatestBlockTime:   resultStatus.SyncInfo.LatestBlockTime.Format(time.RFC3339),
				CatchingUp:        resultStatus.SyncInfo.CatchingUp,
			},
			ValidatorInfo: &pb.StatusResponse_TmStatus_ValidatorInfo{
				Address: string(resultStatus.ValidatorInfo.Address),
				PubKey: &pb.StatusResponse_TmStatus_ValidatorInfo_PubKey{
					Type:  "",
					Value: "",
				},
				VotingPower: resultStatus.ValidatorInfo.VotingPower,
			},
		},
	}

	return response, nil
}

func (m *Manager) NetInfo(context.Context, *empty.Empty) (*pb.NetInfoResponse, error) {
	response := new(pb.NetInfoResponse)
	resultNetInfo, err := m.tmRPC.NetInfo()
	if err != nil {
		return response, status.Error(codes.Internal, err.Error())
	}

	bytes, err := json.Marshal(resultNetInfo)
	if err != nil {
		return response, status.Error(codes.Internal, err.Error())
	}

	err = json.Unmarshal(bytes, response)
	if err != nil {
		return response, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}

func (m *Manager) PruneBlocks(ctx context.Context, req *pb.PruneBlocksRequest) (*empty.Empty, error) {
	//m.blockchain.PruneStates(req.FromHeight, req.ToHeight)
	panic("PruneBlocks")
}

func (m *Manager) DealPeer(ctx context.Context, req *pb.DealPeerRequest) (*empty.Empty, error) {
	res := new(empty.Empty)
	_, err := m.tmRPC.DialPeers([]string{req.Address}, req.Persistent)
	if err != nil {
		return res, status.Error(codes.Internal, err.Error())
	}
	return res, nil
}
