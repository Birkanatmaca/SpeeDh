package email

import (
	"fmt"
	"log"      // Hata durumunda programı durdurmak için
	"net/smtp" // Go'nun standart SMTP kütüphanesi
	"os"
)

// Mailer, e-posta gönderme işlemleri için bir arayüzdür (sözleşmedir).
// Bu arayüzü tutmak iyi bir pratiktir. main.go'nun, e-postanın tam olarak nasıl
// gönderildiğini bilmesine gerek kalmaz, sadece Send metodu olan bir şeye
// ihtiyacı olduğunu bilir. Bu, kodunu gelecekte esnek tutar.
type Mailer interface {
	Send(to, subject, body string) error
}

// gmailMailer, Mailer arayüzünü uygular ve Gmail SMTP bilgilerini tutar.
type gmailMailer struct {
	host      string
	port      string
	fromEmail string
	password  string
	auth      smtp.Auth
}

// NewGmailMailer, .env dosyasından okuduğu bilgilerle yeni bir gmailMailer oluşturur.
func NewGmailMailer() Mailer {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	email := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	// YENİ YAKLAŞIM: Eğer .env bilgileri eksikse, uyarı verip devam etmek yerine,
	// bu kritik bir yapılandırma hatası olduğu için programı durduruyoruz.
	// Bu "fail-fast" (hata varsa hemen dur) prensibidir ve hataları anında fark etmemizi sağlar.
	if host == "" || port == "" || email == "" || password == "" {
		log.Fatal("KRİTİK HATA: SMTP ortam değişkenleri (.env) ayarlanmamış. Lütfen .env dosyasını kontrol edin.")
	}

	// SMTP sunucusunda kimlik doğrulaması yapmak için bir auth nesnesi oluşturuyoruz.
	auth := smtp.PlainAuth("", email, password, host)

	return &gmailMailer{
		host:      host,
		port:      port,
		fromEmail: email,
		password:  password,
		auth:      auth,
	}
}

// Send, Gmail SMTP sunucusu üzerinden gerçek bir e-posta gönderir.
func (m *gmailMailer) Send(to, subject, body string) error {
	// E-postanın başlıklarını (From, To, Subject) ve gövdesini doğru formatta hazırlıyoruz.
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", m.fromEmail, to, subject, body)

	// Sunucu adresi ve portunu birleştiriyoruz (örn: "smtp.gmail.com:587").
	addr := fmt.Sprintf("%s:%s", m.host, m.port)

	// smtp.SendMail fonksiyonu, tüm bu bilgileri kullanarak e-postayı gönderir.
	err := smtp.SendMail(addr, m.auth, m.fromEmail, []string{to}, []byte(message))
	if err != nil {
		return err
	}

	fmt.Printf(">>> E-posta başarıyla gönderildi: %s\n", to)
	return nil
}
