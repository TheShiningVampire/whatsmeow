package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	debugLogs = flag.Bool("debug", false, "Enable debug logs")
)

func main() {
	flag.Parse()
	InitLogger()

	if len(os.Args) < 2 {
		fmt.Println("Usage: ./test <command> [arguments...]")
		fmt.Println("Commands: pair-phone, send")
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "pair-phone":
		if len(args) != 1 {
			fmt.Println("Usage: ./test pair-phone <phone-number>")
			os.Exit(1)
		}
		pairPhone(args[0])
	case "send":
		if len(args) < 2 {
			fmt.Println("Usage: ./test send <jid> <message>")
			os.Exit(1)
		}
		send(args[0], strings.Join(args[1:], " "))
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
