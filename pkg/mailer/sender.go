package mailer

type Sender interface {
	SendOTP(to, code string) error
}
