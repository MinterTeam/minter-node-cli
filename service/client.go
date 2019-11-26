package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MinterTeam/minter-node-cli/pb"
	"github.com/c-bata/go-prompt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"os"
	"strings"
)

type ManagerConsole struct {
	cli *cli.App
}

func NewManagerConsole(cli *cli.App) *ManagerConsole {
	return &ManagerConsole{cli: cli}
}

func (mc *ManagerConsole) Execute(args []string) error {
	return mc.cli.Run(append(make([]string, 1, len(args)+1), args...))
}

func completer(commands cli.Commands) prompt.Completer {
	completions := make([]prompt.Suggest, 0, len(commands))
	for _, command := range commands {
		completions = append(completions, prompt.Suggest{Text: command.Name, Description: command.Description})
	}
	return func(doc prompt.Document) []prompt.Suggest {
		before := doc.TextBeforeCursor()
		wordsBefore := strings.Split(before, " ")
		// the command being entered is the text until the first space
		commandBefore := wordsBefore[0]
		if len(wordsBefore) == 1 {
			return prompt.FilterHasPrefix(completions, commandBefore, true)
		}

		var suggestions []prompt.Suggest
		switch strings.ToLower(commandBefore) {
		case "dial_peer":
			suggestions = append(suggestions, prompt.Suggest{Text: "--address=", Description: "address"})
			suggestions = append(suggestions, prompt.Suggest{Text: "--persistent ", Description: "persistent"})
		case "prune_blocks":
			suggestions = append(suggestions, prompt.Suggest{Text: "--from=", Description: "from"})
			suggestions = append(suggestions, prompt.Suggest{Text: "--to=", Description: "to"})
		default:
			suggestions = append(suggestions, prompt.Suggest{Text: "--json", Description: "echo in json format"})
		}
		return prompt.FilterHasPrefix(suggestions, wordsBefore[len(wordsBefore)-1], true)
	}
}

func (mc *ManagerConsole) Cli() {
	var history []string
	for {
		t := prompt.Input(">>> ", completer(mc.cli.Commands),
			prompt.OptionHistory(history),
			prompt.OptionSelectedSuggestionTextColor(prompt.DarkRed),
		)
		if err := mc.Execute(strings.Fields(t)); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		history = append(history, t)
	}
}

func ConfigureManagerConsole(socketPath string) (*ManagerConsole, error) {
	cc, err := grpc.Dial("passthrough:///unix:///"+socketPath, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := pb.NewManagerServiceClient(cc)

	app := &cli.App{}
	app.CommandNotFound = func(ctx *cli.Context, cmd string) {
		fmt.Println(fmt.Sprintf("No help topic for '%v'", cmd))
	}
	app.UseShortOptionHandling = true
	app.Commands = []*cli.Command{
		{
			Name:    "dial_peer",
			Aliases: []string{"dp"},
			Usage:   "connect a new peer",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "address", Aliases: []string{"a"}, Required: true},
				&cli.BoolFlag{Name: "persistent", Aliases: []string{"p"}, Required: false},
				&cli.BoolFlag{Name: "json", Aliases: []string{"j"}, Required: false},
			},
			Action: func(c *cli.Context) error {
				_, err := client.DealPeer(context.Background(), &pb.DealPeerRequest{
					Address:    c.String("address"),
					Persistent: c.Bool("persistent"),
				})
				if err != nil {
					return err
				}
				if c.Bool("json") {
					fmt.Println("OK")
					return nil
				}
				fmt.Println("OK")
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
				_, err := client.PruneBlocks(context.Background(), &pb.PruneBlocksRequest{
					FromHeight: c.Int64("from"),
					ToHeight:   c.Int64("to"),
				})
				if err != nil {
					return err
				}
				if c.Bool("json") {
					fmt.Println("OK")
					return nil
				}
				fmt.Println("OK")
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
				response, err := client.Status(context.Background(), &empty.Empty{})
				if err != nil {
					return err
				}
				if c.Bool("json") {
					bytes, err := json.Marshal(response)
					if err != nil {
						return err
					}
					fmt.Println(string(bytes))
					return nil
				}
				fmt.Println(proto.MarshalTextString(response))
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
				response, err := client.NetInfo(context.Background(), &empty.Empty{})
				if err != nil {
					return err
				}
				if c.Bool("json") {
					bytes, err := json.Marshal(response)
					if err != nil {
						return err
					}
					fmt.Println(string(bytes))
					return nil
				}
				fmt.Println(proto.MarshalTextString(response))
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
				fmt.Println("test ok")
				return nil
			},
		},
	}

	return NewManagerConsole(app), nil
}
