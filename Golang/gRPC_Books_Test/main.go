package main

import (
	gRPC_Books_Test "gRPC_Books_Test/api/proto"
	"gRPC_Books_Test/pkg/finder"

	//"gRPC_Books_Test/pkg/finder"
	_ "github.com/go-sql-driver/mySQL"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	s := grpc.NewServer()
	srv := &finder.GRPCServer{}
	gRPC_Books_Test.RegisterBooksAndAuthorsServer(s, srv)

	//Создаём слушатель
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
