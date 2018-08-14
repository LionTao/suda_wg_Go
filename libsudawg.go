package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Logout func for wg.suda.edu.cn
func wgLogout(username string, password string) {
	ev, vs := getWgParam()
	wgPost(ev, vs, username, password, "退出网关")
}

// Login func for wg.suda.edu.cn
func wgLogin(username string, password string) {

	ev, vs := getWgParam()
	wgPost(ev, vs, username, password, "登陆网关")
}

// Utility func to get auxiliary param for api interaction
func getWgParam() (eventvalidation, viewstate string) {
	url := "http://wg.suda.edu.cn/indexn.aspx"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var VIEWSTATE, EVENTVALIDATION string
	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		// Get viewstate and eventvalidation for interacting with wg.suda.edu.cn
		name, _ := s.Attr("name")
		if name != "" {
			if name == "__VIEWSTATE" {
				VIEWSTATE, _ = s.Attr("value")
			} else if name == "__EVENTVALIDATION" {
				EVENTVALIDATION, _ = s.Attr("value")
			}
		}
	})

	return EVENTVALIDATION, VIEWSTATE
}

// Utility func to send post form with necessary param to wg.suda.edu.cn
func wgPost(EVENTVALIDATION string, VIEWSTATE string, username string, password string, action string) {
	if EVENTVALIDATION == "" || VIEWSTATE == "" {
		log.Fatal("Oh no! no param was given")
	}

	url := "http://wg.suda.edu.cn/indexn.aspx"

	var r http.Request
	err := r.ParseForm()
	if err != nil {
		log.Fatal("[ERROR] POST Form parsing failed")
	}
	r.Form.Add("__EVENTTARGET", "")
	r.Form.Add("__EVENTARGUMENT=", "")
	r.Form.Add("__VIEWSTATE", VIEWSTATE)
	r.Form.Add("__EVENTVALIDATION", EVENTVALIDATION)
	r.Form.Add("TextBox1", username)
	r.Form.Add("TextBox2", password)
	r.Form.Add("nw", "RadioButton2")
	r.Form.Add("tm", "RadioButton8")
	// Unique action for login or logout
	if action == "登陆网关" {
		r.Form.Add("Button1", action)
	} else if action == "退出网关" {
		r.Form.Add("Button4", action)
	} else {
		log.Fatal("[ERROR] Wrong POST action")
	}

	// Do urlencode and send POST
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
