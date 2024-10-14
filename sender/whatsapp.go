package main

import (
	"context"
	"fmt"
	"time"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"

	_ "github.com/mattn/go-sqlite3"
)

func getClient() (*whatsmeow.Client, error) {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:mdtest.db?_foreign_keys=on", dbLog)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %v", err)
	}
	client := whatsmeow.NewClient(deviceStore, ClientLog)
	return client, nil
}

func pairPhone(phoneNumber string) {
	client, err := getClient()
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return
	}

	err = client.Connect()
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	defer client.Disconnect()

	code, err := client.PairPhone(phoneNumber, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
	if err != nil {
		fmt.Printf("Failed to pair phone: %v\n", err)
		return
	}
	
	fmt.Println("Pairing phone. Please check your phone for the pairing code.")
	fmt.Printf("Pairing code: %s\n", code)

	// Wait for user to complete pairing on their phone
	fmt.Println("Press Enter after you've completed the pairing on your phone...")
	fmt.Scanln()

	// Check if pairing was successful
	if client.Store.ID == nil {
		fmt.Println("Pairing failed")
		return
	}

	fmt.Println("Successfully paired")
	fmt.Printf("JID: %s\n", client.Store.ID)
}

func send(jid, message string) {
	client, err := getClient()
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return
	}

	err = client.Connect()
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	defer client.Disconnect()
	
	recipient, err := types.ParseJID(jid)
	if err != nil {
		fmt.Printf("Invalid JID %s: %v\n", jid, err)
		return
	}
	
	msg := &waProto.Message{Conversation: proto.String(message)}
	
	fmt.Println("Sending message...")
	resp, err := client.SendMessage(context.Background(), recipient, msg)
	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
		return
	}
	fmt.Printf("Message sent (server timestamp: %s)\n", resp.Timestamp)
	
	time.Sleep(3 * time.Second)
}
