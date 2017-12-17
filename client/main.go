package main

import (
	"context"
	"io"
	"log"

	"google.golang.org/grpc"

	pb "github.com/olefile/grpc_sample/customer"
)

const (
	address = "127.0.0.1:50051"
)

func createCustomer(client pb.CustomerClient, customer *pb.CustomerRequest) {
	resp, err := client.CreateCustomer(context.Background(), customer)

	if err != nil {
		log.Fatalf("Could not create Customer %v", err)
	}

	if resp.Success {
		log.Printf("A new Customer has been added with id %d", resp.Id)
	}
}

func getCustomers(client pb.CustomerClient, filter *pb.CustomerFilter) {
	stream, err := client.GetCustomers(context.Background(), filter)

	if err != nil {
		log.Fatalf("Error on getting customers %v", err)
	}

	for {
		customer, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("%v.GetCustomers(_) = _, %v", client, err)
		}

		log.Printf("Customer: %v", customer)
	}
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannot connect to grpc customers %v", err)
	}

	defer conn.Close()

	client := pb.NewCustomerClient(conn)

	customer := &pb.CustomerRequest{
		Id:    101,
		Name:  "Egor Gordoov",
		Email: "egor.gorodov@gmail.com",
		Phone: "+79375837362",
	}

	createCustomer(client, customer)

	filter := &pb.CustomerFilter{Keyword: "Egor"}
	getCustomers(client, filter)
}
