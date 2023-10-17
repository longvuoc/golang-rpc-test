package main

import (
	"context"
	"fmt"
	"io"
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
func (*server) PrimeNumberDecomposition(req *calculatorpb.PNDRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	log.Printf("PrimeNumberDecomposition is calling...")
	k :=int32(2)
	N := req.GetNumber()
	log.Printf(" %v",k)
	log.Printf("%d\n",N)
	for N>1{
		if N%k ==0{
			N=N/k
			stream.Send(&calculatorpb.PNDResponse{
				Result: k,

			})
		}else {
			fmt.Printf("%d\n", k) 
			k ++
			log.Printf("k increase to %v",k)
		}
	}
	return nil
}

func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	var total float32
	var count int
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			resp := &calculatorpb.AverageResponse{
				Result: total / float32(count),
			}
			return stream.SendAndClose(resp)
		}
		if err != nil {
			log.Fatalf("err while Recv Average %v",err)
		}
		log.Println("receive num %v",req)
		total += req.GetNum()
		count++
	}
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

