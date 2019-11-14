package service

import (
	"context"
	"encoding/json"
	"github.com/MinterTeam/minter-go-node/config"
	"github.com/MinterTeam/minter-go-node/core/minter"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/klim0v/minter-node-cli/pb"
	rpc "github.com/tendermint/tendermint/rpc/client"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
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
		return response, status.Error(codes.FailedPrecondition, err.Error())
	}

	bytes, err := json.Marshal(resultNetInfo)
	if err != nil {
		return response, status.Error(codes.FailedPrecondition, err.Error())
	}

	err = json.Unmarshal(bytes, response)
	if err != nil {
		return response, status.Error(codes.FailedPrecondition, err.Error())
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
		return res, status.Error(codes.FailedPrecondition, err.Error())
	}
	return res, nil
}

func NewManager(blockchain *minter.Blockchain, tmRPC *rpc.Local, cfg *config.Config) *Manager {
	return &Manager{blockchain: blockchain, tmRPC: tmRPC, cfg: cfg}
}

func StartCLIServer(socketPath string, manager *Manager, ctx context.Context) error {
	lis, err := net.ListenUnix("unix", &net.UnixAddr{Name: socketPath, Net: "unix"})
	if err != nil {
		return err
	}

	server := grpc.NewServer()

	pb.RegisterManagerServiceServer(server, manager)

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		err := server.Serve(lis)
		if err != nil {
			return err
		}
		return nil
	})

	func() {
		<-ctx.Done()
		server.GracefulStop()
	}()

	err = group.Wait()
	if err != nil {
		return err
	}
}
