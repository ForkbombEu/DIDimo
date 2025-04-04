package workflow

import (
	"time"
)

const (
	OpenIDTestTaskQueue = "OpenIDTestTaskQueue"
	LogsBaseURL         = "https://www.certification.openid.net/api/log/"
)

// EmailConfig holds the email configuration details
type EmailConfig struct {
	SMTPHost      string            `json:"smtp_host"`
	SMTPPort      int               `json:"smtp_port"`
	Username      string            `json:"username"`
	Password      string            `json:"password"`
	SenderEmail   string            `json:"sender_email"`
	ReceiverEmail string            `json:"receiver_email"`
	Subject       string            `json:"subject"`
	Body          string            `json:"body"`
	Attachments   map[string][]byte `json:"attachments"`
}

type GenerateYAMLInput struct {
	Variant  string `json:"variant"`
	Form     any    `json:"form"`
	FilePath string `json:"file_path"`
}

type StepCIRunnerInput struct {
	FilePath string
	Token    string
}

type StepCIRunnerResponse struct {
	Result map[string]any
}

type WorkflowInput struct {
	Variant  string `json:"variant"`
	Form     any    `json:"form"`
	UserMail string `json:"user_mail"`
	AppURL   string `json:"app_url"`
}

type LogWorkflowInput struct {
	AppURL   string
	RID      string
	Token    string
	Interval time.Duration
}

type WorkflowResponse struct {
	Message string
	Logs    []map[string]any
}

type LogWorkflowResponse struct {
	Logs []map[string]any
}

type GetLogsActivityInput struct {
	BaseURL string
	RID     string
	Token   string
}

type TriggerLogsUpdateActivityInput struct {
	AppURL     string
	WorkflowID string
	Logs       []map[string]any
}

type LogUpdateRequest struct {
	WorkflowID string           `json:"workflow_id"`
	Logs       []map[string]any `json:"logs"`
}

type SignalData struct {
	Success bool
	Reason  string
}
