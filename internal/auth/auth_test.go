package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	tests := []struct {
		name        string
		userId      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
		wantErr     bool
	}{
		{
			name:        "Generate token",
			userId:      uuid.New(),
			tokenSecret: "test123",
			expiresIn:   10 * time.Minute,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := MakeJWT(tt.userId, tt.tokenSecret, tt.expiresIn)
			t.Log(token)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	tests := []struct {
		name        string
		userId      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
		wantErr     bool
	}{
		{
			name:        "Generate token",
			userId:      uuid.New(),
			tokenSecret: "test123",
			expiresIn:   10 * time.Minute,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := MakeJWT(tt.userId, tt.tokenSecret, tt.expiresIn)
			t.Log(token)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeJWT() error = %v, wantErr %v", err, tt.wantErr)
			}

			validated, err := ValidateJWT(token, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}

			if (tt.userId != validated) != tt.wantErr {
				t.Errorf("ValidateJWT() matching error - jwt.uuid(%v) != tt.uuid(%v), wantErr %v", validated, tt.userId, tt.wantErr)
			}

		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "Extract token",
			token:   "Bearer ABCDEFG",
			wantErr: false,
		},
		{
			name:    "Extract token2",
			token:   "Bearer HUNTER2",
			wantErr: false,
		},
		{
			name:    "Missing token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "Invalid token",
			token:   "idkMaybeAHeaderOrSomething",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := http.Header{}
			header.Set("Authorization", tt.token)

			_, err := GetBearerToken(header)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
