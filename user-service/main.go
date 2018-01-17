package main

import (
	"fmt"

	pb "../user-service/proto/user"
	"github.com/micro/go-micro"
	// _ "github.com/micro/go-plugins/broker/nats"
)

func main() {

	repo := &UserRepository{}
	repo.New()

	tokenService := &TokenService{repo}

	// Create a new service. Optionally include some options here.
	srv := micro.NewService(

		// This name must match the package name given in your protobuf definition
		micro.Name("go.micro.srv.user"),
		micro.Version("latest"),
	)

	// Init will parse the command line flags.
	srv.Init()

	// Get instance of the broker using our defaults
	pubsub := srv.Server().Options().Broker
	fmt.Println("pubsu", pubsub)

	// Register handler
	pb.RegisterUserServiceHandler(srv.Server(), &service{repo, tokenService, pubsub})

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
