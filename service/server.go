package service

import (
	"context"
	"github.com/MinterTeam/minter-node-cli/pb"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
	"os"
)

func StartCLIServer(socketPath string, manager *Manager, ctx context.Context) error {
	if err := os.RemoveAll(socketPath); err != nil {
		return err
	}

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

	kill := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			server.GracefulStop()
		case <-kill:
		}
		return
	}()

	if err := group.Wait(); err != nil {
		return err
	}

	close(kill)

	return nil
}
