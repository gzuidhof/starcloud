package auth

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/gzuidhof/starcloud/starcloud/logger"
	"github.com/spf13/viper"
)

type FriendlyCaptchaVerifyResponse struct {
	Success bool     `json:"success"`
	Errors  []string `json:"errors,omitempty"`
	Details *string  `json:"details,omitempty"`
}

func VerifyFriendlyCaptcha(c *fiber.Ctx, captchaSolution string) error {
	resp, err := http.PostForm("https://api.friendlycaptcha.com/api/v1/siteverify",
		url.Values{
			"solution": {captchaSolution},
			"secret":   {viper.GetString("friendlycaptcha.api_key")},
		},
	)

	if err != nil {
		logger.GetSugarLogger(c).Errorf("FriendlyCaptcha verification request failed %v", err)
	}

	if resp.StatusCode != 200 { // Intentionally let this through, it's probably a problem in our own credentials
		logger.GetSugarLogger(c).Errorf("FriendlyCaptcha verification failed %v", resp)
		return nil
	}

	decoder := json.NewDecoder(resp.Body)
	var vr FriendlyCaptchaVerifyResponse
	err = decoder.Decode(&vr)
	if err != nil {
		logger.GetSugarLogger(c).Errorf("FriendlyCaptcha verification failed to decode response body %v", err)
		return nil
	}

	if !vr.Success {
		return c.Status(http.StatusBadRequest).SendString("Captcha invalid, please try again. You may need to refresh.")
	}

	return nil
}
