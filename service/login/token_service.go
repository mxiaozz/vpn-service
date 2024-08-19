package login

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"vpn-web.funcworks.net/cst"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model/login"
)

var TokenService = &tokenService{
	header:     gb.Config.Token.Header,
	secret:     gb.Config.Token.Secret,
	expiration: gb.Config.Token.Expiration,
}

type tokenService struct {
	header     string // 令牌自定义标识
	secret     string // 令牌秘钥
	expiration int    // 令牌有效期（默认30分钟）
}

// 从缓存中获取已登录用户
func (s *tokenService) GetLoginUser(ctx *gin.Context) (*login.LoginUser, error) {
	// 从 request 中提取 token
	tokenString, err := s.getToken(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "从 request 中获取 token 失败")
	}
	if tokenString == "" {
		return nil, errors.New("token is empty")
	}

	// 解析并且校验 token
	token, err := s.parseToken(tokenString)
	if err != nil {
		return nil, errors.Wrap(err, "token 解析校验失败")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("claims is not MapClaims")
	}

	// 从 token 中获取用户ID
	userKey, ok := claims[cst.SYS_LOGIN_USER_KEY]
	if !ok || userKey == nil || userKey.(string) == "" {
		return nil, errors.New("userId is not in token")
	}

	// 从 redis 中获取用户信息
	info, err := gb.RedisProxy.Get(s.getTokenKey(userKey.(string)))
	if err == redis.Nil {
		return nil, errors.New("login expired")
	} else if err != nil {
		return nil, errors.Wrap(err, "从缓存中获取用户信息失败")
	} else if info == "" {
		return nil, errors.New("用户信息不正确")
	}

	// 转换为 LoginUser
	var user login.LoginUser
	if err = json.Unmarshal([]byte(info), &user); err != nil {
		return nil, errors.Wrap(err, "将用户缓存信息转换为 LoginUser 对象失败")
	}

	return &user, nil
}

func (s *tokenService) getToken(ctx *gin.Context) (string, error) {
	token := ctx.GetHeader(s.header)
	if token != "" && strings.HasPrefix(token, cst.SYS_TOKEN_PREFIX) {
		token = strings.TrimPrefix(token, cst.SYS_TOKEN_PREFIX)
	}
	if data, err := base64.RawURLEncoding.DecodeString(token); err != nil {
		return "", err
	} else {
		return string(data), nil
	}
}

func (s *tokenService) parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	}, jwt.WithoutClaimsValidation())
}

// 创建登录用户token，并将登录用户信息缓存
func (s *tokenService) CreateToken(loginUser *login.LoginUser) (string, error) {
	loginUser.Token = uuid.NewString()
	loginUser.LoginTime = time.Now().UnixMilli()

	// 缓存登录用户信息
	if err := s.RefreshToken(loginUser); err != nil {
		return "", errors.Wrap(err, "保存 token 缓存失败")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		cst.SYS_LOGIN_USER_KEY: loginUser.Token,
	})
	if tokenString, err := token.SignedString([]byte(s.secret)); err != nil {
		return "", errors.Wrap(err, "生成 token 加签失败")
	} else {
		return token.EncodeSegment([]byte(tokenString)), nil
	}
}

// 删除登录用户缓存信息
func (s *tokenService) DelLoginUser(token string) error {
	if token == "" {
		return nil
	}
	return gb.RedisProxy.Delete(s.getTokenKey(token))
}

func (s *tokenService) getTokenKey(uuid string) string {
	return cst.CACHE_LOGIN_TOKEN_KEY + uuid
}

// 校验 token 是否快过期，提前 10分钟进行续约
func (s *tokenService) VerifyToken(loginUser *login.LoginUser) {
	// 过期时间少于10分钟时续约 token
	if loginUser.ExpireTime-time.Now().UnixMilli() < 10*60*1000 {
		s.RefreshToken(loginUser)
	}
}

// 更新登录用户信息缓存保留时长
func (s *tokenService) RefreshToken(loginUser *login.LoginUser) error {
	if loginUser.LoginTime == 0 {
		loginUser.LoginTime = time.Now().UnixMilli()
	}
	loginUser.ExpireTime = time.Now().UnixMilli() + int64(s.expiration*60*1000)

	data, err := json.Marshal(loginUser)
	if err != nil {
		return errors.Wrap(err, "刷新 token 时长")
	}

	return gb.RedisProxy.SetEx(s.getTokenKey(loginUser.Token), string(data), s.expiration*60)
}
