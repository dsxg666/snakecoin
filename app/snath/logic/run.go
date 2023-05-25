package logic

import (
	"github.com/dsxg666/snakecoin/db"
	"github.com/dsxg666/snakecoin/grpc/pb"
	"github.com/dsxg666/snakecoin/grpc/server"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"log"
	"net"
)

func Run(cmd *cobra.Command) {
	rpcPort, _ := cmd.Flags().GetString("rpc.port")
	listen, _ := net.Listen("tcp", ":"+rpcPort)
	grpcServer := grpc.NewServer()
	pb.RegisterRpcServer(grpcServer, &server.Server{
		MptDB:   db.GetDB(db.MPTirePath),
		ChainDB: db.GetDB(db.ChainPath),
		TxDB:    db.GetDB(db.TxPath),
	})
	err := grpcServer.Serve(listen)
	if err != nil {
		log.Panic("Failed to Server: ", err)
		return
	}
}
