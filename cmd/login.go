package cmd

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"example.com/notion-go/util"

	"github.com/brianstrauch/spotify"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

const RedirectURI = "http://localhost:19024/callback"

const ClientID = "7791b2d111994560b40987bc9088060f"

func NewLoginCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Log in to Spotify.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			token, err := login()
			if err != nil {
				return err
			}
			fmt.Println(token)

			if err := util.SaveToken(token); err != nil {
				return err
			}
			fmt.Println("save token success?")

			api, err := util.Authenticate()
			if err != nil {
				return err
			}

			user, err := api.GetUserProfile()
			if err != nil {
				return err
			}

			cmd.Printf("Logged in as %s.\n", user.DisplayName)

			return nil
		},
	}
}

func login() (*util.Token, error) {
	// 1. Create the code verifier and challenge
	verifier, challenge, err := util.CreatePKCEVerifierAndChallenge()
	if err != nil {
		return nil, err
	}

	// 2. Construct the authorization URI
	state, err := generateRandomState()
	if err != nil {
		return nil, err
	}

	scopes := []string{
		spotify.ScopePlaylistReadPrivate,
		spotify.ScopeUserLibraryModify,
		spotify.ScopeUserModifyPlaybackState,
		spotify.ScopeUserReadPlaybackState,
	}

	uri := util.BuildPKCEAuthURI(util.ClientID, RedirectURI, challenge, state, scopes...)

	// 3. Your app redirects the user to the authorization URI
	if err := browser.OpenURL(uri); err != nil {
		return nil, err
	}

	code, err := listenForCode(state)
	if err != nil {
		return nil, err
	}

	// 4. Your app exchanges the code for an access token
	token, err := util.RequestPKCEToken(util.ClientID, code, RedirectURI, verifier)
	if err != nil {
		return nil, err
	}

	return token, err
}

func listenForCode(state string) (code string, err error) {
	server := &http.Server{Addr: ":19024"}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state || r.URL.Query().Get("error") != "" {
			err = errors.New("login failed")
			fmt.Fprintln(w, "failed!")
		} else {
			code = r.URL.Query().Get("code")
			fmt.Fprintln(w, "login!")
		}

		// Use a separate thread so browser doesn't show a "No Connection" message
		go func() {
			_ = server.Shutdown(context.Background())
		}()
	})

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return "", err
	}

	return
}

func generateRandomState() (string, error) {
	buf := make([]byte, 7)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	state := hex.EncodeToString(buf)
	return state, nil
}
