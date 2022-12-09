package audittrails

import (
	"github.com/gin-gonic/gin"
)

// Digunakan sebagai middleware untuk mengambil data dari tiap request dan response
func (repo *auditTrails) Middleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		auditTrail := NewAuditTrails(c.Writer)
		c.Writer = auditTrail.BodyLog
		auditTrail.GetRequestData(c)

		c.Next()

		auditTrail.GetResponseData()

		// Produce to Kafka or Store to DB here, object to pass --> auditTrail.Data
		auditTrail.Data.Service = repo.serviceName
		repo.Produce(auditTrail.Data)

		if c.Writer.Status() >= 500 {
			c.Set("status", 500)
			c.Set("logerror", auditTrail.Data.Remark)
		}

	}
}
