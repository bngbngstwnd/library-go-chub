package auth

import (
	"errors"
	"net/http"

	"github.com/bngbngstwnd/library-go-chub/constant"
	"github.com/bngbngstwnd/library-go-chub/model/response"
	"github.com/bngbngstwnd/library-go-chub/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Auth struct {
	User  string `json:"user" gorm:"column:user"`
	Pass  string `json:"password" gorm:"column:password"`
	Roles string `json:"roles" gorm:"column:roles"`
}

type Tabler interface {
	TableName() string
}

func (Auth) TableName() string {
	return "user_chub"
}

var authList []Auth

func Initialize(list []Auth) error {
	if len(list) == 0 {
		return errors.New("please provide list of authenticated users")
	}
	authList = list
	return nil
}

func getUser(user string) *Auth {
	for _, v := range authList {
		if v.User == user {
			return &v
		}
	}
	return nil
}

func getUserV3(conn *gorm.DB, user string) *Auth {

	var authUser Auth
	err := conn.Model(&Auth{}).Where("user = ?", user).Find(&authUser).Error
	if err != nil {
		return nil
	}

	return &authUser
}

func ResponseUnauthorized(c *gin.Context) {
	resp := response.BuildUnauthorizedResponse(constant.RESPONSE_MESSAGE_INVALID_AUTH_TOKEN, "")
	c.AbortWithStatusJSON(http.StatusUnauthorized, resp)
	c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
}

func BasicAuth() gin.HandlerFunc {

	return func(c *gin.Context) {

		user, password, hasAuth := c.Request.BasicAuth()

		// check if hasAuth
		if !hasAuth {
			ResponseUnauthorized(c)
			return
		}

		// get user
		fetchedUser := getUser(user)
		if fetchedUser == nil {
			ResponseUnauthorized(c)
			return
		}

		// check pass
		if fetchedUser.Pass != password {
			ResponseUnauthorized(c)
			return
		}
	}
}

func BasicAuth_v2() gin.HandlerFunc {

	return func(c *gin.Context) {

		user, password, hasAuth := c.Request.BasicAuth()

		// check if hasAuth
		if !hasAuth {
			ResponseUnauthorized(c)
			return
		}

		// get user
		fetchedUser := getUser(user)
		if fetchedUser == nil {
			ResponseUnauthorized(c)
			return
		}

		// check passS
		if fetchedUser.Pass != util.GetMD5Hash(password) {
			ResponseUnauthorized(c)
			return
		}
	}
}

func BasicAuth_v3(mysqlConn *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		user, password, hasAuth := c.Request.BasicAuth()

		// check if hasAuth
		if !hasAuth {
			ResponseUnauthorized(c)
			return
		}

		// get user
		fetchedUser := getUserV3(mysqlConn, user)
		if fetchedUser == nil {
			ResponseUnauthorized(c)
			return
		}

		// check passS
		if fetchedUser.Pass != util.GetMD5Hash(password) {
			ResponseUnauthorized(c)
			return
		}

		c.Set("user", fetchedUser.User)
		c.Set("roles", fetchedUser.Roles)
	}
}
