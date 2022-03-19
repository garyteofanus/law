package auth

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/garyteofanus/law-assignment/domain"
	"net/http"
	"strings"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{
		service: service,
	}
}

type errorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"Error_description"`
}

type tokenRequest struct {
	Username     string `form:"username"`
	Password     string `form:"password"`
	GrantType    string `form:"grant_type"`
	ClientID     string `form:"client_id"`
	ClientSecret string `form:"client_secret"`
}
type tokenResponse struct {
	AccessToken  string      `json:"access_token"`
	ExpiresIn    int         `json:"expires_in"`
	TokenType    string      `json:"token_type"`
	Scope        interface{} `json:"scope"`
	RefreshToken string      `json:"refresh_token"`
}

func (h Handler) HandleAuth(c echo.Context) error {
	var tokenRequest tokenRequest
	if err := c.Bind(&tokenRequest); err != nil {
		return c.JSON(401, errorResponse{
			Error:            "invalid_request",
			ErrorDescription: "ada kesalahan masbro!",
		})
	}

	if tokenRequest.GrantType != "password" {
		return c.JSON(401, errorResponse{
			Error:            "invalid_request",
			ErrorDescription: "ada kesalahan masbro!",
		})
	}

	accessToken, refreshToken, err := h.service.AuthorizeUser(domain.User{
		ID:       tokenRequest.ClientID,
		Secret:   tokenRequest.ClientSecret,
		Username: tokenRequest.Username,
		Password: tokenRequest.Password,
	})
	if err != nil {
		return c.JSON(401, errorResponse{
			Error:            "invalid_request",
			ErrorDescription: "ada kesalahan masbro!",
		})
	}

	return c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  accessToken,
		ExpiresIn:    300,
		TokenType:    "Bearer",
		Scope:        nil,
		RefreshToken: refreshToken,
	})
}

type resourceResponse struct {
	AccessToken  string      `json:"access_token"`
	ClientID     string      `json:"client_id"`
	UserID       string      `json:"user_id"`
	FullName     string      `json:"full_name"`
	NPM          string      `json:"npm"`
	Expires      interface{} `json:"expires"`
	RefreshToken string      `json:"refresh_token"`
}

func (h Handler) HandleResource(c echo.Context) error {
	var splitToken []string
	for _, values := range c.Request().Header {
		for _, value := range values {
			if strings.HasPrefix(value, "Bearer ") {
				splitToken = strings.Split(value, "Bearer ")
				break
			}
		}
	}
	if len(splitToken) != 2 {
		return c.JSON(401, errorResponse{
			Error:            "invalid_token",
			ErrorDescription: "Token Salah masbro",
		})
	}
	reqToken := splitToken[1]
	if reqToken == "" {
		return c.JSON(401, errorResponse{
			Error:            "invalid_token",
			ErrorDescription: "Token Salah masbro",
		})
	}

	user, session, err := h.service.AccessResource(reqToken)
	if err != nil {
		return c.JSON(401, errorResponse{
			Error:            "invalid_token",
			ErrorDescription: "Token Salah masbro",
		})
	}

	return c.JSON(http.StatusOK, resourceResponse{
		AccessToken:  session.AccessToken,
		ClientID:     user.ID,
		UserID:       user.Username,
		FullName:     user.Username,
		NPM:          user.NPM,
		Expires:      nil,
		RefreshToken: session.RefreshToken,
	})
}
