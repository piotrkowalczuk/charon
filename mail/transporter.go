package mail

const (
	// TransporterTypeSMTP ...
	TransporterTypeSMTP = "smtp"
)

// Transporter ...
type Transporter interface {
	Send(from string, to string, topic string, textBody string, htmlBody string) error
}
