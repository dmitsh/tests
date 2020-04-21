package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	connect "github.com/dmitsh/grpctest/pkg/proto"
	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Client struct {
	c connect.DataServiceClient
}

func NewClient(addr string) (*Client, error) {
	maxMsgSize := 1024 * 1024 * 20
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)))
	if err != nil {
		return nil, err
	}
	return &Client{
		c: connect.NewDataServiceClient(conn),
	}, nil
}

func (cln *Client) SendDataRequest(messages, units, bufferSize int32) error {

	dataRequest := &connect.DataRequest{
		Messages:   messages,
		Units:      units,
		BufferSize: bufferSize,
	}
	start := time.Now()
	stream, err := cln.c.Data(context.Background(), dataRequest)
	if err != nil {
		return err
	}

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("Client: EOF")
			break
		}
		if err != nil {
			fmt.Println("Client error: ", err.Error())
			return err
		}
	}
	fmt.Println("Client: done in", time.Now().Sub(start).String())
	return nil
}

func main() {
	var address string
	var messages, units, bufferSize int32
	a := kingpin.New(filepath.Base(os.Args[0]), "gRPC test client")
	a.HelpFlag.Short('h')
	a.Flag("address", "gRPC server address.").Short('a').Default("localhost:5432").StringVar(&address)
	a.Flag("messages", "Number of messages.").Short('m').Default("1000").Int32Var(&messages)
	a.Flag("units", "Number of units per message.").Short('u').Default("1000").Int32Var(&units)
	a.Flag("bufferSize", "Unit buffer size.").Short('b').Default("1000").Int32Var(&bufferSize)

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		fmt.Println("Error parsing commandline arguments")
		a.Usage(os.Args[1:])
		os.Exit(2)
	}

	cln, err := NewClient(address)
	if err != nil {
		fmt.Println("Client error:", err.Error())
	}
	cln.SendDataRequest(messages, units, bufferSize)
}
