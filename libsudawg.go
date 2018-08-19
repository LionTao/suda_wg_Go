package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Utility func to get auxiliary param for api interaction
func getWgParam() (eventvalidation, viewstate string, cookie string) {
	url := "http://wg.suda.edu.cn/indexn.aspx"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	cookie = strings.Split(res.Header["Set-Cookie"][0], ";")[0]

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		// Get viewstate and eventvalidation for interacting with wg.suda.edu.cn
		name, _ := s.Attr("name")
		if name != "" {
			if name == "__VIEWSTATE" {
				viewstate, _ = s.Attr("value")
			} else if name == "__EVENTVALIDATION" {
				eventvalidation, _ = s.Attr("value")
			}
		}
	})

	return eventvalidation, viewstate, cookie
}

// Utility func to send post form with necessary param to wg.suda.edu.cn
func wgPost(EVENTVALIDATION string, VIEWSTATE string, user account, action string, cookie string) {
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
	r.Form.Add("TextBox1", user.Username)
	r.Form.Add("TextBox2", user.Password)
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
	request.Header.Set("Cookie", cookie)
	request.Header.Set("Upgrade-Insecure-Requests", "1")
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:61.0) Gecko/20100101 Firefox/61.0")

	//fmt.Println(username)
	//fmt.Println(password)
	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	//
	if err != nil || resp.StatusCode != 200 {
		//var body []byte
		//fmt.Println(resp.StatusCode)
		//fmt.Println(err)
		//body, err = ioutil.ReadAll(resp.Body)
		//if err != nil {
		//	log.Fatal("ohno")
		//}
		////fmt.Println(string(body))
		//err = ioutil.WriteFile("./error.html", body, 0644)
		//if err != nil {
		//
		//}
		log.Fatal("[ERROR] An error occurred in sending request")
	}
	//body, err = ioutil.ReadAll(resp.Body)
	//if err != nil {
	//
	//}
	//err = ioutil.WriteFile("./res.html", body, 0644)
	//if err != nil {
	//
	//}
	////fmt.Println(string(body))

}

// Login func for wg.suda.edu.cn
func wgLogin(user account) {
	var ev, vs, cookie string

	ev, vs, cookie = getWgParam()
	wgPost(ev, vs, user, "登陆网关", cookie)
}

// Logout func for wg.suda.edu.cn
func wgLogout(user account) {
	ev, vs, cookie := getWgParam()
	wgPost(ev, vs, user, "退出网关", cookie)
}
