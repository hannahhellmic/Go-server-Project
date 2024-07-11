package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	pb "awesomeProject/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:5678", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("disconnected: %v", err)
	}
	defer conn.Close()
	c := pb.NewBankAccountServiceClient(conn)

	reader := bufio.NewReader(os.Stdin)

	var cnt int = 0

	for {
		if cnt == 0 {
			fmt.Println("possible operations: create, get, update_balance, change_name, delete")
			cnt++
		} else {
			fmt.Println("something else? yes/no")
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer == "yes" {
				fmt.Print("enter operation: ")
			} else if answer == "no" {
				return
			} else if answer != "no" && answer != "yes" {
				fmt.Println("ROBOT DETECTED! WARNING!")
				return
			}
		}

		operation, _ := reader.ReadString('\n')
		operation = strings.TrimSpace(operation)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		switch operation {
		case "create":
			fmt.Print("name: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)
			fmt.Print("start balance: ")
			balance, _ := reader.ReadString('\n')
			balance = strings.TrimSpace(balance)
			bal, err := strconv.ParseInt(balance, 10, 64)
			if err != nil {
				log.Fatalf("failed to parse balance: %v", err)
				continue
			}
			_, err = c.Create(ctx, &pb.CreateRequest{Name: name, Balance: bal})
			if err != nil {
				log.Printf("could not create account: %v", err)
				continue
			}
			log.Println("account created successfully")

		case "get":
			fmt.Print("name: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)
			getResp, err := c.Get(ctx, &pb.GetRequest{Name: name})
			if err != nil {
				log.Printf("could not get account: %v", err)
				continue
			}
			log.Printf("your account: %v", getResp)

		case "update_balance":
			fmt.Print("name: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)
			fmt.Print("sum change: ")
			balanceStr, _ := reader.ReadString('\n')
			balanceStr = strings.TrimSpace(balanceStr)
			balance, err := strconv.ParseInt(balanceStr, 10, 64)
			if err != nil {
				log.Printf("invalid balance: %v", err)
				continue
			}
			_, err = c.UpdateBalance(ctx, &pb.UpdateBalanceRequest{Name: name, Balance: balance})
			if err != nil {
				log.Printf("could not update balance: %v", err)
				continue
			}
			log.Println("balance has been updated successfully")

		case "change_name":
			fmt.Print("current name: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)
			fmt.Print("new name: ")
			newName, _ := reader.ReadString('\n')
			newName = strings.TrimSpace(newName)
			_, err = c.UpdateName(ctx, &pb.UpdateNameRequest{Name: name, NewName: newName})
			if err != nil {
				log.Printf("could not update name: %v", err)
				continue
			}
			log.Println("account name has been changed successfully")

		case "delete":
			fmt.Print("name: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)
			fmt.Println("are you sure? yes/no")
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer == "yes" {
				_, err = c.Delete(ctx, &pb.DeleteRequest{Name: name})
				if err != nil {
					log.Printf("could not delete account: %v", err)
					continue
				}
				log.Println("account has been deleted successfully")
			} else if answer == "no" {
				continue
			} else {
				log.Printf("incorrect answer")
			}
		default:
			fmt.Println("unknown operation. please try again.")
		}
	}
}
