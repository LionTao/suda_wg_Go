package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Login func for a.suda.edu.cn
func wifiLogin(user account) {
	url := "http://a.suda.edu.cn/index.php/index/login"

	var r http.Request
	err := r.ParseForm()
	if err != nil {
		log.Fatal("[ERROR] POST Form parsing failed")
	}
	r.Form.Add("username", user.Username)
	r.Form.Add("password", base64.StdEncoding.EncodeToString([]byte(user.Password)))

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

// Logout func for a.suda.edu.cn
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
