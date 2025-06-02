package jwts

import (
	"TPC-EDM-Server/global"
	"TPC-EDM-Server/models"
	"TPC-EDM-Server/models/enum"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type Claims struct {
	UserID   uint          `json:"userID"`
	Username string        `json:"username"`
	Role     enum.RoleType `json:"role"`
}
type MyClaims struct {
	Claims
	jwt.RegisteredClaims // v5版本新加的方法
}

func (m MyClaims) GetUser() (user models.UserModel, err error) {
	err = global.DB.Take(&user, m.UserID).Error
	return
}

func GenerateJWT(claims Claims) (string, error) {
	cla := MyClaims{
		claims,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(global.Config.Jwt.Expire) * time.Hour)), // 过期时间24小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                                          // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                                          // 生效时间
		},
	}
	// 使用HS256签名算法
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, cla)
	s, err := t.SignedString([]byte(global.Config.Jwt.Secret))

	return s, err
}

func ParseJwt(tokenString string) (*MyClaims, error) {
	t, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.Jwt.Secret), nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "token is malformed") {
			return nil, errors.New("token无效")
		}
		if strings.Contains(err.Error(), "token is expired") {
			return nil, errors.New("token过期")
		}
		if strings.Contains(err.Error(), "signature is invalid") {
			return nil, errors.New("token错误")
		}
		logrus.Errorf("意料之外的Token错误: %s", err.Error())
		return nil, errors.New("请登录")
	}
	if claims, ok := t.Claims.(*MyClaims); ok && t.Valid {
		return claims, nil
	}
	return nil, errors.New("断言错误")
}

func ParseTokenByGin(c *gin.Context) (*MyClaims, error) {
	token := c.GetHeader("token")
	if token == "" {
		token = c.Query("token")
	}
	return ParseJwt(token)
}

func GetClaims(c *gin.Context) *MyClaims {
	_claims, ok := c.Get("claims")
	if !ok {
		return nil
	}
	claims, ok := _claims.(*MyClaims)
	if !ok {
		return nil
	}
	return claims
}
