
//Dillon Ramsey
//Coinbase API info grabber

package main

import (
    "fmt"
    "net/http"
    "os"
	"os/exec"
    "log"
    "io/ioutil"
	"encoding/json"
	"time"
)

func main () {
    CRYPTO_LIST := [3]string{"ADA-USD", "ETH-USD", "BTC-USD"}
    for {
		c := exec.Command("clear")
		c.Stdout =os.Stdout
		c.Run()
        for _, currencyPair := range CRYPTO_LIST{
            response, err := http.Get(fmt.Sprintf("https://api.coinbase.com/v2/prices/%s/buy", currencyPair))
            if err != nil {
                log.Fatalln(err)
            }
            body, err := ioutil.ReadAll(response.Body)
            if err != nil {
                log.Fatalln(err)
            }
			data := map[string]interface{}{}
			json.Unmarshal(body, &data)
			amount := data["data"].(map[string]interface{})["amount"]
			base := data["data"].(map[string]interface{})["base"]
			fmt.Println(fmt.Sprintf("%s is selling for %s", base, amount))
        }
		time.Sleep(10 * time.Second)
    }
}