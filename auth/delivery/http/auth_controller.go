package http

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"go-echo-api/auth"
	"go-echo-api/infrastructure/response"
	"go-echo-api/middleware"
	"go-echo-api/utils"
	"gopkg.in/go-playground/validator.v9"
	"os"
)

type authController struct {
	authRepository auth.Repository
	authMapper     *auth.Mapper
}

func NewAuthController(s auth.Repository) *authController {
	return &authController{authRepository: s,
		authMapper: auth.NewAuthMapper(),
	}
}

func (c *authController) Login(ctx echo.Context) error {
	var dto auth.LoginDto
	if err := ctx.Bind(&dto); err != nil {
		return response.InternalServerError(ctx, utils.InternalServerError, nil, err.Error())
	}
	if err := ctx.Validate(dto); err != nil {
		errorData := make(echo.Map)
		for _, v := range err.(validator.ValidationErrors) {
			errorData[v.Field()] = v.Tag()
		}
		return response.ValidationError(ctx, utils.ValidationError, nil, err.Error())
	}
	result, err := c.authRepository.Login(dto.Email)
	if err != nil {
		return response.InternalServerError(ctx, utils.InternalServerError, nil, err.Error())
	}
	if !utils.CheckPasswordHash(dto.Password, result.Password) {
		return response.BadRequest(ctx, utils.BadRequest, nil, "Wrong username or password")
	}
	tokens, refreshToken, expire, err := middleware.GenerateTokenPair(result)
	if err != nil {
		return response.InternalServerError(ctx, utils.InternalServerError, nil, err.Error())
	}
	return response.SingleData(ctx, utils.OK, echo.Map{"access_token": tokens, "refresh_token": refreshToken, "expire": expire},
		nil)
}

func (c *authController) Register(ctx echo.Context) error {
	var dto auth.RegisterDto
	if err := ctx.Bind(&dto); err != nil {
		return response.InternalServerError(ctx, utils.InternalServerError, nil, err.Error())
	}
	if err := ctx.Validate(dto); err != nil {
		errorData := make(echo.Map)
		for _, v := range err.(validator.ValidationErrors) {
			errorData[v.Field()] = v.Tag()
		}
		return response.ValidationError(ctx, utils.ValidationError, nil, err.Error())
	}
	result, err := c.authRepository.Register(dto)
	if err != nil {
		return response.InternalServerError(ctx, utils.InternalServerError, nil, err.Error())
	}
	return response.SingleData(ctx, utils.OK, c.authMapper.Map(result), nil)
}

func (c *authController) RefreshToken(ctx echo.Context) error {
	type tokenReqBody struct {
		RefreshToken string `json:"refresh_token"`
	}
	tokenReq := tokenReqBody{}
	_ = ctx.Bind(&tokenReq)

	// Parse takes the token string and a function for looking up the key.
	// The latter is especially useful if you use multiple keys for your application.
	// The standard is to use 'kid' in the head of the token to identify
	// which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenReq.RefreshToken, func(token *jwt.Token) (interface{}, error) {

		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return response.Unauthorized(ctx, utils.Unauthorized, nil, "Token not valid or expired")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get the user record from database or
		// run through your business logic to verify if the user can log in
		email := claims["email"]
		result, err := c.authRepository.Login(email.(string))
		if gorm.IsRecordNotFoundError(err) {
			return response.Unauthorized(ctx, utils.Unauthorized, nil, err.Error())
		}
		newTokenPair, newRefreshToken, newExpire, err := middleware.GenerateTokenPair(result)
		if err != nil {
			return err
		}
		return response.SingleData(ctx, utils.OK, echo.Map{
			"access_token": newTokenPair, "refresh_token": newRefreshToken, "expire": newExpire,
		}, nil)
	}

	return err
}
