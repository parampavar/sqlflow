// Copyright 2020 The SQLFlow Authors. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package modelzooserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	pb "sqlflow.org/sqlflow/pkg/proto"
)

func startServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50055))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterModelZooServerServer(grpcServer, &modelZooServer{})

	grpcServer.Serve(lis)
}

func serverIsReady(addr string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return false
	}
	err = conn.Close()
	return err == nil
}

func waitPortReady(addr string, timeout time.Duration) {
	// Set default timeout to
	if timeout == 0 {
		timeout = time.Duration(1) * time.Second
	}
	for !serverIsReady(addr, timeout) {
		time.Sleep(1 * time.Second)
	}
}

func TestModelZooServer(t *testing.T) {
	go startServer()
	waitPortReady("localhost:50055", 0)

	conn, err := grpc.Dial(":50055", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("create client error: %v", err)
	}
	defer conn.Close()

	client := pb.NewModelZooServerClient(conn)
	res, err := client.ListModelDefs(context.Background(), &pb.ListModelRequest{Start: 0, Size: -1})
	if err != nil {
		t.Fatalf("call ListModelDefs error: %v", err)
	}
	if res == nil {
		t.Errorf("res should not be nil")
	}
}
