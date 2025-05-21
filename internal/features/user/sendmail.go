package user

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendMail(facode int) error {
	host := os.Getenv("EMAIL_SERVER_HOST")
	portStr := os.Getenv("EMAIL_SERVER_PORT")
	password := os.Getenv("EMAIL_SERVER_PASSWORD")
	user := os.Getenv("EMAIL_FROM")
	text := fmt.Sprintf(`
	<html>
		<body style="font-family: sans-serif; color: #333;">
			<h2 style="color: #2c3e50;">Здравствуйте!</h2>
			<p>Ваш код подтверждения для входа:</p>
			<h1 style="color: #2980b9;">%d</h1>
			<p>Пожалуйста, введите этот код в приложении. Он действителен в течение нескольких минут.</p>
			<br/>
			<p style="font-size: 0.9em; color: #888;">Если вы не запрашивали код, просто проигнорируйте это письмо.</p>
			<p>С уважением,<br/>Команда Mahakala</p>
		</body>
	</html>`, facode)

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("Error converting port to integer: %v", err)
		return fmt.Errorf("invalid port number: %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "test@emgushovs.ru")
	m.SetHeader("To", "emgushovs@mail.ru")
	m.SetHeader("Subject", "🔐 Код подтверждения доступа")
	m.SetBody("text/html", text)

	d := gomail.NewDialer(
		host,
		port,
		user,
		password,
	)

	d.SSL = true

	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
