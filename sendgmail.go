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
	gmailInfoFile = flag.String("info", "gmailinfo",
		"File that containing information about the mail to send.")
	dry = flag.Bool("dryrun", false,
		"Do not send mail, just show what will happen")
	msgFile = flag.String("msgfile", "mailmessage",
		"File that containing message of the mail")

	username = flag.String("user", "",
		"Gmail account user name")
	password = flag.String("pass", "",
		"Gmail account password")
	recipients = flag.String("recip", "",
		"Recipients of the mail")
	cc = flag.String("cc", "",
		"CC list of the mail")
	bcc = flag.String("bcc", "",
		"BCC list of the mail")
	subject = flag.String("subject", "",
		"Subject of the mail")
	message = flag.String("message", "",
		"Messge of the mail")

	gmailInfo mailInfo
)

func read_gmailinfo() bool {
	c, err := ioutil.ReadFile(*gmailInfoFile)
	if err != nil {
		fmt.Printf("failed to read mail info file: %s\n", err)
		return false
	}
	if err := json.Unmarshal(c, &gmailInfo); err != nil {
		fmt.Printf("failed to unmarshal mail info: %s\n", err)
		fmt.Printf("The original file content: %s\n", string(c))
		return false
	}
	return true
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

func readMsgfile() string {
	c, err := ioutil.ReadFile(*msgFile)
	if err != nil {
		fmt.Printf("failed to read message file: %s\n", err)
		return ""
	}
	return string(c)
}

func sendgmail(sender string, recipients, cc, bcc, subject, message string) {
	if *dry {
		fmt.Printf("sender: %s\nrecipients: %s\ncc: %s\nbcc: %s\nsubject: %s\nmessage: %s\n",
			sender, recipients, cc, bcc, subject, message)
		return
	}
	if *username == "" || *password == "" {
		fmt.Printf("Gmail account information is wrong\n")
		return
	}
	hostname := "smtp.gmail.com"
	port := 587
	auth := smtp.PlainAuth("", *username, *password, hostname)
	msg := fmt.Sprintf(
		"To: %s\r\nCC: %s\r\nBCC: %s\r\nSubject:%s\r\n\r\n%s\r\n",
		recipients, cc, bcc, subject, message)
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", hostname, port),
		auth, sender, strings.Split(recipients, ", "), []byte(msg))
	if err != nil {
		fmt.Printf("failed to send message: %s\n", err)
	}
}

func main() {
	flag.Parse()

	read := read_gmailinfo()
	if !read {
		save_gmailinfo()
		return
	}

	if *username == "" {
		*username = gmailInfo.Username
	}
	if *password == "" {
		*password = gmailInfo.Password
	}
	if *recipients == "" {
		*recipients = gmailInfo.Recipients
	}
	if *cc == "" {
		*cc = gmailInfo.Cc
	}
	if *bcc == "" {
		*bcc = gmailInfo.Bcc
	}
	if *subject == "" {
		*subject = gmailInfo.Subject
	}

	msg := readMsgfile()
	if msg == "" {
		msg = gmailInfo.Message
	}

	if *message != "" {
		msg = *message
	}
	sendgmail("sendgmail", *recipients, *cc, *bcc, *subject, msg)
}
