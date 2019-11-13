package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/klim0v/minter-node-cli/pb"
	"google.golang.org/grpc"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	s, err := filepath.Abs(filepath.Join(".", "config", "file.sock"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(0)
	}

	var socketPath = flag.String("config", s, "path to dir with config socketPath")
	flag.Args()

	cc, err := grpc.Dial("passthrough:///unix:///"+*socketPath, grpc.WithInsecure())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(0)
	}

	api := NewApi(pb.NewManagerServiceClient(cc))

	stop := make(chan struct{})

	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "exec" {
			if i+1 == len(os.Args) {
				_, err = fmt.Fprintln(os.Stderr, "use 'exec [command]'")
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				return
			}
			response, err := api.runCommand(strings.Join(os.Args[i+1:], " "), stop)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			_, err = fmt.Fprintln(os.Stdout, response)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			return
		}
	}

	//if _, err := fmt.Fprintln(os.Stdout, helpText); err != nil {
	//	fmt.Fprintln(os.Stderr, err)
	//	os.Exit(0)
	//}

	reader := bufio.NewReader(os.Stdin)
	for {
		select {
		case <-stop:
			_, err = fmt.Fprintln(os.Stdout, "Buy!")
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			return
		default:
			fmt.Print("$ ")
			cmdString, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			response, err := api.runCommand(cmdString, stop)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			_, err = fmt.Fprintln(os.Stdout, response)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
		}
	}
}

const helpText = `exit -- 
status -- 
help -- show this help message
Example: status --json`

type Api struct {
	client pb.ManagerServiceClient
}

func NewApi(client pb.ManagerServiceClient) *Api {
	return &Api{client: client}
}

type Reply interface {
	JSON() string
	View() string
	setRequest(interface{})
	load() error
}

type StatusReply struct {
	request  *empty.Empty
	function func(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*pb.StatusResponse, error)
	response *pb.StatusResponse
}

func (s *StatusReply) JSON() string {
	bytes, _ := json.Marshal(s.response)
	return string(bytes)
}

func (s *StatusReply) View() string {
	return proto.MarshalTextString(s.response)
}

func (s *StatusReply) setRequest(request interface{}) {
	s.request = request.(*empty.Empty)
}

func (s *StatusReply) load() (err error) {
	if s.request == nil {
		return errors.New("request is nil")
	}
	s.response, err = s.function(context.Background(), &empty.Empty{})
	return err
}

type NetInfoReply struct {
	request  *empty.Empty
	function func(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*pb.NetInfoResponse, error)
	response *pb.NetInfoResponse
}

func (s *NetInfoReply) JSON() string {
	bytes, _ := json.Marshal(s.response)
	return string(bytes)
}

func (s *NetInfoReply) View() string {
	return proto.MarshalTextString(s.response)
}

func (s *NetInfoReply) setRequest(request interface{}) {
	s.request = request.(*empty.Empty)
}

func (s *NetInfoReply) load() (err error) {
	if s.request == nil {
		return errors.New("request is nil")
	}
	s.response, err = s.function(context.Background(), &empty.Empty{})
	return err
}

type DialPeerReply struct {
	request  *pb.DealPeerRequest
	function func(ctx context.Context, in *pb.DealPeerRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	response *empty.Empty
}

func (s *DialPeerReply) JSON() string {
	bytes, _ := json.Marshal(s.response)
	return string(bytes)
}

func (s *DialPeerReply) View() string {
	return proto.MarshalTextString(s.response)
}

func (s *DialPeerReply) setRequest(request interface{}) {
	s.request = request.(*pb.DealPeerRequest)
}

func (s *DialPeerReply) load() (err error) {
	if s.request == nil {
		return errors.New("request is nil")
	}
	s.response, err = s.function(context.Background(), s.request)
	return err
}

type PruneBlocksReply struct {
	request  *pb.PruneBlocksRequest
	function func(ctx context.Context, in *pb.PruneBlocksRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	response *empty.Empty
}

func (s *PruneBlocksReply) JSON() string {
	bytes, _ := json.Marshal(s.response)
	return string(bytes)
}

func (s *PruneBlocksReply) View() string {
	return proto.MarshalTextString(s.response)
}

func (s *PruneBlocksReply) setRequest(request interface{}) {
	s.request = request.(*pb.PruneBlocksRequest)
}

func (s *PruneBlocksReply) load() (err error) {
	if s.request == nil {
		return errors.New("request is nil")
	}
	s.response, err = s.function(context.Background(), s.request)
	return err
}

type SimpleTextReply struct {
	text string
}

func (s *SimpleTextReply) JSON() string {
	return s.text
}

func (s *SimpleTextReply) View() string {
	return s.text
}

func (s *SimpleTextReply) setRequest(request interface{}) {
	s.text = request.(string)
}

func (s *SimpleTextReply) load() (err error) {
	return nil
}

func (api *Api) runCommand(commandStr string, finish chan<- struct{}) (string, error) {
	commandStr = strings.TrimSuffix(commandStr, "\n")
	command := strings.Fields(commandStr)
	var cmd Reply
	var args interface{}
	switch command[0] {
	case "exit":
		close(finish)
		return "", nil
	case "dial_peer":
		cmd = &DialPeerReply{function: api.client.DealPeer}
		args = &pb.DealPeerRequest{}
	case "prune_blocks":
		cmd = &PruneBlocksReply{function: api.client.PruneBlocks}
		//--from=1 --to=2
		args = &pb.PruneBlocksRequest{
			FromHeight: 0,
			ToHeight:   0,
		}
	case "status":
		cmd = &StatusReply{function: api.client.Status}
		args = new(empty.Empty)
	case "net_info":
		cmd = &NetInfoReply{function: api.client.NetInfo}
		args = new(empty.Empty)
	case "help":
		cmd = &SimpleTextReply{}
		args = helpText
	default:
		cmd = &SimpleTextReply{}
		args = "not found command. Use \"help\" to show all commands"
	}
	cmd.setRequest(args)

	err := cmd.load()
	if err != nil {
		return "", err
	}

	if strings.Contains(commandStr, "--json") {
		return cmd.JSON(), nil
	}

	return cmd.View(), nil
}
