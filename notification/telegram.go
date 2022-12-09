package notification

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/bngbngstwnd/library-go-chub/model/request"
)

type telegram struct {
	host      string
	channelID string
	appName   string
}

func TelegramNotification(host, channelID, appName string) Notification {
	return &telegram{
		host:      host,
		channelID: channelID,
		appName:   appName,
	}
}

// Send : Digunakan untuk mengirim notifikasi ke telegram
func (tele *telegram) Send(status, message string) error {

	url := tele.host + "/broadcast"
	channelID := tele.channelID
	applicationName := tele.appName

	var payload request.TelegramNotificationPayload
	payload.Apps = applicationName
	payload.Message = message
	payload.Status = status

	var param request.TelegramNotificationRequest
	param.ChannelID = channelID
	param.Payload = payload
	s, _ := json.Marshal(param)

	res, err := http.Post(url, "application/json", bytes.NewBuffer(s))

	if err != nil {
		// e := apm.DefaultTracer.NewError(err)
		// e.Send()
		log.Println("Failed sending notification: ", err.Error())
		return err
	}
	defer res.Body.Close()

	log.Println("Telegram status:", res.Status)

	// Print the body to the stdout
	// io.Copy(os.Stdout, res.Body)

	return nil
}
