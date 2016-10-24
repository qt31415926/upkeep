package main

import (
	"./email"
	"encoding/json"
	"fmt"
	"github.com/howeyc/gopass"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type Site struct {
	Url            string
	Text           string
	RecipientEmail string
}
type Configuration struct {
	Email           string
	IntervalMinutes int
	Sites           []Site
}

func fetch(url string, ch chan<- string) {
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprintf("%v", err)
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("error when fetching %s: %v", url, err)
		return
	}
	ch <- fmt.Sprintf("%s", b)
}

func verify(substr string, ch <-chan string) bool {
	body := <-ch
	//	fmt.Printf("body :%s\n", body)
	if strings.Contains(body, substr) {
		return true
	} else {
		return false
	}
}

func parseConfig() Configuration {
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	decoder := json.NewDecoder(file)
	conf := Configuration{}
	err = decoder.Decode(&conf)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	return conf
}

func inputPassword(email string) []byte {
	fmt.Println(email + "'s password:")
	pass, err := gopass.GetPasswd()
	if err != nil {
		fmt.Println("Something went wrong with inputting the password: ", err)
	}
	return pass
}

func main() {
	config := parseConfig()
	pass := inputPassword(config.Email)
	ch := make(chan string)
	for _ = range time.Tick(time.Duration(config.IntervalMinutes) * time.Minute) {
		for _, site := range config.Sites {
			url := site.Url
			text := site.Text
			go fetch(url, ch)
			isValid := verify(text, ch)
			if !isValid {
				email.SmtpSend(config.Email, string(pass[:]), site.RecipientEmail, url)
			}
		}
	}
}
