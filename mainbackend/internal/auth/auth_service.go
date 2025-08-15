package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"mainbackend/internal/model"
	"mainbackend/internal/platform/email"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// generateVerificationCode, 6 haneli rastgele bir sayısal kod üretir.
func generateVerificationCode() (string, error) {
	codeBytes := make([]byte, 6)
	if _, err := rand.Read(codeBytes); err != nil {
		return "", err
	}
	for i := 0; i < 6; i++ {
		codeBytes[i] = '0' + (codeBytes[i] % 10)
	}
	return string(codeBytes), nil
}

type AuthService interface {
	RegisterUser(firstName, lastName, email, password string) error
	VerifyEmail(email, code string) (string, error)
	LoginUser(email, password string) (string, error)
	ResendVerificationCode(email string) error
	RequestPasswordReset(email string) error      // Düzeltme: Fonksiyon adı daha açıklayıcı oldu.
	ResetPassword(code, newPassword string) error // Düzeltme: E-posta parametresi kaldırıldı.
	GetJWTSecret() []byte                         // Middleware için gerekli metot.
}

type authService struct {
	userRepo       UserRepository
	mailer         email.Mailer
	jwtSecret      []byte
	tokenExpiresIn time.Duration
}

func NewAuthService(userRepo UserRepository, mailer email.Mailer, secret string) AuthService {
	return &authService{
		userRepo:       userRepo,
		mailer:         mailer,
		jwtSecret:      []byte(secret),
		tokenExpiresIn: 24 * time.Hour,
	}
}

func (s *authService) RegisterUser(firstName, lastName, email, password string) error {
	if len(password) < 8 {
		return errors.New("şifre en az 8 karakter olmalıdır")
	}

	existingUser, err := s.userRepo.GetUserByEmail(email)
	if err == nil {
		if existingUser.IsActive {
			return errors.New("bu e-posta adresi zaten kullanımda")
		}

		if existingUser.LockedUntil != nil && time.Now().Before(*existingUser.LockedUntil) {
			return fmt.Errorf("çok fazla deneme yapıldı, hesabınız %v tarihine kadar kilitlidir", existingUser.LockedUntil.Format("02-01-2006 15:04"))
		}

		// KURAL: Doğrulanmamış ve 24 saati geçmiş kayıtları sil.
		if !existingUser.IsActive && time.Now().After(existingUser.CreatedAt.Add(24*time.Hour)) {
			if err := s.userRepo.DeleteUserByID(existingUser.ID); err != nil {
				return err
			}
			// Eski kayıt silindi, yeni kaydın oluşturulması için akışa devam et.
		} else {
			return errors.New("bu e-posta ile zaten bir kayıt oluşturulmuş, lütfen hesabınızı doğrulayın veya yeni kod isteyin")
		}
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	code, _ := generateVerificationCode()
	expires := time.Now().Add(15 * time.Minute)

	newUser := &model.User{
		FirstName:           firstName,
		LastName:            lastName,
		Email:               email,
		Password:            string(hashedPassword),
		IsActive:            false,
		VerificationCode:    &code,
		VerificationExpires: &expires,
		VerificationSends:   1, // İlk gönderim
	}

	if err := s.userRepo.CreateUser(newUser); err != nil {
		return err
	}

	subject := fmt.Sprintf("Merhaba %s, SpeeDh Hesabınızı Doğrulayın", firstName)
	body := fmt.Sprintf("SpeeDh hesabınızı doğrulamak için kodunuz: %s\nKod 15 dakika içinde geçersiz olacaktır.", code)
	return s.mailer.Send(newUser.Email, subject, body)
}

func (s *authService) VerifyEmail(email, code string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("geçersiz doğrulama bilgileri")
	}
	if user.IsActive {
		return "", errors.New("hesap zaten aktif")
	}
	if user.VerificationCode == nil || *user.VerificationCode != code {
		return "", errors.New("geçersiz doğrulama kodu")
	}
	if time.Now().After(*user.VerificationExpires) {
		return "", errors.New("doğrulama kodunun süresi dolmuş, lütfen yeni kod isteyin")
	}
	if err := s.userRepo.ActivateUser(user); err != nil {
		return "", err
	}
	return s.generateJWT(user)
}

func (s *authService) ResendVerificationCode(email string) error {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil // E-posta yoksa güvenlik için hata döndürme
	}
	if user.IsActive {
		return errors.New("bu hesap zaten aktif")
	}
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		return fmt.Errorf("çok fazla deneme yapıldı, hesabınız %v tarihine kadar kilitlidir", user.LockedUntil.Format("02-01-2006 15:04"))
	}

	// KURAL: 3'ten fazla kod gönderilmesini engelle
	if user.VerificationSends >= 3 {
		lockTime := time.Now().Add(24 * time.Hour)
		if err := s.userRepo.LockUserAccount(user, lockTime); err != nil {
			return err
		}
		return fmt.Errorf("çok fazla deneme yapıldı, hesabınız 24 saatliğine kilitlendi")
	}

	// Sayaç 3'ten az ise normal işlem yap
	if err := s.userRepo.IncrementVerificationSends(user); err != nil {
		return err
	}

	code, _ := generateVerificationCode()
	expires := time.Now().Add(15 * time.Minute)
	user.VerificationCode = &code
	user.VerificationExpires = &expires
	if err := s.userRepo.SetVerificationCode(user); err != nil {
		return err
	}

	subject := "Yeni SpeeDh Doğrulama Kodunuz"
	body := fmt.Sprintf("Merhaba %s,\nİsteğiniz üzerine yeni doğrulama kodunuz: %s", user.FirstName, code)
	return s.mailer.Send(user.Email, subject, body)
}

func (s *authService) LoginUser(email, password string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("e-posta veya şifre hatalı")
	}
	if !user.IsActive {
		return "", errors.New("hesap aktif değil, lütfen e-postanızı doğrulayın")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("e-posta veya şifre hatalı")
	}
	return s.generateJWT(user)
}

func (s *authService) RequestPasswordReset(email string) error {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil // Güvenlik
	}

	code, _ := generateVerificationCode()
	expires := time.Now().Add(15 * time.Minute)
	user.VerificationCode = &code
	user.VerificationExpires = &expires

	if err := s.userRepo.SetVerificationCode(user); err != nil {
		return err
	}

	subject := "SpeeDh Şifre Sıfırlama Kodu"
	body := fmt.Sprintf("Merhaba %s,\nŞifrenizi sıfırlamak için kodunuz: %s", user.FirstName, code)
	return s.mailer.Send(user.Email, subject, body)
}

// DÜZELTİLDİ: Bu fonksiyon artık e-posta yerine doğrudan kodu kullanıyor.
func (s *authService) ResetPassword(code, newPassword string) error {
	if len(newPassword) < 8 {
		return errors.New("yeni şifre en az 8 karakter olmalıdır")
	}

	user, err := s.userRepo.GetUserByVerificationCode(code)
	if err != nil {
		return errors.New("geçersiz veya kullanılmış sıfırlama kodu")
	}

	if user.VerificationExpires == nil || time.Now().After(*user.VerificationExpires) {
		return errors.New("sıfırlama kodunun süresi dolmuş")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	return s.userRepo.UpdatePassword(user.ID, string(hashedPassword))
}

func (s *authService) generateJWT(user *model.User) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(s.tokenExpiresIn).Unix(),
		"iat": time.Now().Unix(),
	})
	return claims.SignedString(s.jwtSecret)
}
func (s *authService) GetJWTSecret() []byte {
	return s.jwtSecret
}
