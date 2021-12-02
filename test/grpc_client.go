package main

import (
	"context"
	"time"

	"google.golang.org/grpc"

	//"google.golang.org/grpc/reflection"
	"log"
	pb "log_transfer/input/input_grpc"

	//"runtime"
	//"sync"
	//"time"
	"fmt"
)

const (
	networkType = "tcp"
	server      = "127.0.0.1"
	port        = "8133"
	parallel    = 50 //连接并行度
	times       = 100000
)

func main() {
	fmt.Println("start")
	exec()
}

func exec() {
	conn, _ := grpc.Dial("127.0.0.1:8133", grpc.WithInsecure())
	defer conn.Close()
	client := pb.NewGreeterClient(conn)
	for {
		ctx := context.Background()
		response, err := client.SendMessage(ctx, &pb.Request{Tag: "aaxx", Msg: "abc"})
		if err != nil {
			log.Printf("err", err)
		}
		log.Printf("response error  %#v", response.Message)
		time.Sleep(1 * time.Second)
	}

}
