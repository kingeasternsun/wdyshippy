// consignment-cli/cli.go
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	pb "../consignment-service/proto/consignment"
	userpb "../user-service/proto/user"
	"golang.org/x/net/context"
	// "google.golang.org/grpc"
	"github.com/micro/cli"
	micro "github.com/micro/go-micro"
	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/metadata"
)

const (
	address         = "localhost:50051"
	defaultFilename = "consignment.json"
)

func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &consignment)
	return consignment, err
}

func main() {
	// Set up a connection to the server.
	cmd.Init()

	// log.Panicln("jwt")
	// Create new greeter client
	client := pb.NewShippingServiceClient("go.micro.srv.consignment", microclient.DefaultClient)

	service := micro.NewService(
		micro.Flags(
			cli.StringFlag{
				Name:  "name",
				Usage: "You full name",
			},
			cli.StringFlag{
				Name:  "email",
				Usage: "Your email",
			},
			cli.StringFlag{
				Name:  "password",
				Usage: "Your password",
			},
			cli.StringFlag{
				Name:  "company",
				Usage: "Your company",
			},
		),
	)

	// Start as service
	service.Init(

		micro.Action(func(c *cli.Context) {

			userClient := userpb.NewUserServiceClient("go.micro.srv.user", microclient.DefaultClient)

			name := "Ewan Valentine"
			email := "ewan.valentine89@gmail.com"
			password := "test123"
			company := "BBC"

			ur, err := userClient.Create(context.TODO(), &userpb.User{
				Name:     name,
				Email:    email,
				Password: password,
				Company:  company,
				Id:       "id",
			})
			if err != nil {
				log.Fatalf("Could not create: %v", err)
			}
			log.Printf("Created: %s", ur.User.Id)

			ugetAll, err := userClient.GetAll(context.Background(), &userpb.Request{})
			if err != nil {
				log.Fatalf("Could not list users: %v", err)
			}
			for _, v := range ugetAll.Users {
				log.Println(v)
			}

			authResponse, err := userClient.Auth(context.TODO(), &userpb.User{
				Email:    email,
				Password: password,
			})

			if err != nil {
				log.Fatalf("Could not authenticate user: %s error: %v\n", email, err)
			}

			log.Printf("Your access token is: %s \n", authResponse.Token)

			// email := c.String("email")
			// password := c.String("password")
			// company := c.String("company")

			// Contact the server and print out its response.
			file := defaultFilename
			if len(os.Args) > 1 {
				file = os.Args[1]
			}

			consignment, err := parseFile(file)

			if err != nil {
				log.Fatalf("Could not parse file: %v", err)
			}

			ctx := metadata.NewContext(context.Background(), map[string]string{
				"token": authResponse.Token,
			})

			r, err := client.CreateConsignment(ctx, consignment)
			if err != nil {
				log.Fatalf("Could not greet: %v", err)
			}
			log.Printf("Created: %t", r.Created)

			getAll, err := client.GetConsignments(ctx, &pb.GetRequest{})
			if err != nil {
				log.Fatalf("Could not list consignments: %v", err)
			}
			for _, v := range getAll.Consignments {
				log.Println(v)
			}

			os.Exit(0)
		}),
	)

	// Run the server
	if err := service.Run(); err != nil {
		log.Println(err)
	}

}
