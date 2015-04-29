package mail

// Sender ...
type Sender interface {
	Send(string, map[string]interface{}) error
}

// Mailer ...
type Mailer struct {
	from      string
	transport Transporter
}

// NewMailer ...
func NewMailer(from string, transport Transporter) Mailer {
	return Mailer{
		from:      from,
		transport: transport,
	}
}
