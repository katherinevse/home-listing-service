package notifier

import (
	"app/internal/kafka"
	"app/internal/repository/model"
	"fmt"
	"net/smtp"
	"os"
)

type Logger interface {
	Info(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
}

type Notifier struct {
	logger Logger
}

func NewNotifier(logger Logger) *Notifier {
	return &Notifier{logger: logger}
}

// SendNotification отправляет уведомление на электронную почту пользователя о новой квартире.
func (n *Notifier) SendNotification(user model.User, notification kafka.NotificationMessage) error {
	message := fmt.Sprintf(
		"Hello %s, a new apartment (Flat Number: %d) is available in house %d",
		user.Email, notification.FlatNumber, notification.HouseID)

	subject := "New Apartment Available!"

	err := n.sendMessage(user.Email, subject, message)
	if err != nil {
		n.logger.Error("Failed to send notification", "user", user.Email, "error", err)
		return fmt.Errorf("failed to send notification to %s: %v", user.Email, err)
	}

	n.logger.Info("Notification sent successfully", "user", user.Email, "message", message)
	return nil
}

// sendMessage отправляет электронное письмо по SMTP на указанный адрес электронной почты.
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
		n.logger.Error("Error sending email", "to", email, "error", err)
		return fmt.Errorf("error sending email: %v", err)
	}

	n.logger.Info("Email sent successfully", "to", email)
	return nil
}
