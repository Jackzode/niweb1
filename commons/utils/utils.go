package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/Jackzode/painting/commons/constants"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
	"time"
	"unicode"
)

func EncryptPassword(Pass string) (string, error) {
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(Pass), bcrypt.DefaultCost)
	// This encrypted string can be saved to the database and can be used as password matching verification
	return string(hashPwd), err
}

func GenerateTraceId() string {
	newUUID, _ := uuid.NewUUID()
	return newUUID.String()
}

func GetTraceIdFromHeader(ctx *gin.Context) string {
	trace := ctx.GetHeader(constants.TraceID)
	return trace
}

func ParseToken(tokenString string) (*CustomClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaim{}, func(token *jwt.Token) (i interface{}, e error) {
		return secret, nil
	})
	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token == nil || !token.Valid {
		return nil, TokenInvalid
	}
	if claims, ok := token.Claims.(*CustomClaim); ok {
		return claims, nil
	}
	return nil, TokenInvalid
}

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("token invalid")
	secret           = []byte("lawyer")
)

type CustomClaim struct {
	jwt.RegisteredClaims
	UserName string
	Role     int
	Uid      string
}

func CreateToken(UserName, Uid string, Role int) (string, error) {

	claims := CustomClaim{
		UserName: UserName,
		Uid:      Uid,
		Role:     Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"lawyer"},                               // 受众
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)),                // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(constants.ExpireDate)), // 过期时间 7天  配置文件
			Issuer:    constants.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func GetUidFromTokenByCtx(ctx *gin.Context) string {
	v, ok := ctx.Get(constants.TokenClaim)
	if !ok {
		return ""
	}
	claim := v.(*CustomClaim)
	return claim.Uid
}

func IsChinese(str string) bool {
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			return true
		}
	}
	return false
}

func UsernameSuffix() string {
	bytes := make([]byte, 2)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

var (
	usernameReg = regexp.MustCompile(`^[a-z0-9._-]{4,30}$`)
)

func IsInvalidUsername(username string) bool {
	return !usernameReg.MatchString(username)
}

func ConvertUserStatus(status, mailStatus int) string {
	switch status {
	case 1:
		if mailStatus == constants.EmailStatusToBeVerified {
			return constants.UserInactive
		}
		return constants.UserNormal
	case 9:
		return constants.UserSuspended
	case 10:
		return constants.UserDeleted
	}
	return constants.UserNormal
}

func IsReservedUsername(username string) bool {
	return false
}

func IsUsersIgnorePath(username string) bool {
	return false
}

func JsonObj2String(data interface{}) string {
	marshal, _ := json.Marshal(data)
	return string(marshal)
}

func FromJsonString2Obj(data string, obj interface{}) error {
	return json.Unmarshal([]byte(data), obj)
}

func ExtractToken(ctx *gin.Context) (token string) {
	token = ctx.GetHeader("Authorization")
	if len(token) == 0 {
		token = ctx.Query("Authorization")
	}
	return strings.TrimPrefix(token, "lawyer-")
}

func GetLang(ctx *gin.Context) string {
	acceptLanguage := ctx.GetHeader(constants.AcceptLanguageFlag)
	if len(acceptLanguage) == 0 {
		return constants.DefaultLanguage
	}
	return acceptLanguage
}

func EnShortID(id string) string {

	return id
}

func DeShortID(id string) string {

	return id
}

func IsNotZeroString(s string) bool {
	return len(s) > 0 && s != "0"
}
