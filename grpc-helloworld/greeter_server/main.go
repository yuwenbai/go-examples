///*
// *
// * Copyright 2015, Google Inc.
// * All rights reserved.
// *
// * Redistribution and use in source and binary forms, with or without
// * modification, are permitted provided that the following conditions are
// * met:
// *
// *     * Redistributions of source code must retain the above copyright
// * notice, this list of conditions and the following disclaimer.
// *     * Redistributions in binary form must reproduce the above
// * copyright notice, this list of conditions and the following disclaimer
// * in the documentation and/or other materials provided with the
// * distribution.
// *     * Neither the name of Google Inc. nor the names of its
// * contributors may be used to endorse or promote products derived from
// * this software without specific prior written permission.
// *
// * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
// *
// */

package main

import (
	"encoding/json"
	"log"
	"net"
	"os"

	pb "grpc-helloworld/helloworld"

	"io/ioutil"

	"strconv"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type Config struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

// server is used to implement helloworld.GreeterServer.
type server struct{}

type errorString struct {
	s string
	A string
}

func (e *errorString) Error() string {
	return e.s + e.A
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {

	log.Println("get request")
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func (s *server) SayHello2(in *pb.HelloRequest, gs pb.Greeter_SayHello2Server) error {
	name := in.Name
	for i := 0; i < 100000000; i++ {
		gs.Send(&pb.HelloReply{Message: "Hello " + name + strconv.Itoa(i)})

		log.Println("Hello " + name + strconv.Itoa(i))

		//time.Sleep(time.Second)
	}
	return nil
}

func readConfigFile(configObj interface{}, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	content, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(content, configObj)
	if err != nil {
		log.Println("config_reader readConfigFile json.Unmarshal ", fileName, " error: ", err.Error())
		return err
	}

	return nil
}

func readcfg() *Config {
	var cfg Config
	readConfigFile(&cfg, "gretter_server.json")
	return &cfg
}

func main() {
	cfg := readcfg()
	log.Println("listenning at ", cfg.IP+":"+cfg.Port)
	lis, err := net.Listen("tcp", cfg.IP+":"+cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	//reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
