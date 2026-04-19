package services

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/api/idtoken"
)

type GoogleIDTokenVerifier struct{}

func NewGoogleIDTokenVerifier() *GoogleIDTokenVerifier {
	return &GoogleIDTokenVerifier{}
}

func (v *GoogleIDTokenVerifier) VerifyIDToken(ctx context.Context, idToken, audience string) (*GoogleIdentity, error) {
	cleanToken := strings.TrimSpace(idToken)
	cleanAudience := strings.TrimSpace(audience)
	if cleanToken == "" || cleanAudience == "" {
		return nil, errors.New("google token or audience is empty")
	}

	payload, err := idtoken.Validate(ctx, cleanToken, cleanAudience)
	if err != nil {
		return nil, err
	}

	email, _ := payload.Claims["email"].(string)
	emailVerified, _ := payload.Claims["email_verified"].(bool)
	if strings.TrimSpace(email) == "" || !emailVerified {
		return nil, errors.New("google email is missing or unverified")
	}

	name, _ := payload.Claims["name"].(string)
	var avatarURL *string
	if picture, ok := payload.Claims["picture"].(string); ok {
		picture = strings.TrimSpace(picture)
		if picture != "" {
			avatarURL = &picture
		}
	}

	return &GoogleIdentity{
		Subject:   payload.Subject,
		Email:     email,
		Name:      name,
		AvatarURL: avatarURL,
	}, nil
}
