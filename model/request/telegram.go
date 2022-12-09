package request

// TelegramNotificationRequest : Parameter notifikasi ke Telegram
type TelegramNotificationRequest struct {
	ChannelID string                      `json:"channel_id"`
	Payload   TelegramNotificationPayload `json:"payload"`
}

// TelegramNotificationPayload : Parameter payload notifikasi ke Telegram
type TelegramNotificationPayload struct {
	Message string `json:"message"`
	Apps    string `json:"apps"`
	Status  string `json:"status"`
}
