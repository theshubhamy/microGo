package account

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/segmentio/ksuid"
)

type SessionData struct {
	UserID      string `json:"userId"`
	IP          string `json:"ip"`
	UserAgent   string `json:"userAgent"`
	Fingerprint string `json:"fingerprint"`
	CreatedAt   int64  `json:"createdAt"`
}

type Service interface {
	PostAccount(ctx context.Context, name, email, phone, password string) (*Account, error)
	LoginAccount(ctx context.Context, emailorphone, password, ip, userAgent string) (*Account, string, string, error)
	GetAccountById(ctx context.Context, id string) (*Account, error)
	GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
	LogoutBySession(ctx context.Context, userID, sessionID string) error
	LogoutAllSessions(ctx context.Context, userID string) error
	ListActiveSessions(ctx context.Context, userID string) ([]SessionData, error)
}

type accountService struct {
	repository  Repository
	redisClient *redis.Client
}

func NewService(r Repository, redisClient *redis.Client) Service {
	return &accountService{r, redisClient}
}

func (as *accountService) PostAccount(ctx context.Context, name, email, phone, password string) (*Account, error) {
	password_hash, err := HashPassword(password)
	if err != nil {
		log.Println("Error hashing password:", err)
	}

	account := &Account{
		Name:     name,
		Email:    email,
		Phone:    phone,
		Password: password_hash,
		ID:       ksuid.New().String(),
	}

	if err := as.repository.PutAccount(ctx, *account); err != nil {
		return nil, err
	}
	return account, nil
}

func (as *accountService) LoginAccount(ctx context.Context, emailOrPhone, password, ip, userAgent string) (*Account, string, string, error) {
	queryKey, err := checkPhoneorEmail(emailOrPhone)
	if err != nil {
		log.Println(err)
	}

	account, err := as.repository.GetAccount(ctx, queryKey, emailOrPhone)
	if err != nil {
		return nil, "", "", err
	}
	if !CompareHashPassword(account.Password, password) {
		return nil, "", "", errors.New("invalid Credentials")
	}

	accessToken, refreshToken, err := GenerateJWT(account.ID)
	if err != nil {
		return nil, "", "", err
	}
	// Generate fingerprint
	fingerprint := generateFingerprint(ip, userAgent)

	// Create session ID (random UUID recommended)
	sessionID := uuid.New().String()

	// Create session data
	session := SessionData{
		UserID:      account.ID,
		IP:          ip,
		UserAgent:   userAgent,
		Fingerprint: fingerprint,
		CreatedAt:   time.Now().Unix(),
	}

	// Store session in Redis for 7 days
	sessionJSON, _ := json.Marshal(session)
	as.redisClient.Set(ctx, "session:"+sessionID, sessionJSON, 7*24*time.Hour)
	as.redisClient.SAdd(ctx, "user-sessions:"+account.ID, sessionID)

	return account, accessToken, refreshToken, nil
}

func (as *accountService) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	accounts, err := as.repository.ListAccounts(ctx, skip, take)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (as *accountService) GetAccountById(ctx context.Context, id string) (*Account, error) {
	account, err := as.repository.GetAccount(ctx, "id", id)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (as *accountService) LogoutBySession(ctx context.Context, userID, sessionID string) error {
	sessionKey := fmt.Sprintf("session:%s", sessionID)
	userSessionsKey := fmt.Sprintf("user-sessions:%s", userID)

	// Delete the session key
	err := as.redisClient.Del(ctx, sessionKey).Err()
	if err != nil {
		return fmt.Errorf("failed to delete session key: %w", err)
	}

	err = as.redisClient.SRem(ctx, userSessionsKey, sessionID).Err()
	if err != nil {
		return fmt.Errorf("failed to remove session from user set: %w", err)
	}

	return nil
}

func (as *accountService) LogoutAllSessions(ctx context.Context, userID string) error {
	userSessionsKey := fmt.Sprintf("user-sessions:%s", userID)

	sessionIDs, err := as.redisClient.SMembers(ctx, userSessionsKey).Result()
	if err != nil {
		return fmt.Errorf("failed to get user sessions: %w", err)
	}

	for _, sid := range sessionIDs {
		sessionKey := fmt.Sprintf("session:%s", sid)
		if err := as.redisClient.Del(ctx, sessionKey).Err(); err != nil {
			fmt.Printf("Warning: failed to delete session %s: %v\n", sid, err)
		}
	}

	if err := as.redisClient.Del(ctx, userSessionsKey).Err(); err != nil {
		return fmt.Errorf("failed to delete user sessions set: %w", err)
	}

	return nil
}

func (as *accountService) ListActiveSessions(ctx context.Context, userID string) ([]SessionData, error) {
	userSessionsKey := fmt.Sprintf("user-sessions:%s", userID)

	sessionIDs, err := as.redisClient.SMembers(ctx, userSessionsKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get session IDs: %w", err)
	}

	sessions := make([]SessionData, 0, len(sessionIDs))
	for _, sid := range sessionIDs {
		sessionKey := fmt.Sprintf("session:%s", sid)

		val, err := as.redisClient.Get(ctx, sessionKey).Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			return nil, fmt.Errorf("failed to get session data for %s: %w", sid, err)
		}

		var s SessionData
		if err := json.Unmarshal([]byte(val), &s); err != nil {
			continue
		}

		sessions = append(sessions, s)
	}

	return sessions, nil
}
