/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/wangkebin/grpc-jwt-hello/myjwt"
	pb "github.com/wangkebin/grpc-jwt-hello/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
	mu    sync.Mutex
	names []string
}

func auth(ctx context.Context) error {
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("no metadata found in context")

	}
	tokens := headers["jwt"]
	if len(tokens) < 1 {
		return errors.New("no JWT token found in metadata")
	}
	fmt.Println("received JWT token", tokens[0])

	validator, err := myjwt.NewValidator(os.Args[1])
	if err != nil {
		fmt.Printf("unable to create validator: %v\n", err)
		os.Exit(1)
	}

	if _, err := validator.ValidateTkString(tokens[0]); err != nil {
		fmt.Printf("unable to create validator: %v\n", err)
		os.Exit(1)
	}

	return nil
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	if auth(ctx) != nil {
		fmt.Println("NOT authorized for ", in.GetName())
		return &pb.HelloReply{Message: "Unauthorized " + in.GetName()}, nil
	}
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *server) Chat(stream pb.Greeter_ChatServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		name := in.Name
		s.mu.Lock()
		s.names = append(s.names, name)
		namelist := strings.Join(s.names, ",")
		s.mu.Unlock()

		rep := &pb.HelloReply{
			Message: namelist,
		}
		if err := stream.Send(rep); err != nil {
			return err
		}

	}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
