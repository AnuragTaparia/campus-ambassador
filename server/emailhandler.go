package server

import(
	"fmt"
	"net/smtp"
	"os"
)

func SendCode(code,email string) error{
	fmt.Println("Sending code now")
	sender := os.Getenv("DEV_EMAIL")
	pwd := os.Getenv("DEV_PWD")


	toList := []string{email}

	from := fmt.Sprintf("From: <%s>\r\n", sender)
	to := fmt.Sprintf("To: <%s>\r\n", email)
	subject := "Password Reset Code\r\n"
	body := fmt.Sprintf("Password verification code: %s \nThis code is valid for 5 minutes",code)

	host := "smtp.gmail.com"
	port := 587

	msg := from+to+subject+"\r\n"+body

	msgBody := []byte(msg)
	auth := smtp.PlainAuth("", sender, pwd, host)

	loc := fmt.Sprintf("%s:%d",host,port)
	err := smtp.SendMail(loc, auth, sender, toList, msgBody)
	
	if err!=nil{
		return err
	}
	fmt.Println("Reset code has been sent")

	return nil 
}