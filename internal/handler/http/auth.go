package handler

import (
	"errors"
	"unicode"

	"github.com/Zhiyenbek/users-auth-service/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (h *handler) TestAuth(c *gin.Context) {
	c.JSON(200, sendResponse(0, "YOU ARE AUTHORIZED", nil))
}

func (h *handler) SignIn(c *gin.Context) {
	req := &models.UserSignInRequest{}
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		h.logger.Errorf("ERROR: invalid input, some fields are incorrect: %s\n", err.Error())
		c.AbortWithStatusJSON(400, sendResponse(-1, nil, models.ErrInvalidInput))
		return
	}

	switch {
	case !validatePassword(req.Password):
		h.logger.Error("invalid password")
		c.JSON(400, sendResponse(-1, nil, models.ErrInvalidPasswordFormat))
		return
	}
	var errMsg error
	tokens, err := h.service.AuthService.Login(req)
	if err != nil {
		h.logger.Errorf("Error occurred while login: %v", err)
		switch {
		case errors.Is(err, models.ErrWrongPassword):
			errMsg = models.ErrWrongPassword
		default:
			c.JSON(500, sendResponse(-1, nil, models.ErrInternalServer))
			return
		}
		c.JSON(200, sendResponse(-1, nil, errMsg))
		return
	}
	c.SetCookie("access_token", tokens.AccessToken.TokenValue, int(tokens.AccessToken.ExpiresAt.Seconds()), "/", "swe.works", true, true)
	c.SetCookie("refresh_token", tokens.RefreshToken.TokenValue, int(tokens.RefreshToken.ExpiresAt.Seconds()), "/refresh-token", "*.swe.works", true, true)
	c.JSON(200, sendResponse(0, nil, nil))
}
func (h *handler) RefreshToken(c *gin.Context) {
	rtToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.AbortWithStatusJSON(401, sendResponse(-1, nil, models.ErrInvalidToken))
		return
	}
	tokens, err := h.service.AuthService.RefreshToken(rtToken)
	if err != nil {
		c.AbortWithStatusJSON(401, sendResponse(-1, nil, models.ErrInvalidInput))
		return
	}
	c.SetCookie("access_token", tokens.AccessToken.TokenValue, int(tokens.AccessToken.ExpiresAt.Seconds()), "/", "swe.works", true, true)
	c.SetCookie("refresh_token", tokens.RefreshToken.TokenValue, int(tokens.RefreshToken.ExpiresAt.Seconds()), "/refresh-token", "*.swe.works", true, true)

	c.JSON(200, sendResponse(0, nil, nil))
}

// validatePassword - function that validates password. Password being validated by these requirements:
// 1.Password must have upper case characters
// 2.Password must have special characters
// function returns boolean
func validatePassword(pass string) bool {
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	for _, char := range pass {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasNumber && hasSpecial
}
