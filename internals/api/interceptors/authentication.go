package interceptors

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"jsantdev.com/grpc_sm_api/pkg/utils"
)

func AuthenticationInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("Received request in Authentication interceptor")

	// Skip specific rpc
	skipMethods := map[string]bool{
		"/main.ExecService/Login":          true,
		"/main.ExecService/ForgotPassword": true,
		"/main.ExecService/ResetPassword":  true,
	}

	if skipMethods[info.FullMethod] {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no metadata available")
	}

	authHeader, ok := md["authorization"]
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authorization header missing")
	}

	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	token = strings.TrimSpace(token)

	isLoggedOut := utils.JwtStore.IsLoggedOut(token)

	if isLoggedOut {
		return nil, status.Error(codes.Unauthenticated, "token is invalid")
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			log.Println("Token expired")
			return nil, status.Error(codes.Unauthenticated, "token expired")
		} else if errors.Is(err, jwt.ErrTokenMalformed) {
			log.Println("Token malformed")
			return nil, status.Error(codes.Unauthenticated, "token malformed")
		}

		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	if !parsedToken.Valid {
		log.Println("Invalid token")
		return nil, status.Error(codes.Unauthenticated, "invalid expired")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Invalid token claims")
		return nil, status.Error(codes.Unauthenticated, "invalid token claims")
	}

	userId, ok := claims["uid"].(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid userid claims")
	}

	userName, ok := claims["username"].(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid username claims")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid email claims")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid role claims")
	}

	expiresAtF64, ok := claims["exp"].(float64)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid expiresAt claims")
	}
	expiresAtI64 := int64(expiresAtF64)
	expiresAt := fmt.Sprintf("%v", expiresAtI64)

	ctxValue := context.WithValue(ctx, utils.ContextKey("uid"), userId)
	ctxValue = context.WithValue(ctxValue, utils.ContextKey("username"), userName)
	ctxValue = context.WithValue(ctxValue, utils.ContextKey("email"), email)
	ctxValue = context.WithValue(ctxValue, utils.ContextKey("role"), role)
	ctxValue = context.WithValue(ctxValue, utils.ContextKey("expiresAt"), expiresAt)

	log.Println("Sending response from Authentication interceptor")
	return handler(ctxValue, req)
}
