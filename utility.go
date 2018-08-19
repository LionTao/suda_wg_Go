package main

import (
	"fmt"
	"os/exec"
	"runtime"
)

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
		//log.Println(err.Error())
		return false
	}
	fmt.Println("Net Status: OK")
	return true
}
