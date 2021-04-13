package pkg

import (
	"encoding/json"
	"errors"
	"github.com/cristalhq/jwt"
	"gitlab.com/pbobby001/postit-schedule-status/pkg/logs"
	"net/http"
	"os"
	"time"
)

func WebSocketTokenValidateToken(tokenString string, tenantNamespace string) error {
	logs.Logger.Info(tokenString)
	logs.Logger.Info(tenantNamespace)

	jwtToken, err := jwt.Parse([]byte(tokenString))
	if err != nil {
		return err
	}

	var jwtClaims *jwt.StandardClaims
	claims := jwtToken.RawClaims()
	err = json.Unmarshal(claims, &jwtClaims)
	if err != nil {
		return err
	}

	//var newToken string
	if jwtClaims.ExpiresAt.Time().Before(time.Now()) {
		logs.Logger.Info("Token expired! getting a new one....")

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, os.Getenv("AUTHENTICATION_SERVER_URL")+"/refresh-token", nil)
		if err != nil {
			return err
		}

		req.Header.Set("token", tokenString)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode == 200 {
			logs.Logger.Info("refresh-token: ", resp.Header.Get("refresh-token"))
		} else {
			return err
		}
	}

	validator := jwt.NewValidator(
		jwt.AudienceChecker([]string{"postit-audience", tenantNamespace}),
	)

	if jwtClaims.Audience[1] != tenantNamespace {
		return errors.New("invalid tenant namespace")
	}

	err = validator.Validate(jwtClaims)
	if err != nil {
		return err
	}

	return nil
}
