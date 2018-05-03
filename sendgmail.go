package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"strings"
)

type mailInfo struct {
	Username string
	Password string
	Recipients string
	Cc         string
	Bcc        string
	Subject    string
	Message    string
}

var (
	gmailInfoFile = flag.String("account", "gmailinfo",
		"File that containing information about the mail to send.")
	dry = flag.Bool("dryrun", false,
		"Do not send mail, just show what will happen")

	gmailInfo mailInfo
)

func read_gmailinfo() {
	c, err := ioutil.ReadFile(*gmailInfoFile)
	if err != nil {
		fmt.Printf("failed to read mail info file: %s\n", err)
		return
	}
	if err := json.Unmarshal(c, &gmailInfo); err != nil {
		fmt.Printf("failed to unmarshal mail info: %s\n", err)
		return
	}
}

func save_gmailinfo() {
	bytes, err := json.Marshal(gmailInfo)
	if err != nil {
		fmt.Printf("failed to marshal account: %s\n", err)
		return
	}

	if err := ioutil.WriteFile(*gmailInfoFile, bytes, 0600); err != nil {
		fmt.Printf("failed to write account info: %s\n", err)
		return
	}
}

func sendgmail(sender string, receipients, cc, bcc, subject, message string) {
	if *dry {
		fmt.Printf("sender: %s\nrecipients: %s\ncc: %s\nbcc: %s\nsubject: %s\nmessage: %s\n",
			sender, receipients, cc, bcc, subject, message)
		return
	}
	username := gmailInfo.Username
	password := gmailInfo.Password
	if username == "" || password == "" {
		fmt.Printf("Mail info not read\n")
		return
	}
	hostname := "smtp.gmail.com"
	port := 587
	auth := smtp.PlainAuth("", username, password, hostname)
	msg := fmt.Sprintf(
		"To: %s\r\nCC: %s\r\nBCC: %s\r\nSubject:%s\r\n\r\n%s\r\n",
		receipients, cc, bcc, subject, message)
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", hostname, port),
		auth, sender, strings.Split(receipients, ", "), []byte(msg))
	if err != nil {
		fmt.Printf("failed to send message: %s\n", err)
	}
}

func main() {
	flag.Parse()

	read_gmailinfo()
	sendgmail("sendgmail", gmailInfo.Recipients, gmailInfo.Cc,
		gmailInfo.Bcc, gmailInfo.Subject, gmailInfo.Message)
	save_gmailinfo()
}
