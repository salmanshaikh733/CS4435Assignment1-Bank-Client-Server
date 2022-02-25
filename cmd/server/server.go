package main

import (
	pb "bank/proto"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"strconv"
)

// server is used to implement the server.
type server struct {
	pb.UnimplementedBankServer
}

type BankAccount struct {
	Name      string
	AccountID int64
	Balance   float32
}

var accounts []BankAccount

func (s *server) Deposit(ctx context.Context, in *pb.BankRequest) (*pb.BankResponse, error) {
	for i := range accounts {
		if accounts[i].AccountID == in.GetAccountNum() {
			accounts[i].Balance = accounts[i].Balance + in.GetAmount()
		}
	}

	dat, _ := json.MarshalIndent(accounts, "", "")
	ioutil.WriteFile("./bin/accountsUpdated.json", dat, 0644)

	return &pb.BankResponse{Success: "success"}, nil
}

func (s *server) Withdraw(ctx context.Context, in *pb.BankRequest) (*pb.BankResponse, error) {
	for i := range accounts {
		if accounts[i].AccountID == in.GetAccountNum() {
			accounts[i].Balance = accounts[i].Balance - in.GetAmount()
		}
	}

	dat, _ := json.MarshalIndent(accounts, "", "")
	ioutil.WriteFile("./bin/accountsUpdated.json", dat, 0644)

	return &pb.BankResponse{Success: "success"}, nil
}

func (s *server) Interest(ctx context.Context, in *pb.BankRequest) (*pb.BankResponse, error) {
	for i := range accounts {
		if accounts[i].AccountID == in.GetAccountNum() {
			intRate := in.GetAmount() / 100
			intAmount := intRate * accounts[i].Balance
			accounts[i].Balance = accounts[i].Balance + intAmount
		}
	}

	dat, _ := json.MarshalIndent(accounts, "", "")
	ioutil.WriteFile("./bin/accountsUpdated.json", dat, 0644)

	return &pb.BankResponse{Success: "success"}, nil
}

func main() {

	file, err := ioutil.ReadFile("./bin/accounts.json")

	if err != nil {
		fmt.Println(err.Error())
	}

	err2 := json.Unmarshal(file, &accounts)

	if err2 != nil {
		fmt.Println(err2)
	}

	flag.Parse()

	content, _ := ioutil.ReadFile("port")
	text := string(content)
	log.Println(text)
	portNum, _ := strconv.Atoi(text)
	var (
		port = flag.Int("port", portNum, "The server port")
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBankServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
