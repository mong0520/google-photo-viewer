package services

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/mong0520/google-photo-viewer/models"
	"golang.org/x/oauth2"
)

type SessionService struct {
}

var sessionService *SessionService

func InitSessionService(r *gin.Engine, redisUri string) error {
	store, err := redis.NewStore(10, "tcp", redisUri, "", []byte("secret"))
	if err != nil {
		return err
	}

	r.Use(sessions.Sessions("my_session", store))
	return nil
}

func GetSessionService() *SessionService {
	return sessionService
}

func (m *SessionService) GetUserInfo(c *gin.Context) *models.UserInfo {
	session := sessions.Default(c)
	userInfoVal := session.Get("user-info")

	if userInfoVal == nil {
		return nil
	}

	// fmt.Println(userInfoVal)
	ret := userInfoVal.(*models.UserInfo)
	return ret
}

func (m *SessionService) GetOAuth2Conf(c *gin.Context) *oauth2.Config {
	session := sessions.Default(c)
	confVal := session.Get("conf")
	if confVal == nil {
		return nil
	}

	ret := confVal.(*oauth2.Config)
	return ret
}

func (m *SessionService) GetOAuth2Token(c *gin.Context) *oauth2.Token {
	session := sessions.Default(c)
	tokenVal := session.Get("token")
	if tokenVal == nil {
		return nil
	}

	token := tokenVal.(*oauth2.Token)
	return token
}

func (m *SessionService) SetSessionEntity(c *gin.Context, key interface{}, val interface{}) {
	session := sessions.Default(c)
	session.Set(key, val)
	session.Save()
}
