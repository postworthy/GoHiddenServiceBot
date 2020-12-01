package main

import (
	"bufio"
	"os"
)

func main() {
	bot := &SshServiceManagerTelegramBot{}
	bot.Init()
	bot.Run()
	for {
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}
