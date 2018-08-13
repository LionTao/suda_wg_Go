package main

import (
	"encoding/base64"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "sudawg-agent"
	app.Usage = "Login Program for wg.suda.edu.cn"
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

	//app.Action = func(ctx *cli.Context) error{
	//	login(ctx)
	//	return nil
	//}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func logout(ctx *cli.Context) {
	wifiLogout()
	time.Sleep(2)
	if NetWorkStatus() {
		username, passsword := getAccount(ctx)
		wgLogout(username, passsword)
	}
	log.Println("Offine now")
}

func getAccount(ctx *cli.Context) (uname, pwd string) {
	var username, password string
	if ctx.String("username") != "" {
		//fmt.Println("here")
		username = ctx.String("username")
	} else {
		fmt.Print("Username:")
		fmt.Scanln(&username)
	}

	if ctx.String("password") != "" {
		//fmt.Println("here")
		username = ctx.String("password")
	} else {
		fmt.Print("Password:")
		fmt.Scanln(&password)
	}
	return username, password
}

func login(ctx *cli.Context) {
	username, password := getAccount(ctx)

	//default is wg
	if ctx.String("portal") == "wifi" {
		log.Println("Logging into sudawifi")
		wifiLogin(username, password)
	} else {
		log.Println(ctx.String("portal"))
		log.Println("Logging into sudawg")
		wgLogin(username, password)
	}

	if NetWorkStatus() {
		fmt.Println("Login success")
	} else {
		fmt.Println("[WARNING] Connection test Failed")
	}
}

func wgLogout(username string, password string) {
	ev, vs := get_wg_param()
	wg_POST(ev, vs, username, password, "退出网关")
}
func wgLogin(username string, password string) {

	ev, vs := get_wg_param()
	wg_POST(ev, vs, username, password, "登陆网关")
}

func get_wg_param() (eventvalidation, viewstate string) {
	url := "http://wg.suda.edu.cn/indexn.aspx"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	//body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		// handle error
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(body))
	//fmt.Println()
	var VIEWSTATE, EVENTVALIDATION string
	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		//fmt.Println(s.Attr("name"))
		name, _ := s.Attr("name")
		if name != "" {
			if name == "__VIEWSTATE" {
				//fmt.Println("VIEWSTATE")
				//fmt.Println(s.Attr("value"))
				VIEWSTATE, _ = s.Attr("value")
			} else if name == "__EVENTVALIDATION" {
				//fmt.Println("EVENTVALIDATION")
				//fmt.Println(s.Attr("value"))
				EVENTVALIDATION, _ = s.Attr("value")
			}
		}
	})

	return EVENTVALIDATION, VIEWSTATE
}

func wg_POST(EVENTVALIDATION string, VIEWSTATE string, username string, password string, action string) {
	if EVENTVALIDATION == "" || VIEWSTATE == "" {
		log.Fatal("Oh no! no param was given")
	}

	url := "http://wg.suda.edu.cn/indexn.aspx"

	var r http.Request
	r.ParseForm()
	r.Form.Add("__EVENTTARGET", "")
	r.Form.Add("__EVENTARGUMENT=", "")
	r.Form.Add("__VIEWSTATE", VIEWSTATE)
	r.Form.Add("__EVENTVALIDATION", EVENTVALIDATION)
	r.Form.Add("TextBox1", username)
	r.Form.Add("TextBox2", password)
	r.Form.Add("nw", "RadioButton2")
	r.Form.Add("tm", "RadioButton8")
	if action == "登陆网关" {
		r.Form.Add("Button1", action)
	} else if action == "退出网关" {
		r.Form.Add("Button4", action)
	}

	bodystr := strings.TrimSpace(r.Form.Encode())
	request, err := http.NewRequest("POST", url, strings.NewReader(bodystr))
	if err != nil {
		log.Fatal("POST Failed")
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Connection", "Keep-Alive")

	resp, err := http.DefaultClient.Do(request)
	if err != nil || resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		fmt.Println(err)
		log.Fatal("[ERROR] An error occurred in sending request")

	}

}

func NetWorkStatus() bool {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "baidu.com", "-n", "1", "-w", "5")
	} else {
		cmd = exec.Command("ping", "baidu.com", "-c", "1", "-W", "5")
	}
	//fmt.Println("NetWorkStatus Start:", time.Now().Unix())
	err := cmd.Run()
	//fmt.Println("NetWorkStatus End  :", time.Now().Unix())
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else {
		fmt.Println("Net Status: OK")
	}
	return true
}

func wifiLogin(username string, password string) {
	url := "http://a.suda.edu.cn/index.php/index/login"

	var r http.Request
	r.ParseForm()
	r.Form.Add("username", username)
	r.Form.Add("password", base64.StdEncoding.EncodeToString([]byte(password)))

	bodystr := strings.TrimSpace(r.Form.Encode())
	request, err := http.NewRequest("POST", url, strings.NewReader(bodystr))
	if err != nil {
		log.Fatal("POST Failed")
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Connection", "Keep-Alive")

	resp, err := http.DefaultClient.Do(request)
	if err != nil || resp.StatusCode != 200 {
		log.Println("[WARNING] Wifi Login Failed")
		fmt.Println(resp.StatusCode)
		fmt.Println(err)
		log.Fatal("[ERROR] An error occurred in sending request")
	}
}

func wifiLogout() {
	url := "http://a.suda.edu.cn/index.php/index/logout"
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		log.Println("[WARNING] Wifi Logout Failed")
		fmt.Println(resp.StatusCode)
		fmt.Println(err)
		log.Fatal("[ERROR] An error occurred in sending request")
	}
}
