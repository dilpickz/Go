//Dillon Ramsey
//Yara auto deploy
//This code is a rewrite from the Python auto deplopyment code.
//This code can be compiled into an executable that can run on any windows machine
//The performance speed of the Go rewrite is significantly faster, but less configurable
package main

import (
    "fmt"
	"net/http"
	"log"
	"io/ioutil"
	"crypto/tls"
	"os"
	"mime/multipart"
	"bytes"
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

func deployYara(yaraFile, addr, token string, hc *http.Client) string {
	url := fmt.Sprintf("https://%s/wsapis/v2.0.0/customioc/yara/add/common?target_type=all", addr)
	file , err := os.Open(yaraFile)
	if err != nil {
		log.Fatalln(err)
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}
	file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileWriter, err := writer.CreateFormFile("filename", yaraFile)
	if err != nil {
		log.Fatalln(err)
	}
	fileWriter.Write(fileContents)
	writer.Close()

	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Add("X-FeApi-Token", token)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response, err := hc.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	rb := string(content)
	for name, values := range request.Header {
		for _, value := range values {
			fmt.Println(name, value)
		}
	}
	return rb
}

func main () {
	var username string
	var password string
	var yaraFile string
	vx_addrs := [4]string{//insert vx addresses here}
	fe_addrs := [5]string{//insert all other devices here}
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	hc := &http.Client{Transport: customTransport}

	fmt.Print("Enter username: ")
	fmt.Scanln(&username)
	fmt.Print("Enter password: ")
	fmt.Scanln(&password)
	fmt.Print("Enter name of yara file: ")
	fmt.Scanln(&yaraFile)

	for _, addr := range vx_addrs {
		url := fmt.Sprintf("https://%s/wsapis/common/v2.0.0/auth/login?", addr)
		fe_token := serverAuth(url, username, password, hc)
		fmt.Println(deployYara(yaraFile, addr, fe_token, hc))

	}

	for _, addr := range fe_addrs {
		url := fmt.Sprintf("https://%s/wsapis/v2.0.0/auth/login?", addr)
		fe_token := serverAuth(url, username, password, hc)
		fmt.Println(deployYara(yaraFile, addr, fe_token, hc))
	}
}
