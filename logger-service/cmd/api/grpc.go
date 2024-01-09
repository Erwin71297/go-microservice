package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"log-service/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(context context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	log.Println("enter write log")
	input := req.GetLogEntry()

	//Write the log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	log.Println("log Entry: ", logEntry)

	err := l.Models.LogEntry.InsertGRPC(logEntry)
	if err != nil {
		log.Println("error inserting data")
		//res := &logs.LogResponse{Result: "Failed"}
		return nil, err
	}

	// Return Response
	res := &logs.LogResponse{Result: "Logged!"}
	return res, nil
}

func gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen to gRPC: %v", err)
	}

	s := grpc.NewServer()

	logs.RegisterLogServiceServer(s, &LogServer{Models: data.Models{}})

	log.Printf("gRPC Server started on port %s", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
