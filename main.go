package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
  	"regexp"
	"net/http"
	"io/ioutil"
	"log"
	"bytes"
	"strings"
	"encoding/json"
  	"github.com/gookit/color"
	"github.com/bwmarrin/discordgo"
	"github.com/denisbrodbeck/machineid"
)


var T,_ = ioutil.ReadFile("token.txt")
var Token string = strings.Replace(string(T), "\n", "", 1)

var Redeemed []string

func main() {
  fmt.Println("Developed by swag")
  id, err := machineid.ProtectedID("claimer")
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Machine-ID: " + id)
  dg, err := discordgo.New(Token)
	_, err = dg.User("@me")
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func redeem_code(code string, channel_id string) {
		client := &http.Client{}
		values := map[string]string{"channel_id": channel_id, "payment_source_id": "123"}
		jsonValue, _ := json.Marshal(values)
		req, err := http.NewRequest("POST", "https://discordapp.com/api/v6/entitlements/gift-codes/"+code+"/redeem", bytes.NewBuffer(jsonValue))
		req.Header.Add("content-type", "application/json")
		req.Header.Add("Authorization", Token)
		resp, err := client.Do(req)
		defer resp.Body.Close()

		//bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
	//bodyString := string(bodyBytes)
	//fmt.Println(bodyString)
		switch status := resp.StatusCode; status {
		case 200:
			color.Success.Println("Redeemed Nitro: " + code)
		case 404:
			color.Danger.Println("Invalid Code: " + code)
		case 400:
			color.Danger.Println("Invalid Code: " + code)
		}
		Redeemed = append(Redeemed, code)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

  match1, _ := regexp.MatchString("(?:https?:)?discord(?:app.com/gifts/|.gift/)([^\\s]+)", m.Content)
  if match1 == true {
    re := regexp.MustCompile("(?:https?:)?discord(?:app.com/gifts/|.gift/)([^\\s]+)")
    match := re.FindStringSubmatch(m.Content)
	code := match[1]
	//fmt.println(m.channelID)
	_, found := Find(Redeemed, code)
	if !found && len(code) == 16 || len(code) == 24 {
		redeem_code(code, m.ChannelID)
	} else {
		fmt.Println(code, "- Already Attempted or Invalid Format")
	}
}
}
func Find(slice []string, val string) (int, bool) {
    for i, item := range slice {
        if item == val {
            return i, true
        }
    }
    return -1, false
}
