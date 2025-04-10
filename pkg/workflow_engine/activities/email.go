package activities

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	workflowengine "github.com/forkbombeu/didimo/pkg/workflow_engine"
	"gopkg.in/gomail.v2"
)

type SendMailActivity struct{}

func (a *SendMailActivity) Configure(ctx context.Context, input *workflowengine.ActivityInput) error {
	SMTPHost := os.Getenv("SMTP_HOST")
	if SMTPHost == "" {
		return errors.New("SMTP_HOST environment variable not set")
	}
	input.Config["smtp_host"] = SMTPHost
	SMTPPort := os.Getenv("SMTP_PORT")
	if SMTPPort == "" {
		return errors.New("SMTP_PORT environment variable not set")
	}
	input.Config["smtp_port"] = SMTPPort
	sender := os.Getenv("MAIL_SENDER")
	if sender == "" {
		return errors.New("MAIL_SENDER environment variable not set")
	}
	input.Config["sender"] = sender
	return nil
}

func (a *SendMailActivity) Execute(ctx context.Context, input workflowengine.ActivityInput) (workflowengine.ActivityResult, error) {
	var result workflowengine.ActivityResult

	m := gomail.NewMessage()
	m.SetHeader("From", input.Config["sender"])
	m.SetHeader("To", input.Config["recipient"])
	m.SetHeader("Subject", input.Payload["subject"].(string))
	m.SetBody("text/html", input.Payload["email"].(string))

	// Attach any files if necessary
	attachments, ok := input.Payload["attachments"].(map[string][]byte)
	if ok {
		for filename, attachedBytes := range attachments {
			attached := gomail.SetCopyFunc(func(w io.Writer) error {
				_, err := w.Write(attachedBytes)
				return err
			})
			m.Attach(filename, attached)
		}
	}

	SMTPPort, err := strconv.Atoi(input.Config["smtp_port"])
	if err != nil {
		return workflowengine.Fail(&result, "SMTP_PORT environment variable not an integer")
	}

	d := gomail.NewDialer(
		input.Config["smtp_host"],
		SMTPPort,
		os.Getenv("MAIL_USERNAME"),
		os.Getenv("MAIL_PASSWORD"),
	)

	if err := d.DialAndSend(m); err != nil {
		return workflowengine.Fail(&result, fmt.Sprintf("failed to send email: %v", err))
	}

	result.Output = "Email sent successfully"
	return result, nil
}
