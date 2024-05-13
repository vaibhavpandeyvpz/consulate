package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"io"
	"strconv"
	"sync"
	"time"
)

var _slackClient *slack.Client
var _slackClientCreator sync.Once

func SlackClient() *slack.Client {
	_slackClientCreator.Do(func() {
		_slackClient = slack.New(GetConfig().Slack.BotToken)
	})

	return _slackClient
}

func VerifyRequestFromSlack(c *gin.Context) bool {
	var signature string
	if signature = c.GetHeader("X-Slack-Signature"); signature == "" {
		return false
	}

	var ts string
	if ts = c.GetHeader("X-Slack-Request-Timestamp"); ts == "" {
		return false
	}

	tsInt, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return false
	}

	now := time.Now()
	if now.After(time.Unix(tsInt, 0).Add(time.Minute)) {
		return false // too old
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return false
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	str := append([]byte(fmt.Sprintf("v0:%s:", ts)), body...)
	mac := hmac.New(sha256.New, []byte(GetConfig().Slack.SigningSecret))
	mac.Write(str)
	signature2 := fmt.Sprintf("v0=%s", hex.EncodeToString(mac.Sum(nil)))

	return hmac.Equal([]byte(signature2), []byte(signature))
}
