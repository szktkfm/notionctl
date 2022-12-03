package util

import (
	secure "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/brianstrauch/spotify"
	"github.com/spf13/viper"
)

const accountsBaseURL = "https://accounts.spotify.com"

const ClientID = "7791b2d111994560b40987bc9088060f"

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func CreatePKCEVerifierAndChallenge() (string, string, error) {
	verifier, err := generateRandomVerifier()
	if err != nil {
		return "", "", err
	}

	challenge := calculateChallenge(verifier)

	return string(verifier), challenge, nil
}

func generateRandomVerifier() ([]byte, error) {
	seed, err := generateSecureSeed()
	if err != nil {
		return nil, err
	}
	rand.Seed(seed)

	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_.-~"

	verifier := make([]byte, 128)
	for i := 0; i < len(verifier); i++ {
		idx := rand.Intn(len(chars))
		verifier[i] = chars[idx]
	}

	return verifier, nil
}

func generateSecureSeed() (int64, error) {
	buf := make([]byte, 8)
	_, err := secure.Read(buf)
	if err != nil {
		return 0, err
	}

	seed := int64(binary.BigEndian.Uint64(buf))
	return seed, nil
}

func calculateChallenge(verifier []byte) string {
	hash := sha256.Sum256(verifier)
	challenge := base64.URLEncoding.EncodeToString(hash[:])
	return strings.TrimRight(challenge, "=")
}

// BuildPKCEAuthURI constructs the URI which users will be redirected to, to authorize the app.
func BuildPKCEAuthURI(clientID, redirectURI, challenge, state string, scopes ...string) string {
	q := url.Values{}
	q.Add("client_id", clientID)
	q.Add("response_type", "code")
	q.Add("redirect_uri", redirectURI)
	q.Add("code_challenge_method", "S256")
	q.Add("code_challenge", challenge)
	q.Add("state", state)
	q.Add("scope", strings.Join(scopes, " "))
	// q.Add("owner", "user")

	return accountsBaseURL + "/authorize?" + q.Encode()
}

// RequestPKCEToken allows a user to exchange an authorization code for an access token.
func RequestPKCEToken(clientID, code, redirectURI, verifier string) (*Token, error) {
	query := make(url.Values)
	query.Set("client_id", clientID)
	query.Set("grant_type", "authorization_code")
	query.Set("code", code)
	query.Set("redirect_uri", redirectURI)
	query.Set("code_verifier", verifier)
	body := strings.NewReader(query.Encode())

	return postToken(body)
}

// RefreshPKCEToken allows a user to exchange a refresh token for an access token.
func RefreshPKCEToken(refreshToken, clientID string) (*Token, error) {
	query := make(url.Values)
	query.Set("grant_type", "refresh_token")
	query.Set("refresh_token", refreshToken)
	query.Set("client_id", clientID)
	body := strings.NewReader(query.Encode())

	return postToken(body)
}

func postToken(body io.Reader) (*Token, error) {
	res, err := http.Post(accountsBaseURL+"/api/token", "application/x-www-form-urlencoded", body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	token := new(Token)
	err = json.NewDecoder(res.Body).Decode(token)

	return token, err
}

func Authenticate() (*spotify.API, error) {
	if time.Now().Unix() > viper.GetInt64("expiration") {
		if err := refresh(); err != nil {
			return nil, err
		}
	}

	token := viper.GetString("token")
	if token == "" {
		return nil, errors.New("not logged in")
	}

	return spotify.NewAPI(token), nil
}

func refresh() error {
	refresh := viper.GetString("refresh_token")

	token, err := RefreshPKCEToken(refresh, ClientID)
	if err != nil {
		return err
	}

	return SaveToken(token)
}

func SaveToken(token *Token) error {
	expiration := time.Now().Unix() + int64(token.ExpiresIn)

	viper.Set("expiration", expiration)
	viper.Set("token", token.AccessToken)
	viper.Set("refresh_token", token.RefreshToken)

	return viper.WriteConfig()
}
