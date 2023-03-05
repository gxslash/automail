package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"strings"
	"text/template"

	"github.com/xuri/excelize/v2"
)

func main() {

	receivers := getReceivers("participants.xlsx")
	for index, receiver := range receivers {
		if index > 80 {
			fmt.Println(fmt.Sprintf("Sending to %s with email %s", receiver.Name, receiver.Email))
			fmt.Println(index)
			sendGoMail(receiver)
		}
	}
}

type Receiver struct {
	Name  string
	Email string
}

type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}

func readExcelCols(filePath string) [][]string {

	f, err := excelize.OpenFile(filePath)

	if err != nil {
		log.Fatal(err)
	}

	col, err := f.GetCols("Form Responses 1")

	if err != nil {
		log.Fatal(err)
	}

	return col
}

func getReceivers(filePath string) []Receiver {

	col := readExcelCols(filePath)

	receivers := make([]Receiver, len(col[1]), len(col[1]))
	for index, name := range col[1] {
		if index == 0 {
			continue
		}
		receivers[index-1].Name = strings.TrimSpace(name)
	}

	for index, email := range col[3] {
		if index == 0 {
			continue
		}
		receivers[index-1].Email = strings.TrimSpace(email)
	}
	return receivers
}

func prepareBody(name string) string {
	var body bytes.Buffer
	templatePath := "./content.html"
	t, err := template.ParseFiles(templatePath)
	t.Execute(&body, struct{ Name string }{Name: name})

	if err != nil {
		log.Fatal(err)
	}

	return body.String()
}

func sendGoMail(receiver Receiver) {
	user := "qsbitu22@gmail.com"
	password := "ihcduscobatgsyzn"

	request := Mail{
		Sender:  "qsbitu22@gmail.com",
		To:      []string{receiver.Email},
		Subject: "QSB ITU Python for Quantum Computing Certificate",
		Body:    prepareBody(receiver.Name),
	}

	addr := "smtp.gmail.com:587"
	host := "smtp.gmail.com"

	filePath := fmt.Sprintf("./certificates/%s.pdf", strings.TrimSpace(receiver.Name))
	data := buildMail(request, filePath)
	auth := smtp.PlainAuth("", user, password, host)
	err := smtp.SendMail(addr, auth, request.Sender, request.To, data)

	if err != nil {
		fmt.Println(receiver.Name)
		log.Fatal(err)
	}
	fmt.Println("Email sent successfully")
}

func buildMail(mail Mail, filePath string) []byte {

	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("From: %s\r\n", mail.Sender))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";")))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", mail.Subject))

	boundary := "my-boundary-779"
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n",
		boundary))

	buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
	buf.WriteString("MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n")
	buf.WriteString(fmt.Sprintf("\r\n%s", mail.Body))

	buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
	buf.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	buf.WriteString("Content-Transfer-Encoding: base64\r\n")
	buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\r\n", strings.Split(filePath, "/")[2]))
	buf.WriteString("Content-ID: <words.txt>\r\n\r\n")

	data := readAttachmentFile(filePath)

	b := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(b, data)
	buf.Write(b)
	buf.WriteString(fmt.Sprintf("\r\n--%s", boundary))

	buf.WriteString("--")

	return buf.Bytes()
}

func readAttachmentFile(fileName string) []byte {

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	return data
}
