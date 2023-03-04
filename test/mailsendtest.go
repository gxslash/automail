/*
https://articles.wesionary.team/sending-emails-with-go-golang-using-smtp-gmail-and-oauth2-185ee12ab306

https://zetcode.com/golang/email-smtp/

https://dzone.com/articles/go-language-library-for-reading-and-writing-micros

https://www.youtube.com/watch?v=H0HZc4FgX7E
*/

package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"gopkg.in/gomail.v2"
)

func sendMailSimple(subject string, body string, to []string) {
	auth := smtp.PlainAuth(
		"",
		"qsbitu22@gmail.com",
		"ihcduscobatgsyzn",
		"smtp.gmail.com",
	)

	msg := "Subject: " + subject + "\n" + body

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"qsbitu22@gmail.com",
		to,
		[]byte(msg),
	)

	if err != nil {
		fmt.Println(err)
	}
}

func sendMailSimpleHTML(subject string, templatePath string, to []string) {
	var body bytes.Buffer
	t, err := template.ParseFiles(templatePath)
	t.Execute(&body, struct{ Name string }{Name: "Robby"})

	auth := smtp.PlainAuth(
		"",
		"qsbitu22@gmail.com",
		"ihcduscobatgsyzn",
		"smtp.gmail.com",
	)

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

	msg := "Subject: " + subject + "\n" + headers + "\n\n" + body.String()

	err = smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"qsbitu22@gmail.com",
		to,
		[]byte(msg),
	)

	if err != nil {
		fmt.Println(err)
	}
}

func sendGoMail(templatePath string) {
	var body bytes.Buffer
	t, err := template.ParseFiles(templatePath)
	t.Execute(&body, struct{ Name string }{Name: "Yunus"})

	if err != nil {
		fmt.Println(err)
		return
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "qsbitu22@gmail.com")
	m.SetHeader("To", "yunus.e.gunduz@gmail.com")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", body.String())
	// m.Attach("./file.png")
	d := gomail.NewDialer("smtp.gmail.com", 587, "qsbitu22@gmail.com", "ihcduscobatgsyzn")

	fmt.Println(m)

	if err := d.DialAndSend(); err != nil {
		fmt.Println("err")
		panic(err)
	}
}

func main() {

	// sendMailSimple("Subject", "Body", []string{"yunus.e.gunduz@gmail.com", "imgxslash@gmail.com"})
	// sendMailSimpleHTML("Subject", "./content.html", []string{"yunus.e.gunduz@gmail.com"})

	sendGoMail("./content.html")
}
