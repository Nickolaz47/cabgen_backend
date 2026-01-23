package services

// import (
// 	"fmt"
// 	"sync"

// 	"github.com/CABGenOrg/cabgen_backend/internal/config"
// 	"github.com/CABGenOrg/cabgen_backend/internal/email"
// 	"github.com/CABGenOrg/cabgen_backend/internal/logging"
// 	"github.com/CABGenOrg/cabgen_backend/internal/repository"
// )

// func SendActivationUserEmail(userToActivate string, emailSender email.EmailSender) error {
// 	admins, err := repository.UserRepo.GetAllAdminUsers()
// 	if err != nil {
// 		return err
// 	}

// 	var adminEmailConfigs []email.EmailConfig
// 	for _, a := range admins {
// 		body := `
// 		Prezados administradores,
// 		<p>Um novo usuário foi criado no CAGBen.</p>
// 		<p>Por favor, acesse o site e realize a ativação.</p>
// 		Obrigado.
// 		`
// 		if a.Email != "" {
// 			emailConfig := email.EmailConfig{
// 				Sender:    config.SenderEmail,
// 				Recipient: a.Email,
// 				Subject:   "Novo Usuário Criado: " + userToActivate,
// 				Body:      body,
// 			}
// 			adminEmailConfigs = append(adminEmailConfigs, emailConfig)
// 		}
// 	}

// 	var wg sync.WaitGroup

// 	for _, emailConfig := range adminEmailConfigs {
// 		wg.Add(1)

// 		go func(cfg email.EmailConfig) {
// 			defer wg.Done()

// 			err := email.SendEmail(cfg, emailSender)
// 			if err != nil {
// 				logging.FileLogger.Error(
// 					fmt.Sprintf("Failed to send activation user email to %s: %v", cfg.Recipient, err),
// 				)
// 			}
// 		}(emailConfig)
// 	}
// 	wg.Wait()

// 	return nil
// }
