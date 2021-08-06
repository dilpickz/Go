package main

import (
    "fmt"
	"net/http"
	"log"
	"io/ioutil"
	"crypto/tls"
	"encoding/json"
)

func serverAuth(url, username, password string, hc *http.Client) string {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.SetBasicAuth(username, password)
	resp, err := hc.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	token := resp.Header.Get("X-FeApi-Token")
	return token
}

func statusCheck(addr, fe_token string, hc *http.Client) string{
	url := fmt.Sprintf("https://%s/wsapis/v2.0.0/health/system", addr)
	request, err := http.NewRequest("GET", url , nil)
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set("X-FeApi-Token", fe_token)
	resp2, err := hc.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	//defer request.Body.Close()
	body, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		log.Fatalln(err)
	}
	data := map[string]interface{}{}
	json.Unmarshal(body, &data)
	status := data["status"]
	info := data[]
	applianceInfo := data["appliance"].(map[string]interface{})["name"]
	sb := fmt.Sprintf("The status for %[1]s is %[2]s", applianceInfo, status)
	return sb
}

func main() {
	var username string
	var password string
	vx_addrs := [4]string{//add vx addresses here}
	fe_addrs := [5]string{//add other devices here}
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	hc := &http.Client{Transport: customTransport}

	fmt.Print("Enter username: ")
	fmt.Scanln(&username)
	fmt.Print("Enter password: ")
	fmt.Scanln(&password)

	for _, addr := range vx_addrs {
		url := fmt.Sprintf("https://%s/wsapis/common/v2.0.0/auth/login?", addr)
		fe_token := serverAuth(url, username, password, hc)
		serverStatus := statusCheck(addr, fe_token, hc)
		fmt.Println(serverStatus)
	}

	for _, addr := range fe_addrs {
		url := fmt.Sprintf("https://%s/wsapis/v2.0.0/auth/login?", addr)
		fe_token := serverAuth(url, username, password, hc)
		serverStatus := statusCheck(addr, fe_token, hc)
		fmt.Println(serverStatus)
	}
}

