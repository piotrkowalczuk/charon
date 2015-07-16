package lib

// Transporter ...
type Transporter interface {
	Send(from string, to string, topic string, mailBodies map[string]string) error
}
