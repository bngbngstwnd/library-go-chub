package notification

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Notification interface {
	Send(status, message string) error
	Middleware() gin.HandlerFunc
}

// Digunakan sebagai middleware untuk mengirim notifikasi ke Telegram jika terjadi error
func (tele *telegram) Middleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		c.Next()

		status := c.GetInt("status")
		if status >= 500 {
			message := c.GetString("logerror")
			go tele.Send("Error", fmt.Sprintf("%v | %v", status, message))
		}
	}
}
