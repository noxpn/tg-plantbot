package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"tg-bot-echo/lib/e"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

)

const (
	cmdPic    = "/pic"
	cmdUptime = "/uptime"
	ver = "0.0.6"
)

var startTime time.Time

func uptime() time.Duration {
    return time.Since(startTime)
}

func init() {
    startTime = time.Now()
}

func main() {

	var cfg Config

	if err := cfg.readConfigFile("conf/plantbot.cfg"); err!=nil {
		log.Fatalln(err)
	}

    logfile, err := os.OpenFile(cfg.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatal(err)
    }
	
	log.SetOutput(logfile)
	fmt.Printf("Started: GoPiBot v%s,  timeout:%d\n", ver, cfg.PollTimeout)
	log.Printf("Started: GoPiBot v%s,  timeout:%d\n", ver, cfg.PollTimeout)
    fmt.Println("logfile: ",logfile.Name())
	//os.Exit(5)

	

	offset := 0

	for {
		log.Println("CHK: ", offset)
		updates, err := getUpdates(&cfg, offset)
		if err != nil {
			log.Println("getUpdates: ", err.Error())
		}
		for _, u := range updates {
			if err := process(u, &cfg); err != nil {
				log.Println(err)
			}
			offset = u.UpdateID + 1
		}

		time.Sleep(time.Duration(cfg.PollTimeout) * time.Second)
	}

}
// change Message Field in update
func changeUpdateMessage(u Update, t string) Update {
	update := Update{
		Message: Message{
			Chat: Chat{Id: u.Message.Chat.Id},
			Text: t}}
	return update
}

// get updates fom telegram server
func getUpdates(cfg *Config, offset int) ([]Update, error) {

	head_url:= makeUrl(cfg, 0, "getUpdates")
	full_url := head_url.String() + "?offset=" + strconv.Itoa(offset)
	resp, err := http.Get(full_url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var resResponse RestResponse
	if err := json.Unmarshal(body, &resResponse); err != nil {
		return nil, err
	}
	if len(resResponse.Result) > 0 {
		log.Println("GOT: ", len(resResponse.Result))
		for _, v := range resResponse.Result {
			if cfg.Debug {
				logMessage(v)
			}
			
			log.Printf("< %v\n",v)
		}
	}

	return resResponse.Result, nil
}

// send text respond to telegram server
func sendTextRespond(cfg *Config, u Update) (string, error) {
	data, err := makeTextRespond(u)
	if err != nil {
		return "", err
	}

	url := makeUrl(cfg, 0, "sendMessage")
	res, err := http.Post(url.String(), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	log.Println(">", string(body))
	
	return string(body), nil
}

// process incoming message
func process(update Update, cfg *Config) error {
	incomingMessage := strings.Split(update.Message.Text, " ")
	if len(incomingMessage) < 2 {
		return fmt.Errorf("process update: no password")
	}
	switch incomingMessage[0] {
	case cmdPic:
		if incomingMessage[1] != cfg.BotPwd {
			return fmt.Errorf("process update: bad password")
		}
		photo, err := getPhoto(cfg.PictureDir, cfg.FailbackDir)
		if err!= nil {
			upd := changeUpdateMessage(update, err.Error())
			rpl, err := sendTextRespond(cfg, upd) 
			if err != nil {
				return err
			}
			log.Println(rpl)
			return err
		}
		u := makeUrl(cfg, update.Message.Chat.Id, "sendPhoto")
		rpl, err := sendP(u, photo)
		if err != nil {
			log.Println(err)	
		}
		log.Println("<", rpl)
	case cmdUptime:
		if incomingMessage[1] != cfg.BotPwd {
			return fmt.Errorf("process update: bad password")
		}
		str := fmt.Sprintf("bot uptime: %s", uptime().String())
		update.Message.Text = str
		rpl, err := sendTextRespond(cfg, update)
		if err != nil {
			return err
		}
		log.Println("<", rpl)
		return err
	
	default:
		if incomingMessage[1] != cfg.BotPwd {
			return fmt.Errorf("process update: bad password")
		}
		str := "\nunknown command:\n\n"
		update.Message.Text = str + update.Message.Text
		rpl, err := sendTextRespond(cfg, update)
		if err != nil {
			return err
		}
		log.Println("<", rpl)
		return err
	}
	return nil
}

// convert json:update -> []byte
func makeTextRespond(u Update) ([]byte, error) {
	var respond BotRespond
	respond.ChatID = u.Message.Chat.Id
	respond.Text = u.Message.Text
	data, err := json.Marshal(respond)
	if err != nil {
		return nil, e.Wrap("MakeTextRespond: ", err)
	}
	return data, nil
}


func makeUrl(cfg *Config, chat_id int, method string) *url.URL {
	
	u, err := url.Parse(cfg.BotUrl + cfg.Key)
	if err != nil {
		fmt.Println("failed to parse destination url: %w", err)
	}
	u.Path = path.Join(u.Path, method)
	q := u.Query()
	if chat_id != 0 {
		q.Set("chat_id", strconv.Itoa(chat_id))
		u.RawQuery = q.Encode()
	}
	
	return u
}
