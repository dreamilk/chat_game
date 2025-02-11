package main

import (
	"flag"
	"fmt"

	"chat_game/cmd"
)

func Run(cmds map[string]func()) {
	flag.Parse()
	arg := flag.Arg(0)
	if cmd, ok := cmds[arg]; ok {
		cmd()
	} else {
		fmt.Println("command not found, available commands:")
		for k := range cmds {
			fmt.Println(k)
		}
		fmt.Println("usage: ./chat_game server")
	}
}

func main() {
	Run(map[string]func(){
		"server": cmd.Server,
		"worker": cmd.Worker,
	})
}
