package auth

import (
	"testing"
	"fmt"
	"time"
	"net/http"

	"github.com/google/uuid"
)


func TestHashing(t *testing.T) {
	tests := []struct {
		password string
	}{
		{"password12345"},
		{"astje3267ldhz!"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("Testing Password: %s", tt.password), func(t *testing.T) {
			hashStr, err := HashPassword(tt.password)
			if err != nil {
				t.Errorf("HashPassword function returned error for %s: %v", tt.password, err)
				return
			}
			err = CheckPasswordHash(tt.password, hashStr)
			if err != nil {
				t.Errorf("CheckPasswordHash function returned error for %s: %v", tt.password, err)
				return
			}
		})
	}
}


func TestJWT(t *testing.T) {
	tests := []struct {
		name string
		userID uuid.UUID
		tokenSecret string
		expiresIn time.Duration
		wantErr bool
	}{
		{
			name: "Test failure",
			userID: uuid.New(),
			tokenSecret: "TestToken123",
			expiresIn: 100 * time.Millisecond,
			wantErr: true,
		},
		{
			name: "Test success",
			userID: uuid.New(),
			tokenSecret: "AnotherToken5",
			expiresIn: 1 * time.Hour,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwt, err := MakeJWT(tt.userID, tt.tokenSecret, tt.expiresIn)
			if err != nil {
				t.Errorf("MakeJWT() error = %v", err)
			}
			time.Sleep(200 * time.Millisecond)
			gotUserID, err := ValidateJWT(jwt, tt.tokenSecret)
			if err != nil && tt.wantErr != true {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.userID != gotUserID && tt.wantErr != true {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.userID)
				return
			}
		})
	}
}


func TestGetBearerToken(t *testing.T) {	
	tests := []struct {
		name string
		headers http.Header
		wantErr bool
	}{
		{
			name: "Should Pass",
			headers: func() http.Header {
				header := http.Header{}
				header.Add("Authorization", "Bearer testBearerToken1")
				return header
			}(), 
			wantErr: false,
		},
		{
			name: "Should Fail",
			headers: func() http.Header {
				header := http.Header{}
				header.Add("NoAuthorizationKey", "No Bearer Token")
				return header
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetBearerToken(tt.headers)
			if err != nil && tt.wantErr != true {
				t.Errorf("GetBearerToken() error - %v", err)
				return
			}
		})
	}
}

