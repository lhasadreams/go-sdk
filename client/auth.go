package client

import "fmt"

const defaultTokenExpiryTime = 3600

// authConfig representing information like key_id, secret and token
// used for authenticating requests
type authConfig struct {
	keyID      string
	secret     string
	token      string
	expiration int
}

// WithApiKeys sets the key_id and secret used to generate API access tokens
func WithApiKeys(id, secret string) Option {
	return clientFunc(func(c *client) {
		if c.auth == nil {
			c.auth = &authConfig{}
		}
		c.auth.keyID = id
		c.auth.secret = secret
	})
}

// WithToken sets the token used to authenticate the API requests
func WithToken(token string) Option {
	return clientFunc(func(c *client) {
		c.auth.token = token
	})
}

// WithExpirationTime configures the token expiration time
func WithExpirationTime(t int) Option {
	return clientFunc(func(c *client) {
		c.auth.expiration = t
	})
}

// GenerateToken generates a new access token
func (c *client) GenerateToken() (response tokenResponse, err error) {
	if c.auth.keyID == "" || c.auth.secret == "" {
		err = fmt.Errorf("unable to generate access token: auth keys missing")
		return
	}

	body, err := jsonReader(tokenRequest{c.auth.keyID, c.auth.expiration})
	if err != nil {
		return
	}

	err = c.RequestDecoder("POST", apiTokens, body, &response)
	if err != nil {
		return
	}

	if len(response.Data) > 0 {
		// @afiune how do we handle cases where there is more than one token
		c.auth.token = response.Data[0].Token
	}

	return
}

// GenerateTokenWithKeys generates a new access token with the provided keys
func (c *client) GenerateTokenWithKeys(keyID, secretKey string) (tokenResponse, error) {
	c.auth.keyID = keyID
	c.auth.secret = secretKey
	return c.GenerateToken()
}

type tokenResponse struct {
	Data    []tokenData `json:"data"`
	Ok      bool        `json:"ok"`
	Message string      `json:"message"`
}

type tokenData struct {
	ExpiresAt string `json:"expiresAt"`
	Token     string `json:"token"`
}

type tokenRequest struct {
	KeyId      string `json:"keyId"`
	ExpiryTime int    `json:"expiryTime"`
}
