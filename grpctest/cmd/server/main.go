package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	connect "github.com/dmitsh/grpctest/pkg/proto"
	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Server struct {
	s    *grpc.Server
	addr string
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (srv Server) StartServer() error {
	fmt.Println("starting server")
	if srv.s != nil {
		fmt.Println("server already running")
		return nil
	}
	lis, err := net.Listen("tcp", srv.addr)
	if err != nil {
		return err
	}

	srv.s = grpc.NewServer()

	connect.RegisterDataServiceServer(srv.s, srv)

	//go func() {
	if err := srv.s.Serve(lis); err != nil {
		fmt.Println("failed to serve. err", err)
		srv.s = nil
	}
	//}()
	return nil
}

func (srv Server) Data(req *connect.DataRequest, s connect.DataService_DataServer) error {
	fmt.Println("Server: sending", req.Messages, "messages with", req.Units, "units,", req.BufferSize, "bytes each.")
	var i int32
	ws := &connect.DataMessage{
		Units: make([]*connect.DataUnit, req.Units),
	}
	for i = 0; i < req.Units; i++ {
		ws.Units[i] = &connect.DataUnit{Buffer: make([]byte, req.BufferSize)}
	}
	start := time.Now()
	for i = 0; i < req.Messages; i++ {
		if err := s.Send(ws); err != nil {
			fmt.Println("Server: Error", err.Error())
			return err
		}
	}
	fmt.Println("Server: done in", time.Now().Sub(start).String())
	return nil
}

func main() {
	var address string
	a := kingpin.New(filepath.Base(os.Args[0]), "gRPC test server")
	a.HelpFlag.Short('h')
	a.Flag("address", "gRPC server address.").Short('a').Default(":5432").StringVar(&address)

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		fmt.Println("Error parsing commandline arguments")
		a.Usage(os.Args[1:])
		os.Exit(2)
	}

	srv := NewServer(address)
	if err := srv.StartServer(); err != nil {
		fmt.Println("Server: error", err.Error())
	}
}
