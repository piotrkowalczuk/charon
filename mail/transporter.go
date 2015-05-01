package mail

// Transporter ...
type Transporter interface {
	Send(string, string, string, map[string]string) error
}
