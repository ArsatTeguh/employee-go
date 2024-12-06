package helper

import (
	"bytes"
	"html/template"
	"path"

	"gopkg.in/gomail.v2"
)

type BodyRegister struct {
	Email    string
	Password string
	Name     string
}

func SendEmailRegister(user BodyRegister) {
	m := gomail.NewMessage()
	m.SetHeader("From", "arsyatteguh@gmail.com")
	m.SetHeader("To", "arsatteguh@gmail.com")
	// m.SetAddressHeader("Cc", "dan@example.com", "Dan")

	// get HTML
	var filepath = path.Join("HTML", "register.html")
	t := TemplateEmailRegister(filepath, user)

	m.SetHeader("Subject", "PT. Lorem")
	m.SetBody("text/html", t.String())
	// m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer("smtp.gmail.com", 587, "arsyatteguh@gmail.com", "duog cufu wpty cezi")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func TemplateEmailRegister(url string, user BodyRegister) *bytes.Buffer {
	body := new(bytes.Buffer)
	t, err := template.ParseFiles(url)
	if err != nil {
		panic(err.Error())
	}

	var data = map[string]interface{}{
		"Email":    user.Email,
		"Password": user.Password,
	}

	t.Execute(body, data)

	return body
}
