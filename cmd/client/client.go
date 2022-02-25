package main

import (
	"bufio"
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	pb "bank/proto"
)

func main() {
	flag.Parse()

	content, _ := ioutil.ReadFile("port")
	text := string(content)
	address := "localhost:" + text

	var (
		addr = flag.String("addr", address, "the address to connect to")
	)

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBankClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	file, err := os.Open("./bin/input")
	if err != nil {
		log.Fatalf("could not open file")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		input := strings.Fields(scanner.Text())
		accountNum, _ := strconv.ParseInt(input[1], 10, 64)
		amount, _ := strconv.ParseInt(input[2], 10, 64)

		if input[0] == "deposit" {
			_, err := c.Deposit(ctx, &pb.BankRequest{AccountNum: accountNum, Amount: float32(amount)})
			if err != nil {
				log.Fatalf("could not do deposit: %v", err)
			}
		} else if input[0] == "withdraw" {
			_, err := c.Withdraw(ctx, &pb.BankRequest{AccountNum: accountNum, Amount: float32(amount)})
			if err != nil {
				log.Fatalf("could not do withdrawl: %v", err)
			}
		} else if input[0] == "interest" {
			if amount > 10 {
				log.Fatalf("Interest rate above 10, ERROR")
			} else {
				_, err := c.Interest(ctx, &pb.BankRequest{AccountNum: accountNum, Amount: float32(amount)})
				if err != nil {
					log.Fatalf("could not do interest: %v", err)
				}
			}
		}
	}
}
