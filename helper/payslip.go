package helper

import (
	"io"

	"gopkg.in/gomail.v2"
)

func SendPayslip(file []byte, email string, periodic string) {

	m := gomail.NewMessage()
	m.SetHeader("From", "arsyatteguh@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Payroll")
	m.SetBody("text/plain", "Here we attach your e-slip (PDF file) for the month: "+periodic)

	// Attach PDF directly from bytes
	m.Attach("payslip.pdf", gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(file)
		return err
	}))

	d := gomail.NewDialer("smtp.gmail.com", 587, "arsyatteguh@gmail.com", "duog cufu wpty cezi")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
