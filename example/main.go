package main

import (
	"fmt"
	"log"

	"gxslash.com/automail"
)

func main() {

	nameIndex := 1  // stores the column index which has names
	emailIndex := 3 // stores the column index which has emails
	receiversFilePath := "participants.xlsx"
	templatePath := "./content.html"

	user := "mymail@gmail.com"
	password := "gofindthis"

	userCred := automail.UserCred{
		User: user,
		Pass: password,
	}

	receivers, e := automail.GetReceivers(receiversFilePath, nameIndex, emailIndex)

	if e != nil {
		fmt.Println("Receivers are not found")
	}

	for index, receiver := range receivers {
		fmt.Println(fmt.Sprintf("%d index: Sending to %s with email %s", index, receiver.Name, receiver.Email))
		attachmentFilePath := fmt.Sprintf("./appendix/%s.pdf", receiver.Name)
		e := automail.SendGMail(receiver, userCred, attachmentFilePath, templatePath)
		if e != nil {
			log.Fatal(e)
		}
	}

}
