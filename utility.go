package main

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/urfave/cli"
)

// Utility func for acquiring user account
func getAccount(ctx *cli.Context) (uname, pwd string) {
	var username, password string
	if ctx.String("username") != "" {
		username = ctx.String("username")
	} else {
		fmt.Print("Username:")
		_, err := fmt.Scanln(&username)
		for err != nil {
			fmt.Print("Username:")
			_, err = fmt.Scanln(&username)
		}
	}

	if ctx.String("password") != "" {
		username = ctx.String("password")
	} else {
		fmt.Print("Password:")
		_, err := fmt.Scanln(&password)
		for err != nil {
			fmt.Print("Password:")
			_, err = fmt.Scanln(&password)
		}
	}
	return username, password
}

// NetWorkStatus : It checks the Connection to baidu.com
func NetWorkStatus() bool {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "baidu.com", "-n", "1", "-w", "5") // #nosec
	} else {
		cmd = exec.Command("ping", "baidu.com", "-c", "1", "-W", "5") // #nosec
	}
	err := cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	fmt.Println("Net Status: OK")
	return true
}
