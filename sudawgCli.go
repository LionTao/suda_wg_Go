package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli"
)

func main() {
	// Create Cli app
	app := cli.NewApp()
	app.Name = "sudawg-agent"
	app.Usage = "Login Program for wg.suda.edu.cn"

	// Global options for login without login command
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "username,u",
			Value: "",
			Usage: "Your username for login",
		},
		cli.StringFlag{
			Name:  "password,p",
			Value: "",
			Usage: "Password of your account",
		},
		cli.StringFlag{
			Name:  "portal,P",
			Value: "wg",
			Usage: "Short name of the portal",
		},
	}

	// Define cli commands with command options
	app.Commands = []cli.Command{
		{
			Name:    "login",
			Aliases: []string{"l"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "username,u",
					Value: "",
					Usage: "Your username for login",
				},
				cli.StringFlag{
					Name:  "password,p",
					Value: "",
					Usage: "Password of your account",
				},
				cli.StringFlag{
					Name:  "portal,P",
					Value: "wg",
					Usage: "Short name of the portal",
				},
			},
			Action: func(ctx *cli.Context) error {
				login(ctx)
				return nil
			},
		},
		{
			Name:    "logout",
			Aliases: []string{"l"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "username,u",
					Value: "",
					Usage: "Your username for login",
				},
				cli.StringFlag{
					Name:  "password,p",
					Value: "",
					Usage: "Password of your account",
				},
				cli.StringFlag{
					Name:  "portal,P",
					Value: "wg",
					Usage: "Short name of the portal",
				},
			},
			Action: func(ctx *cli.Context) error {
				logout(ctx)
				return nil
			},
		},
	}

	//Default action if no commands were given
	app.Action = func(ctx *cli.Context) error {
		login(ctx)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

// Master func of logout
func logout(ctx *cli.Context) {
	wifiLogout()
	time.Sleep(2)
	if NetWorkStatus() {
		username, passsword := getAccount(ctx)
		wgLogout(username, passsword)
	}
	log.Println("Offine now")
}

// Master func of login
func login(ctx *cli.Context) {
	username, password := getAccount(ctx)

	// default portal is wg
	if ctx.String("portal") == "wifi" {
		//change portal to sudawifi if specified
		log.Println("Logging into sudawifi")
		wifiLogin(username, password)
	} else {
		log.Println("Logging into sudawg")
		wgLogin(username, password)
	}

	// Connection test
	if NetWorkStatus() {
		fmt.Println()
		fmt.Println("===============================")
		fmt.Println("Congratulations! Login success")
		fmt.Println("===============================")
	} else {
		fmt.Println("[WARNING] Connection test Failed")
	}
}
