package helper

import (
	"backend/models"
	"bytes"
	"html/template"
	"path"

	"gopkg.in/gomail.v2"
)

// duog cufu wpty cezi
func SendEmail(employee models.Employee) {
	m := gomail.NewMessage()
	m.SetHeader("From", "arsyatteguh@gmail.com")
	m.SetHeader("To", "arsatteguh@gmail.com")
	// m.SetAddressHeader("Cc", "dan@example.com", "Dan")

	// get HTML
	var filepath = path.Join("HTML", "index.html")
	t := TemplateEmail(filepath, employee)

	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", t.String())
	// m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer("smtp.gmail.com", 587, "arsyatteguh@gmail.com", "duog cufu wpty cezi")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func TemplateEmail(url string, employee models.Employee) *bytes.Buffer {
	leave := employee.Leave[0]

	body := new(bytes.Buffer)
	t, err := template.ParseFiles(url)
	if err != nil {
		panic(err.Error())
	}
	start := leave.StartDate.Format("2006-01-02")
	end := leave.EndDate.Format("2006-01-02")

	var data = map[string]interface{}{
		"Name":   employee.Name,
		"Cuti":   leave.LeaveType,
		"Start":  start,
		"End":    end,
		"Status": leave.Status,
	}

	t.Execute(body, data)

	return body
}
