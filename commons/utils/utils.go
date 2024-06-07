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

//
//var AlphanumericSet = []rune{
//	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
//	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
//	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
//}
//
//var AlphanumericIndex map[rune]int
//
//func init() {
//	AlphanumericIndex = make(map[rune]int, len(AlphanumericSet))
//	for i, ru := range AlphanumericSet {
//		AlphanumericIndex[ru] = i
//	}
//}
//
//func EnShortID(id int64, salt int64) string {
//	id = id + salt
//	var code []rune
//	for id > 0 {
//		idx := id % int64(len(AlphanumericSet))
//		code = append(code, AlphanumericSet[idx])
//		id = id / int64(len(AlphanumericSet))
//	}
//	return string(code)
//}
//func DeShortID(code string, salt int64) int64 {
//	var id int64
//	runes := []rune(code)
//	for i := len(runes) - 1; i >= 0; i-- {
//		ru := runes[i]
//		idx := AlphanumericIndex[ru]
//		id = id*int64(len(AlphanumericSet)) + int64(idx)
//	}
//	id = id - salt
//	return id
//}

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
