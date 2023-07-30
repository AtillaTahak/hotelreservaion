package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("JWTAuthentication")
	token, ok := c.GetReqHeaders()["X-Access-Token"]
	if !ok {
		return fmt.Errorf("unauthorized")
	}
	claims, err := validateToken(token)
	if err != nil {
		return err
	}
	expiry := claims["expiry"].(float64)
	if time.Now().Unix() > int64(expiry) {
		return fmt.Errorf("unauthorized")
	}
	c.Locals("user", claims)
	return c.Next()
}

func validateToken(tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parser JWT toke", err)
		return nil, fmt.Errorf("unauthorized")
	}
	if !token.Valid {
		fmt.Println("invalid token")
		return nil, fmt.Errorf("unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("invalid token")
		return nil, fmt.Errorf("unauthorized")
	}
	return claims, nil

}
