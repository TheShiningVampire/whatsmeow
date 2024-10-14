package main

import (
	waLog "go.mau.fi/whatsmeow/util/log"
)

var ClientLog waLog.Logger

func InitLogger() {
	ClientLog = waLog.Stdout("Client", "DEBUG", true)
}
