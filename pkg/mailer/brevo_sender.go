package mailer

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"
)

const defaultSMTPPort = 587

// BrevoSender sends OTP emails using Brevo SMTP (compatible with net/smtp).
type BrevoSender struct {
	Host    string
	Port    int
	User    string
	Pass    string
	From    string
	Subject string
	TTL     time.Duration
	AppName string
}

func NewBrevoSenderFromEnv() (*BrevoSender, error) {
	port := defaultSMTPPort
	if v := os.Getenv("SMTP_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			port = p
		}
	}

	host := os.Getenv("SMTP_HOST")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")

	if host == "" || pass == "" || from == "" {
		return nil, fmt.Errorf("smtp credentials not configured")
	}
	if user == "" {
		user = "apikey"
	}

	subject := os.Getenv("SMTP_SUBJECT")
	if subject == "" {
		subject = "Your Registration OTP"
	}

	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "YourApp"
	}

	ttl := parseDurationEnv([]string{"OTP_TTL"}, 5*time.Minute)

	return &BrevoSender{
		Host:    host,
		Port:    port,
		User:    user,
		Pass:    pass,
		From:    from,
		Subject: subject,
		TTL:     ttl,
		AppName: appName,
	}, nil
}

func (s *BrevoSender) SendOTP(to, code string) error {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	auth := smtp.PlainAuth("", s.User, s.Pass, s.Host)

	msg := buildOTPMessage(s.From, to, s.Subject, s.AppName, code, s.TTL)

	return smtp.SendMail(addr, auth, extractEmail(s.From), []string{to}, msg)
}

func buildOTPMessage(from, to, subject, appName, code string, ttl time.Duration) []byte {
	minutes := int(ttl.Minutes())
	if minutes <= 0 {
		minutes = 5
	}

	textBody := fmt.Sprintf(
		"Your registration OTP code is: %s\nThis code expires in %d minutes.\nIf you did not request this, please ignore this email.\n",
		code,
		minutes,
	)
	htmlBody := fmt.Sprintf(`<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>%s OTP</title>
</head>
<body style="margin:0;padding:0;background:#f6f7fb;">
  <table role="presentation" cellpadding="0" cellspacing="0" width="100%%" style="background:#f6f7fb;padding:24px 0;">
    <tr>
      <td align="center">
        <table role="presentation" cellpadding="0" cellspacing="0" width="600" style="max-width:600px;background:#ffffff;border-radius:12px;overflow:hidden;font-family:Arial,Helvetica,sans-serif;">
          <tr>
            <td style="padding:24px 32px;background:#0f172a;color:#ffffff;">
              <div style="font-size:18px;font-weight:bold;">%s</div>
              <div style="font-size:12px;opacity:.8;">Registration verification</div>
            </td>
          </tr>
          <tr>
            <td style="padding:32px;color:#111827;">
              <div style="font-size:16px;margin-bottom:12px;">Hi,</div>
              <div style="font-size:14px;line-height:1.6;margin-bottom:18px;">
                Use the OTP below to complete your registration. This code expires in <strong>%d minutes</strong>.
              </div>
              <div style="font-size:28px;letter-spacing:6px;font-weight:bold;background:#f3f4f6;padding:16px 20px;border-radius:10px;display:inline-block;">
                %s
              </div>
              <div style="font-size:12px;color:#6b7280;margin-top:18px;">
                If you did not request this, please ignore this email.
              </div>
            </td>
          </tr>
          <tr>
            <td style="padding:16px 32px;background:#f9fafb;color:#6b7280;font-size:11px;">
              %s
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`,
		appName,
		appName,
		minutes,
		code,
		appName,
	)

	boundary := "otp-boundary"

	var buf bytes.Buffer
	buf.WriteString("From: " + from + "\r\n")
	buf.WriteString("To: " + to + "\r\n")
	buf.WriteString("Subject: " + subject + "\r\n")
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: multipart/alternative; boundary=" + boundary + "\r\n\r\n")

	buf.WriteString("--" + boundary + "\r\n")
	buf.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
	buf.WriteString(textBody + "\r\n")
	buf.WriteString("--" + boundary + "\r\n")
	buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
	buf.WriteString(htmlBody + "\r\n")
	buf.WriteString("--" + boundary + "--")

	return buf.Bytes()
}

func extractEmail(from string) string {
	start := strings.IndexByte(from, '<')
	end := strings.IndexByte(from, '>')
	if start >= 0 && end > start {
		return from[start+1 : end]
	}
	return from
}

func parseDurationEnv(keys []string, fallback time.Duration) time.Duration {
	for _, key := range keys {
		value := strings.TrimSpace(os.Getenv(key))
		if value == "" {
			continue
		}
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
		if i, err := strconv.Atoi(value); err == nil {
			return time.Duration(i) * time.Second
		}
	}
	return fallback
}
