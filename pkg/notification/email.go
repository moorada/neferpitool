package notification

import (
	"bytes"
	"html/template"
	"net/smtp"

	"github.com/moorada/neferpitool/pkg/log"
)

var auth smtp.Auth

func EmailChanges(tpl TemplateData, req Request) (err error) {

	auth = smtp.PlainAuth("", req.From, req.Password, "smtp.gmail.com")

	if len(tpl.DatasWhois) == 0 {
		tpl.TextWhois = "There aren't whois changes"
	}
	if len(tpl.DatasStatus) == 0 {
		tpl.TextStatus = "There aren't status changes"
	}

	err = req.parseTemplate("./config/emailTemplates/table.html", tpl)

	if err != nil {
		err = req.parseTemplateString(tpl)
	}
	if err != nil {
		log.Fatal(err.Error())
	} else {
		ok, err := req.sendEmail()
		if !ok {
			log.Error("Error To sendEmail, %s\n, %s\n, %s", err.Error(), req.From, req.Password)
		} else {
			log.Debug("Email sent To %s, Subject: %s", req.To, req.Subject)
		}
	}
	return
}

func EmailReport(tpl TemplateData, req Request) (err error) {

	auth = smtp.PlainAuth("", req.From, req.Password, "smtp.gmail.com")

	if len(tpl.DatasWhois) == 0 {
		tpl.TextWhois = "No whois changes"
	}
	if len(tpl.DatasStatus) == 0 {
		tpl.TextStatus = "No status changes"
	}
	if len(tpl.DatasExpiry) == 0 {
		tpl.TextExpiry = "No Typo-domains in expiration"
	}

	err = req.parseTemplate("./config/emailTemplates/report.html", tpl)

	if err != nil {
		err = req.parseTemplateString(tpl)
	}
	if err != nil {
		log.Fatal(err.Error())
	} else {
		ok, err := req.sendEmail()
		if !ok {
			log.Error("Error To sendEmail, %s\n, %s\n, %s", err.Error(), req.From, req.Password)
		} else {
			log.Debug("Email sent To %s, Subject: %s", req.To, req.Subject)
		}
	}
	return
}

func (r *Request) sendEmail() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.Subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)
	addr := "smtp.gmail.com:587"

	if err := smtp.SendMail(addr, auth, r.From, r.To, msg); err != nil {
		return false, err
	}
	return true, nil
}

func (r *Request) parseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)

	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}

func (r *Request) parseTemplateString(data interface{}) error {

	tmpl, err := template.New("test").Parse(templateString)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	if err = tmpl.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil

}

const (
	pathSecret = "./config/secret.json"
)
