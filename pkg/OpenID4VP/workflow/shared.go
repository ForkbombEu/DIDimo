package workflow

import "github.com/forkbombeu/didimo/pkg/OpenID4VP/testdata"

const OpenIDTestTaskQueue = "OpenIDTestTaskQueue"

// EmailConfig holds the email configuration details
type EmailConfig struct {
	SMTPHost      string
	SMTPPort      int
	Username      string
	Password      string
	SenderEmail   string
	ReceiverEmail string
	Subject       string
	Body          string
	Attachments   map[string][]byte
}

type GenerateYAMLInput struct {
	Variant  string
	Form     testdata.Form
	FilePath string
}

type StepCIRunnerInput struct {
	FilePath string
	Token    string
}

type StepCIRunnerResponse struct {
	Result map[string]any
}

type WorkflowInput struct {
	Variant  string
	Form     testdata.Form
	UserMail string
	AppURL   string
}

type WorkflowResponse struct {
	Message string
}

type SignalData struct {
	Success bool
	Reason  string
}
