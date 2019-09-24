package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type (
	WebhookTestSuite struct {
		suite.Suite
		r *http.Request
		w *httptest.ResponseRecorder
	}
)

const (
	secret       = "client_secret"
	id           = "webhook_test"
	event  Event = EventOrderCreated
)

func TestWebhookTestSuite(t *testing.T) {
	suite.Run(t, new(WebhookTestSuite))
}

func (suite *WebhookTestSuite) SetupTest() {
	body := &Body{
		ID:    id,
		Event: event,
	}
	buf, _ := json.Marshal(body)
	reader := bytes.NewReader(buf)
	suite.r = httptest.NewRequest("POST", "/", reader)
	suite.r.Header.Add("Content-Type", "application/json; charset=UTF-8")

	// TODO chech it
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(fmt.Sprintf("%d.%s", body.Created, body.ID)))
	sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	suite.r.Header.Add("X-Ecwid-Webhook-Signature", sig)

	suite.w = httptest.NewRecorder()
}

// ////////////////////////////////////////////////////////////////////////////

func (suite *WebhookTestSuite) TestOK() {

	webhook := New(secret).Add(event, func(body *Body) error {
		suite.Equal(id, body.ID)
		return nil
	})

	webhook.ServeHTTP(suite.w, suite.r)

	suite.Equal(http.StatusOK, suite.w.Code)
}

// TODO add more tests
