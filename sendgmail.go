package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"strings"
)

type account struct {
	Username string
	Password string
}

type mailcontent struct {
	Recipients string
	Cc         string
	Bcc        string
	Subject    string
	Message    string
}

var (
	gmailAccountFile = flag.String("account", "gmailinfo",
		"File that containing account name and password.")
	gmailContentFile = flag.String("content", "mailcontent",
		"File that containing recipients, subject, and message.")
	dry = flag.Bool("dryrun", false,
		"Do not send mail, just show what will happen")
)

var gmailAccount account
var gmailContent mailcontent

func read_gmailinfo() {
	c, err := ioutil.ReadFile(*gmailAccountFile)
	if err != nil {
		fmt.Printf("failed to read mail info file: %s\n", err)
		return
	}
	if err := json.Unmarshal(c, &gmailAccount); err != nil {
		fmt.Printf("failed to unmarshal mail info: %s\n", err)
		return
	}
}

func save_gmailinfo() {
	bytes, err := json.Marshal(gmailAccount)
	if err != nil {
		fmt.Printf("failed to marshal account: %s\n", err)
		return
	}

	if err := ioutil.WriteFile(*gmailAccountFile, bytes, 0600); err != nil {
		fmt.Printf("failed to write account info: %s\n", err)
		return
	}
}

func read_gmailContent() {
	c, err := ioutil.ReadFile(*gmailContentFile)
	if err != nil {
		fmt.Printf("failed to read mail content file: %s\n", err)
		return
	}
	if err := json.Unmarshal(c, &gmailContent); err != nil {
		fmt.Printf("failed to unmarshal mail content: %s\n", err)
		return
	}
}

func sendgmail(sender string, receipients, cc, bcc, subject, message string) {
	if *dry {
		fmt.Printf("sender: %s\nrecipients: %s\ncc: %s\nbcc: %s\nsubject: %s\nmessage: %s\n",
			sender, receipients, cc, bcc, subject, message)
		return
	}
	username := gmailAccount.Username
	password := gmailAccount.Password
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
	read_gmailContent()
	sendgmail("sendgmail", gmailContent.Recipients, gmailContent.Cc,
		gmailContent.Bcc, gmailContent.Subject, gmailContent.Message)
	save_gmailinfo()
}
