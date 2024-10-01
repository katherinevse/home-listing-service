package notifier

import (
	"app/internal/kafka"
	"app/internal/repository/model"
	"fmt"
	"log"
	"net/smtp"
	"os"
)

type Notifier struct{}

func New() *Notifier {
	return &Notifier{}
}

// SendNotification sends a notification to the user's email about a new apartment listing
func (n *Notifier) SendNotification(user model.User, notification kafka.NotificationMessage) error {
	message := fmt.Sprintf(
		"Hello %s, a new apartment (Flat Number: %d) is available in house %d",
		user.Email, notification.FlatNumber, notification.HouseID)

	//theme
	subject := "New Apartment Available!"

	err := n.sendMessage(user.Email, subject, message)
	if err != nil {
		return fmt.Errorf("failed to send notification to %s: %v", user.Email, err)
	}

	log.Printf("Notification sent to %s: %s", user.Email, message)
	return nil
}

// sendMessage sends an email via SMTP to the provided email address
func (n *Notifier) sendMessage(email, subject, message string) error {
	smtpHost := "smtp.mail.ru"
	smtpPort := "587"
	from := "katedasha.twins@bk.ru"
	password := os.Getenv("SMTP_PASSWORD")

	if password == "" {
		return fmt.Errorf("SMTP credentials are not set")
	}

	auth := smtp.PlainAuth("", from, password, smtpHost)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", email, subject, message))

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, msg)
	if err != nil {
		log.Printf("Error sending email to %s: %v", email, err)
		return fmt.Errorf("error sending email: %v", err)
	}

	log.Printf("Email sent successfully to %s", email)
	return nil
}
