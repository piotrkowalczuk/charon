package mail

const (
	// tplNameForWelcome ...
	tplNameForWelcome      = "welcome.html"
	tplNameForWelcomeTopic = "welcome_topic.txt"
)

type welcomeMailTplParams struct {
	// Username is name of user to whom email will be sent
	Username string

	// Mail is email address of user to whom email will be sent
	Mail string
}
