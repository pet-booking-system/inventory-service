package interceptors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"invservice/internal/logger"
	"io"
	"net/http"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var protectedMethods = map[string]string{
	"/inventory.InventoryService/UpdateResourceStatus": "admin",
	"/inventory.InventoryService/CreateResource":       "admin",
	"/inventory.InventoryService/DeleteResource":       "admin",
}

type contextKey string

const UserIDKey contextKey = "userID"

type AuthResponse struct {
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	ExpiresAt string `json:"expiresAt"`
}

func validateTokenWithAuthService(token string) (*AuthResponse, error) {
	url := "http://auth-service:8080/api/v1/auth/validate"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("Failed to create request to auth service: ", err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Auth service not reachable: ", err)
		return nil, errors.New("auth service not reachable")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Info("Token validation failed with status: ", resp.StatusCode)
		return nil, fmt.Errorf("invalid token or unauthorized: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read auth service response body: ", err)
		return nil, err
	}

	var authResp AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		logger.Error("Failed to unmarshal auth response: ", err)
		return nil, err
	}

	return &authResp, nil
}

func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		requiredRole, protected := protectedMethods[info.FullMethod]
		if protected {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				logger.Error("Missing metadata in gRPC context")
				return nil, status.Error(codes.Unauthenticated, "missing metadata")
			}

			authHeaders := md.Get("authorization")
			if len(authHeaders) == 0 {
				logger.Error("Missing authorization header")
				return nil, status.Error(codes.Unauthenticated, "missing authorization header")
			}

			token := strings.TrimPrefix(authHeaders[0], "Bearer ")
			token = strings.TrimSpace(token)

			authResp, err := validateTokenWithAuthService(token)
			if err != nil {
				logger.Error("Token validation failed: ", err)
				return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("token validation failed: %v", err))
			}

			if authResp.Role != requiredRole {
				logger.Info("Permission denied: required role ", requiredRole, ", but got ", authResp.Role)
				return nil, status.Error(codes.PermissionDenied, "insufficient permissions")
			}

			ctx = context.WithValue(ctx, UserIDKey, authResp.UserID)
		}

		return handler(ctx, req)
	}
}
