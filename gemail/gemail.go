package gemail

import (
    "crypto/tls"
    "mime"
    "net/smtp"
    "net/textproto"
    "path/filepath"

    "github.com/jordan-wright/email"
    "github.com/sirupsen/logrus"
)

// AddAttachments
// attachments (allow initial to nil)
// filename  result.json/result.xlsx
// content
func AddAttachments(attachments []*email.Attachment, filename string, content []byte) []*email.Attachment {
    return append(attachments, &email.Attachment{
        Filename:    filename,
        ContentType: mime.TypeByExtension(filepath.Ext(filename)),
        Header:      textproto.MIMEHeader{},
        Content:     content,
    })
}

// SendMail
// to, cc (send_to@xx.com   carbon_copy@xx.com)
// subject theme
// text
// attachments (use AddAttachments)
func SendMail(to, cc []string, subject string, text []byte, attachments []*email.Attachment) (err error) {
    client := email.NewEmail()
    client.From = reportUser
    client.To = to
    client.Subject = subject
    client.Cc = cc
    client.Text = text
    client.Attachments = attachments

    err = client.SendWithTLS(mailHost+":465", smtp.PlainAuth("", reportUser, mailPassword, mailHost), &tls.Config{ServerName: mailHost})
    if err != nil {
        logrus.Errorf("send email failed|from=%v|to=%v|msg=%v|err=%v", reportUser, to, string(text), err)
    } else {
        logrus.Infof("send email success|from=%v|to=%v|msg=%v|err=%v", reportUser, to, string(text), err)
    }

    return
}
