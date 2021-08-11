package config

type Config struct {
	TelegramToken   string `json:"telegram_token"`
	WebhookEndpoint string `json:"webhook_endpoint"`
	CertificateFile string `json:"certificate_file"`
	PrivateKeyFile  string `json:"private_key_file"`
	DatabasePath    string `json:"database_path"`
}
