package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// SendOTPEmail sends the OTP via Gmail SMTP.
// It returns an error if SMTP is not configured or sending fails.
func SendOTPEmail(toEmail, otp string) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")

	if port == "" {
		port = "587"
	}
	if from == "" {
		from = user
	}

	// 🔍 ALWAYS log what we see in env, even if wrong
	log.Printf("SMTP debug: host=%q port=%q user=%q from=%q", host, port, user, from)

	// ❌ If critical envs are missing → return error (NO more silent fallback)
	if host == "" || user == "" || pass == "" {
		return fmt.Errorf("SMTP not configured correctly (missing HOST/USER/PASS)")
	}

	auth := smtp.PlainAuth("", user, pass, host)

	subject := "Your Crypto Wallet OTP"
	body := fmt.Sprintf("Your one-time password (OTP) is: %s\n\nIt will expire in 10 minutes.", otp)

	msg := []byte(
		"To: " + toEmail + "\r\n" +
			"From: " + from + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			body + "\r\n",
	)

	addr := host + ":" + port
	if err := smtp.SendMail(addr, auth, from, []string{toEmail}, msg); err != nil {
		log.Printf("❌ failed to send OTP email to %s: %v", toEmail, err)
		return fmt.Errorf("smtp send error: %w", err)
	}

	log.Printf("✅ OTP email sent to %s", toEmail)
	return nil
}
