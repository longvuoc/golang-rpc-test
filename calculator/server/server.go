package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"test/calculator/calculatorpb"

	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context,req *calculatorpb.SumRequest) (*calculatorpb.SumResponse,error){
	log.Println("sum called ...")
	resp := &calculatorpb.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}
	return resp,nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50069")
	if err != nil {
		log.Fatalf("err while create listen %v", err)
	}
	s := grpc.NewServer()
	fmt.Println("Calculator is running")
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("err while serve %v", err)
	}
}
