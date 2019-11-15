package service

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/MinterTeam/minter-node-cli/pb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"os"
	"strings"
)

type Api struct {
	client pb.ManagerServiceClient
}

func NewApi(client pb.ManagerServiceClient) *Api {
	return &Api{client: client}
}

func RunCli(socketPath string, agrs []string) {
	cc, err := grpc.Dial("passthrough:///unix:///"+socketPath, grpc.WithInsecure())
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}

	api := NewApi(pb.NewManagerServiceClient(cc))

	app := &cli.App{}
	app.UseShortOptionHandling = true
	app.Commands = []*cli.Command{
		{
			Name:    "dial_peer",
			Aliases: []string{"dp"},
			Usage:   "connect a new peer",
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "persistent", Aliases: []string{"p"}, Required: false},
				&cli.BoolFlag{Name: "json", Aliases: []string{"j"}, Required: false},
			},
			Action: func(c *cli.Context) error {
				_, err := api.client.DealPeer(context.Background(), &pb.DealPeerRequest{
					Address:    c.String("pear"), //todo
					Persistent: c.Bool("persistent"),
				})
				if err != nil {
					return err
				}
				if c.Bool("json") {
					_, _ = fmt.Fprintln(os.Stdout, "OK")
					return nil
				}
				_, _ = fmt.Fprintln(os.Stdout, "OK")
				return nil
			},
		},
		{
			Name:    "prune_blocks",
			Aliases: []string{"pb"},
			Usage:   "delete block information",
			Flags: []cli.Flag{
				&cli.IntFlag{Name: "from", Aliases: []string{"f"}, Required: true},
				&cli.IntFlag{Name: "to", Aliases: []string{"t"}, Required: true},
				&cli.BoolFlag{Name: "json", Aliases: []string{"j"}, Required: false},
			},
			Action: func(c *cli.Context) error {
				_, err := api.client.PruneBlocks(context.Background(), &pb.PruneBlocksRequest{
					FromHeight: c.Int64("from"),
					ToHeight:   c.Int64("to"),
				})
				if err != nil {
					return err
				}
				if c.Bool("json") {
					_, _ = fmt.Fprintln(os.Stdout, "OK")
					return nil
				}
				_, _ = fmt.Fprintln(os.Stdout, "OK")
				return nil
			},
		},
		{
			Name:    "status",
			Aliases: []string{"s"},
			Usage:   "display the current status of the blockchain",
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "json", Aliases: []string{"j"}, Required: false},
			},
			Action: func(c *cli.Context) error {
				response, err := api.client.Status(context.Background(), &empty.Empty{})
				if err != nil {
					return err
				}
				if c.Bool("json") {
					bytes, err := json.Marshal(response)
					if err != nil {
						return err
					}
					_, _ = fmt.Fprintln(os.Stdout, string(bytes))
					return nil
				}
				_, _ = fmt.Fprintln(os.Stdout, proto.MarshalTextString(response))
				return nil
			},
		},
		{
			Name:    "net_info",
			Aliases: []string{"ni"},
			Usage:   "display network data",
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "json", Aliases: []string{"j"}, Required: false},
			},
			Action: func(c *cli.Context) error {
				response, err := api.client.NetInfo(context.Background(), &empty.Empty{})
				if err != nil {
					return err
				}
				if c.Bool("json") {
					bytes, err := json.Marshal(response)
					if err != nil {
						return err
					}
					_, _ = fmt.Fprintln(os.Stdout, string(bytes))
					return nil
				}
				_, _ = fmt.Fprintln(os.Stdout, proto.MarshalTextString(response))
				return nil
			},
		},
		{
			Name:    "exit",
			Aliases: []string{"e"},
			Usage:   "exit",
			Action: func(c *cli.Context) error {
				os.Exit(0)
				return nil
			},
		},
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "test",
			Action: func(c *cli.Context) error {
				_, _ = fmt.Fprintln(os.Stdout, "test ok")
				return nil
			},
		},
	}

	for i := 0; i < len(agrs); i++ {
		if agrs[i] == "exec" {
			if err := app.Run(agrs[i:]); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
			}
			return
		}
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		cmd, err := reader.ReadString('\n')
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			continue
		}
		if err = app.Run(append([]string{""}, strings.Fields(cmd)...)); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
	}
}
