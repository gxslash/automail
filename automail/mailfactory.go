package automail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"strings"
	"text/template"

	"golang.org/x/net/html"
)

type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}

type UserCred struct {
	User string
	Pass string
}

// prepareBody generates an HTML body for an email by parsing a template file and executing it with a name parameter.
// The contentHtmlPath parameter specifies the path to the HTML template file, and the name parameter is used as a parameter for executing the template.
// Returns the resulting HTML body as a string.
func prepareBody(contentHtmlPath string, name string) string {
	var body bytes.Buffer
	t, err := template.ParseFiles(contentHtmlPath)
	t.Execute(&body, struct{ Name string }{Name: name})

	if err != nil {
		log.Fatal(err)
	}

	return body.String()
}

// SendGMail sends an email using the Gmail SMTP server.
// The email's contents are specified by the parameters and attachments can also be included.
// The receiver parameter specifies the recipient's email address, and the user parameter is used to authenticate with the Gmail SMTP server.
// The attachmentFilePath parameter specifies the path to the file to attach to the email, and the contentHtmlPath parameter specifies the path to the HTML template file to use for the email's content.
// Returns an error if there is a problem sending the email, or nil if the email was sent successfully.
func SendGMail(receiver Receiver, user UserCred, attachmentFilePath string, contentHtmlPath string) error {
	request := Mail{
		Sender:  user.User,
		To:      []string{receiver.Email},
		Subject: getTitle(contentHtmlPath),
		Body:    prepareBody(contentHtmlPath, receiver.Name),
	}

	data := buildMail(request, attachmentFilePath)
	addr := "smtp.gmail.com:587"
	host := "smtp.gmail.com"
	auth := smtp.PlainAuth("", user.User, user.Pass, host)
	if err := smtp.SendMail(addr, auth, request.Sender, request.To, data); err != nil {
		return fmt.Errorf("failed to send email to %s: %w", receiver.Name, err)
	}

	fmt.Println("Email sent successfully")
	return nil
}

// buildMail builds the data for an email message that includes attachments.
// The mail parameter specifies the email message contents, and the filePath parameter specifies the path to the attachment file.
// Returns the data for the email message as a byte slice.
func buildMail(mail Mail, filePath string) []byte {
	var buf bytes.Buffer
	writeHeaders(&buf, mail)
	writeBody(&buf, mail)
	writeAttachment(&buf, filePath)
	return buf.Bytes()
}

// writeHeaders writes the email message headers to the provided buffer.
// The buf parameter is a pointer to a buffer that will receive the email message headers.
// The mail parameter specifies the email message contents.
func writeHeaders(buf *bytes.Buffer, mail Mail) {
	buf.WriteString(fmt.Sprintf("From: %s\r\n", mail.Sender))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";")))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", mail.Subject))
	boundary := "my-boundary-779"
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
	buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
}

// writeBody writes the email message body to the provided buffer.
// The buf parameter is a pointer to a buffer that will receive the email message body.
// The mail parameter specifies the email message contents.
func writeBody(buf *bytes.Buffer, mail Mail) {
	buf.WriteString("MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n")
	buf.WriteString(fmt.Sprintf("\r\n%s", mail.Body))
}

// writeAttachment writes an attachment to the email message.
// The buf parameter is a pointer to a buffer that will receive the attachment data.
// The filePath parameter specifies the path to the attachment file.
func writeAttachment(buf *bytes.Buffer, filePath string) {
	data := readAttachmentFile(filePath)
	b := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(b, data)
	fileName := strings.Split(filePath, "/")[2]
	buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", "my-boundary-779"))
	buf.WriteString("Content-Type: application/pdf; charset=\"utf-8\"\r\n")
	buf.WriteString("Content-Transfer-Encoding: base64\r\n")
	buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\r\n", fileName))
	buf.WriteString("Content-ID: <words.txt>\r\n\r\n")
	buf.Write(b)
	buf.WriteString("\r\n--my-boundary-779--")
}

// getTitle gets the title of an HTML document by parsing the document and searching for a title element.
// The path parameter specifies the path to the HTML file.
// Returns the title of the HTML document as a string.
func getTitle(path string) string {

	text, err := readHtmlFromFile(path)

	if err != nil {
		log.Fatal(err)
	}

	title := parseTitle(text)
	return title
}

// readAttachmentFile reads the contents of a file and returns them as a byte slice.
// The fileName parameter specifies the path to the file to read.
// Returns the contents of the file as a byte slice, or an error if the file cannot be read.
func readAttachmentFile(fileName string) []byte {

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

// parseTitle parses an HTML document and searches for a title element, returning the text of the title element.
// The text parameter is the HTML document to parse.
// Returns the text of the title element, or an empty string if no title element is found.
func parseTitle(text string) string {
	tokenizer := html.NewTokenizer(strings.NewReader(text))

	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			break
		} else if tt == html.StartTagToken {
			token := tokenizer.Token()
			if token.Data == "title" {
				return parseTitleText(tokenizer)
			}
		}
	}
	return ""
}

// parseTitleText searches for the text within a title element in an HTML document.
// The tokenizer parameter is an HTML tokenizer that has just processed a start tag for a title element.
// Returns the text within the title element as a string.
func parseTitleText(tokenizer *html.Tokenizer) string {
	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken || tt == html.EndTagToken {
			break
		} else if tt == html.TextToken {
			token := tokenizer.Token()
			return token.Data
		}
	}
	return ""
}

// readHtmlFromFile reads an HTML file from disk and returns its contents as a string.
// The fileName parameter specifies the path to the HTML file to read.
// Returns the contents of the HTML file as a string, or an error if the file cannot be read.
func readHtmlFromFile(fileName string) (string, error) {

	bs, err := ioutil.ReadFile(fileName)

	if err != nil {
		return "", err
	}

	return string(bs), nil
}
