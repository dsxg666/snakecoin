package logic

import (
	"fmt"
	"github.com/dsxg666/snakecoin/db"
	"net"

	"github.com/dsxg666/snakecoin/grpc/pb"
	"github.com/dsxg666/snakecoin/grpc/server"
	"google.golang.org/grpc"
)

func Run() {
	listen, _ := net.Listen("tcp", ":8545")
	grpcServer := grpc.NewServer()
	pb.RegisterRpcServer(grpcServer, &server.Server{
		MptDB:   db.GetDB(db.MPTirePath),
		ChainDB: db.GetDB(db.ChainPath),
		TxDB:    db.GetDB(db.TxPath),
	})
	err := grpcServer.Serve(listen)
	if err != nil {
		fmt.Println("Failed to Server: ", err)
		return
	}
}
