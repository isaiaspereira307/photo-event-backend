package services

import (
	"fmt"
	"time"

	"github.com/isaiaspereira307/photo-event-backend/configs"
	"github.com/resend/resend-go/v2"
)

var (
	emailClient *resend.Client
	logger      *configs.Logger
)

// InitializeEmailService configura o cliente de email
func InitializeEmailService() {
	logger = configs.GetLogger("email_service")

	apiKey := configs.GetEmailApiKey()
	if apiKey == "" {
		logger.Error("api_key não está definida")
		return
	}

	emailClient = resend.NewClient(apiKey)
	logger.Info("Serviço de email inicializado com sucesso")
}

// SendEmail envia um email através da API Resend
func SendEmail(to, subject, htmlContent string) error {
	params := &resend.SendEmailRequest{
		From:    "noreply@fluxojuridicos.com",
		To:      []string{to},
		Subject: subject,
		Html:    htmlContent,
	}

	_, err := emailClient.Emails.Send(params)
	if err != nil {
		logger.Errorf("Erro ao enviar email: %v", err)
		return err
	}

	logger.Infof("Email enviado com sucesso para %s", to)
	return nil
}

// Send2FAVerificationCode envia um código de verificação 2FA por email
func Send2FAVerificationCode(email, code string) error {
	subject := "Código de Verificação - Fluxo Jurídicos"
	htmlContent := fmt.Sprintf(`
        <div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e0e0e0; border-radius: 5px;">
            <h2 style="color: #333; text-align: center;">Seu Código de Verificação</h2>
            <p style="color: #555; font-size: 16px;">Para concluir seu login, utilize o código:</p>
            <div style="background-color: #f5f5f5; padding: 15px; text-align: center; font-size: 24px; font-weight: bold; letter-spacing: 5px; margin: 20px 0;">
                %s
            </div>
            <p style="color: #555; font-size: 14px;">Este código é válido por 10 minutos.</p>
            <p style="color: #555; font-size: 14px;">Se você não solicitou este código, por favor ignore este email.</p>
            <div style="text-align: center; margin-top: 30px; color: #888; font-size: 12px;">
                <p>© %d Fluxo Jurídicos. Todos os direitos reservados.</p>
            </div>
        </div>
    `, code, time.Now().Year())

	return SendEmail(email, subject, htmlContent)
}

func SendPasswordResetEmail(email, resetToken string) error {
	subject := "Redefinição de Senha - Fluxo Jurídicos"

	htmlContent := fmt.Sprintf(`
        <div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e0e0e0; border-radius: 5px;">
            <h2 style="color: #333; text-align: center;">Redefinição de Senha</h2>
            <p style="color: #555; font-size: 16px;">Você solicitou a redefinição de sua senha. Clique no botão abaixo para redefinir sua senha:</p>
            <p style="color: #555; font-size: 14px;">Este link é válido por 30 minutos.</p>
            <div style="text-align: center; margin: 25px 0;">
                <a href="https://app.fluxojuridicos.com/reset-password?token=%s" style="background-color: #3498db; color: white; padding: 15px 30px; text-decoration: none; border-radius: 5px; font-weight: bold; display: inline-block;">
                    Redefinir Senha
                </a>
            </div>
            <p style="color: #555; font-size: 14px;">Se você não conseguir clicar no botão, copie e cole o link abaixo no seu navegador:</p>
            <p style="color: #3498db; font-size: 14px; word-break: break-all;">https://app.fluxojuridicos.com/reset-password?token=%s</p>
            <p style="color: #555; font-size: 14px;">Se você não solicitou esta redefinição, por favor ignore este email e verifique a segurança de sua conta.</p>
            <div style="text-align: center; margin-top: 30px; color: #888; font-size: 12px;">
                <p>© %d Fluxo Jurídicos. Todos os direitos reservados.</p>
            </div>
        </div>
    `, resetToken, resetToken, time.Now().Year())

	return SendEmail(email, subject, htmlContent)
}
