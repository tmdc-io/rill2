package auth

import (
	"bytes"
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
	"os"
	"strings"
)

var ErrForbidden = status.Error(codes.Unauthenticated, "action not allowed")

type AuthorizationRequest struct {
	Token string `json:"Token"`
}

type AuthorizationResponse struct {
	Allow  bool                 `json:"allow,omitempty" binding:"required"`
	Result *AuthorizationResult `json:"result,omitempty"`
	Error  *AuthorizationError  `json:"error,omitempty"`
	Valid  bool                 `json:"valid"`
}

type AuthorizationResult struct {
	ID   string      `json:"id"`
	Data interface{} `json:"data,omitempty"`
	Tags []string    `json:"tags,omitempty"`
}

type AuthorizationError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func HeimdallAuthorization(bearerToken string) error {
	heimdallUrl := strings.Trim(os.Getenv("HEIMDALL_URL"), "/")
	if len(heimdallUrl) == 0 {
		return status.Error(codes.NotFound, "environment Variable HEIMDALL_URL is missing or empty")
	}

	client := &http.Client{}

	// expecting Bearer Token as 'Bearer token'
	token := strings.Split(bearerToken, " ")
	if len(token) != 2 {
		return status.Error(codes.Internal, "invalid token")
	}
	payload := AuthorizationRequest{
		Token: token[1],
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", heimdallUrl+"/api/v1/authorize", bytes.NewBuffer(payloadBytes))

	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer res.Body.Close()

	var ar AuthorizationResponse
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	err = json.Unmarshal(body, &ar)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	if !ar.Allow {
		return ErrForbidden
	}
	return nil
}
