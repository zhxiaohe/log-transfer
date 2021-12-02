package input

import (
	"context"
	"log"
	pb "log_transfer/input/input_grpc"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcComponent struct {
	lis   net.Listener
	gs    *grpc.Server
	event chan map[string]interface{}
}

func grpcInit() *grpcComponent {
	lis, err := net.Listen("tcp", "0.0.0.0:8133")
	if err != nil {
		log.Fatal("faild", err)
	}
	gs := grpc.NewServer()
	us := &grpcComponent{lis: lis, gs: gs}
	pb.RegisterGreeterServer(gs, us)
	reflection.Register(gs)
	return us
}

func (g *grpcComponent) SendMessage(ctx context.Context, req *pb.Request) (*pb.Reply, error) {
	// log.Printf("recve: %s %s", req.Msg, req.Tag)
	g.event <- map[string]interface{}{"tag": req.Tag, "value": req.Msg}
	return &pb.Reply{
		Message: "sucess",
	}, nil
}

func (g *grpcComponent) start(event chan map[string]interface{}) error {
	g.event = event
	log.Println("grpc server start ")
	if err := g.gs.Serve(g.lis); err != nil {
		log.Fatal("faild to server", err)
	}
	return nil
}
