package helper

import (
	"bytes"
	"html/template"
	"path"

	"gopkg.in/gomail.v2"
)

type FormatEmailLeave struct {
	LeaveType string
	StartDate string
	EndDate   string
	Status    *string
	Employee  string
}

// duog cufu wpty cezi
func SendEmail(leave FormatEmailLeave) {
	m := gomail.NewMessage()
	m.SetHeader("From", "arsyatteguh@gmail.com")
	m.SetHeader("To", "arsatteguh@gmail.com")
	// m.SetAddressHeader("Cc", "dan@example.com", "Dan")

	// get HTML
	var filepath = path.Join("HTML", "index.html")
	t := TemplateEmail(filepath, leave)

	m.SetHeader("Subject", "PT. Lorem")
	m.SetBody("text/html", t.String())
	// m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer("smtp.gmail.com", 587, "arsyatteguh@gmail.com", "duog cufu wpty cezi")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func TemplateEmail(url string, leave FormatEmailLeave) *bytes.Buffer {

	body := new(bytes.Buffer)
	t, err := template.ParseFiles(url)
	if err != nil {
		panic(err.Error())
	}

	var data = map[string]interface{}{
		"Name":   leave.Employee,
		"Cuti":   leave.LeaveType,
		"Start":  leave.StartDate,
		"End":    leave.EndDate,
		"Status": leave.Status,
	}

	t.Execute(body, data)

	return body
}
