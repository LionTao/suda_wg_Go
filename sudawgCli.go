package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
)

//type tomlConfig struct {
//	user account
//}

type account struct {
	Username string
	Password string
}

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
					Name:  "file,f",
					Usage: "Path to config file",
				},
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
		user := getAccount(ctx)
		wgLogout(user)
	}
	log.Println("Offine now")
}

// Master func of login
func login(ctx *cli.Context) {
	var user account
	user = getAccount(ctx)

	// default portal is wg
	if ctx.String("portal") == "wifi" {
		//change portal to sudawifi if specified
		log.Println("Logging into sudawifi")
		wifiLogin(user)
	} else {
		log.Println("Logging into sudawg")
		wgLogin(user)
	}

	// Connection test
	if NetWorkStatus() {
		fmt.Println("[SUCCESS] Login success")
	} else {
		fmt.Println("[WARNING] Connection test Failed")
	}
}

// Utility func for acquiring user account
func getAccount(ctx *cli.Context) (user account) {
	var err error
	if ctx.String("file") != "" {
		//var config tomlConfig
		var (
			fp       *os.File
			fcontent []byte
		)
		if fp, err = os.Open("./test.toml"); err != nil {
			fmt.Println("open error ", err)
		}

		if fcontent, err = ioutil.ReadAll(fp); err != nil {
			fmt.Println("ReadAll error ", err)
		}

		//temp := new(account)
		if _, err = toml.Decode(string(fcontent), &user); err != nil {
			fmt.Println("toml.Unmarshal error ", err)
		}
		return user
	}
	if ctx.String("username") != "" {
		user.Username = ctx.String("username")
	} else {
		fmt.Print("Username:")
		_, err = fmt.Scanln(&user.Username)
		for err != nil {
			fmt.Print("Username:")
			_, err = fmt.Scanln(&user.Username)
		}
	}

	if ctx.String("password") != "" {
		user.Password = ctx.String("password")
	} else {
		fmt.Print("Password:")
		_, err = fmt.Scanln(&user.Password)
		for err != nil {
			fmt.Print("Password:")
			_, err = fmt.Scanln(&user.Password)
		}
	}
	return user
}
