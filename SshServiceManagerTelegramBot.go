package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
	"net/http"
)

type SshServiceManagerTelegramBot struct {
	TelegramBot
}

func (bot *SshServiceManagerTelegramBot) Init() {
	bot.TelegramBot = TelegramBot{}
	bot.TelegramBot.Init()

	http.HandleFunc("/newsshconnection", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			msgTxt := r.URL.Query().Get("msg")
			if bot.trustedChatID > 0 {
				msg := tgbotapi.NewMessage(bot.trustedChatID, string(msgTxt))
				bot.botApi.Send(msg)
			}
		}))

	go http.ListenAndServe("127.0.0.1:8080", nil)

	startTor := os.Getenv("START_TOR")

	if strings.EqualFold(startTor, "true") {
		log.Println("Starting Tor Service")
		torCmd := exec.Command("tor")
		torCmd.Stdout = log.Writer()
		torCmd.Stderr = log.Writer()
		go torCmd.Run()
	} else {
		torCmd := exec.Command("sh", "-c", "ps | grep tor")
		torCmd.Stdout = log.Writer()
		torCmd.Stderr = log.Writer()
		go torCmd.Run()
	}

	startSsh := os.Getenv("START_SSH")

	if strings.EqualFold(startSsh, "true") {
		log.Println("Starting SSH Service")
		sshCmd := exec.Command("/usr/sbin/sshd")
		sshCmd.Stdout = log.Writer()
		sshCmd.Stderr = log.Writer()
		go sshCmd.Run()
	}

	exists := func(path string) (bool, error) {
		_, err := os.Stat(path)
		if err == nil { return true, nil }
		if os.IsNotExist(err) { return false, nil }
		return false, err
	}

	if doesExist, _ := exists("/var/lib/tor/ssh/"); doesExist == true {
		//bot.ResetHiddenServices()
		log.Println("The Directory Exists: /var/lib/tor/ssh/")
	} else {
		log.Println("The Directory Doesn't Exist: /var/lib/tor/ssh/")
	}

	bot.TelegramBot.RegisterMessageHandler("*", func(update tgbotapi.Update){
		if strings.EqualFold(update.Message.Text, "/kill") {
			os.Exit(0)
		} else if strings.EqualFold(update.Message.Text, "/reset") {
			bot.ResetHiddenServices()
			for {
				if doesExist, _ := exists("/var/lib/tor/ssh/"); doesExist == false {
					log.Println("The Directory Doesn't Exist: /var/lib/tor/ssh/")
					log.Println("Sleep for 5 seconds...")
					time.Sleep(5 * time.Second)
				} else {
					log.Println("The Directory Exists: /var/lib/tor/ssh/")
					break
				}
			}
		}

		if doesExist, _ := exists("/var/lib/tor/ssh/"); doesExist == true {
			hostName, err := ioutil.ReadFile("/var/lib/tor/ssh/hostname")
			if err != nil {
				log.Fatal(err)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "torify ssh -l root " + string(hostName))
			bot.botApi.Send(msg)

			rndPasswd, err := ioutil.ReadFile("/root/.rndpasswd")
			if err != nil {
				log.Fatal(err)
			}
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, string(rndPasswd))
			bot.botApi.Send(msg)
		} else {
			log.Println("The Directory Doesn't Exist: /var/lib/tor/ssh/")
		}
	})


}

func (bot *SshServiceManagerTelegramBot) ResetHiddenServices() {
	log.Println("Resetting Hidden Service and Root Password")
	bot.RemoveOldServices()

	pwdCmd := exec.Command("sh", "-c", "cat /dev/urandom | env LC_CTYPE=C tr -dc 'a-zA-Z' | fold -w 26 | head -n 1 > /root/.rndpasswd")
	pwdCmd.Stdout = log.Writer()
	pwdCmd.Stderr = log.Writer()
	pwdCmd.Run()
	pwdCmd2 := exec.Command("sh", "-c", "cat /root/.rndpasswd > /root/.rndpasswdx2")
	pwdCmd2.Stdout = log.Writer()
	pwdCmd2.Stderr = log.Writer()
	pwdCmd2.Run()
	pwdCmd3 := exec.Command("sh", "-c", "cat /root/.rndpasswd >> /root/.rndpasswdx2")
	pwdCmd3.Stdout = log.Writer()
	pwdCmd3.Stderr = log.Writer()
	pwdCmd3.Run()
	pwdCmd4 := exec.Command("sh", "-c", "cat /root/.rndpasswdx2 | passwd")
	pwdCmd4.Stdout = log.Writer()
	pwdCmd4.Stderr = log.Writer()
	pwdCmd4.Run()

	torCmd := exec.Command("pkill", "-sighup", "tor")
	torCmd.Stdout = log.Writer()
	torCmd.Stderr = log.Writer()
	torCmd.Run()
}

func (bot *SshServiceManagerTelegramBot) RemoveOldServices() {
	log.Println("Removing Old Service Directory")
	rmCmd := exec.Command("rm", "-rf", "/var/lib/tor/ssh/")
	rmCmd.Stdout = log.Writer()
	rmCmd.Stderr = log.Writer()
	rmCmd.Run()
}
