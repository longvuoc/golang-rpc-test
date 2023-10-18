package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"test/calculator/calculatorpb"
	"time"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50069", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("err while dial %v",err)
	}
	defer cc.Close()
	client := calculatorpb.NewCalculatorServiceClient(cc)
	log.Printf("service client %f",client)
	callFindMax(client)
}

func callSum(c calculatorpb.CalculatorServiceClient){
	log.Println("calling sum api")
	resp, err := c.Sum(context.Background(), &calculatorpb.SumRequest{
		Num1:5,
		Num2:6,
	})
	if err != nil {
		log.Fatalf("CALL SUM API ERR %v",err)
	}
	log.Printf("sum api response %v\n",resp.GetResult())
}

func callPND(c calculatorpb.CalculatorServiceClient){
		log.Println("calling sum api")

	stream,err := c.PrimeNumberDecomposition(context.Background(), &calculatorpb.PNDRequest{
		Number: 120,
	})
	if err != nil {  
		log.Fatalf("CallPND %v",err)
	}
	for {
		resp, recErr := stream.Recv()
		if recErr == io.EOF{
		log.Println("SERVER FINISH STREAM")
		return
		}

		log.Printf("prime number %v",resp.GetResult())
	}
}

func callAverage(c calculatorpb.CalculatorServiceClient){
			log.Println("calling sum api")
			stream,err := c.Average(context.Background())
			if err != nil {
				log.Fatalf("call average err %v",err)
			}
			listReq := []calculatorpb.AverageRequest{
				calculatorpb.AverageRequest{
					Num:5,
				},calculatorpb.AverageRequest{
					Num:10,
				},calculatorpb.AverageRequest{
					Num:12,
				},calculatorpb.AverageRequest{
					Num:3,
				},calculatorpb.AverageRequest{
					Num:4.2,
				},
			}
			fmt.Print(listReq)
			for _, req := range listReq{
				err := stream.Send(&req)
				if err != nil {
					log.Fatalf("send average request err %v",err)
				}
			}
			resp, err := stream.CloseAndRecv()
			if err != nil {
				log.Fatalf("receive average request err %v",err)
			}
			log.Printf("average response %+v", resp)
}

func callFindMax(c calculatorpb.CalculatorServiceClient){
	log.Println("calling find max...")
	stream,err := c.FindMax(context.Background())
	if err != nil {
		log.Fatalf("Call average err %v",err)
	}

	listReq := []calculatorpb.FindMaxRequest{
				calculatorpb.FindMaxRequest{
					Num:5,
				},calculatorpb.FindMaxRequest{
					Num:10,
				},calculatorpb.FindMaxRequest{
					Num:12,
				},calculatorpb.FindMaxRequest{
					Num:18,
				},calculatorpb.FindMaxRequest{
					Num:4,
				},
			}

			waitc := make(chan struct{})

			go func (){
				for _, req := range listReq{
					err := stream.Send(&req)
					if err != nil {
						log.Fatalf("Send ")
						break
					}
					time.Sleep(1000 * time.Millisecond)
				}
				stream.CloseSend()
			}()
			go func (){
				for {
					resp, err := stream.Recv()
					if err == io.EOF{
						log.Println("ending find max api ...")
						close(waitc)
					}
					if err != nil {
						log.Fatalf("recv find max %v",err)
						break
					}
					log.Println("max: %v",resp.GetMax())
				}
			}()
			<-waitc
}