package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/VishwasMallikarjuna/luminous-uploads/model"
	"github.com/VishwasMallikarjuna/luminous-uploads/utils"
)

func ValidateToken(authHeader string) (*oidc.IDToken, error) {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", utils.AppConfig.TenantID))
	if err != nil {
		utils.Logger.Errorf("Failed to get provider: %v", err)
		return nil, err
	}

	oidcConfig := &oidc.Config{
		ClientID: utils.AppConfig.ClientID,
	}

	verifier := provider.Verifier(oidcConfig)
	token := strings.TrimPrefix(authHeader, "Bearer ")
	idToken, err := verifier.Verify(ctx, token)
	if err != nil {
		utils.Logger.Errorf("Invalid token: %v", err)
		return nil, err
	}

	return idToken, nil
}

func ValidateUploadImageToken(c echo.Context) (*model.CustomClaims, error) {
	authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
	if authHeader == "" {
		utils.Logger.Warn("Authorization header missing")
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Authorization header missing")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.ParseWithClaims(tokenString, &model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(utils.AppConfig.SecretKey), nil
	})
	if err != nil || !token.Valid {
		utils.Logger.Errorf("Invalid or expired token: %v", err)
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
	}

	if claims, ok := token.Claims.(*model.CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
}
