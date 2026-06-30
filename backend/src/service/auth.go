package service

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lovelymondayz/members/backend/src/config"
	"github.com/lovelymondayz/members/backend/src/models"
	"github.com/lovelymondayz/members/backend/src/repository"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService struct {
	userRepo *repository.UserRepository
	storeRepo *repository.StoreRepository
	cfg      *config.Config
	oauth2Config *oauth2.Config
}

type GoogleUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func NewAuthService(userRepo *repository.UserRepository, storeRepo *repository.StoreRepository, cfg *config.Config) *AuthService {
	oauth2Config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}

	return &AuthService{
		userRepo:     userRepo,
		storeRepo:    storeRepo,
		cfg:          cfg,
		oauth2Config: oauth2Config,
	}
}

// ─── Password ──────────────────────────────────────────

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func CheckPassword(hashed, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}

// ─── JWT ───────────────────────────────────────────────

type JWTClaims struct {
	UserID  string  `json:"user_id"`
	Role    string  `json:"role"`
	Email   string  `json:"email"`
	StoreID *string `json:"store_id,omitempty"`
	jwt.RegisteredClaims
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	claims := JWTClaims{
		UserID: user.UserID.String(),
		Role:   s.getRoleName(user.RoleID),
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Attach store_id for admin users
	store, err := s.storeRepo.FindByAdminID(user.UserID.String())
	if err == nil && store != nil {
		storeID := store.StoreID.String()
		claims.StoreID = &storeID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *AuthService) getRoleName(roleID uint) string {
	switch roleID {
	case 1:
		return models.RoleSuperAdmin
	case 2:
		return models.RoleAdmin
	default:
		return models.RoleMember
	}
}

// ─── Email/Password Auth ───────────────────────────────

func (s *AuthService) LoginWithPassword(email, password string) (*models.User, string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	if !user.IsActive {
		return nil, "", fmt.Errorf("account deactivated")
	}

	if !CheckPassword(user.Password, password) {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token")
	}

	return user, token, nil
}

func (s *AuthService) RegisterUser(name, email, password, roleStr string) (*models.User, string, error) {
	// Determine role
	roleID := uint(3) // default member
	if roleStr == "admin" {
		roleID = 2
	} else if roleStr == "super_admin" {
		roleID = 1
	}

	hashedPass, err := HashPassword(password)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password")
	}

	user := &models.User{
		RoleID:   roleID,
		Email:    email,
		Name:     name,
		Password: hashedPass,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token")
	}

	return user, token, nil
}

// ─── Google OAuth ──────────────────────────────────────

func (s *AuthService) GetGoogleAuthURL() string {
	return s.oauth2Config.AuthCodeURL(generateState())
}

func (s *AuthService) HandleGoogleCallback(code string) (*models.User, string, error) {
	// Exchange code for token
	token, err := s.oauth2Config.Exchange(nil, code)
	if err != nil {
		return nil, "", fmt.Errorf("failed to exchange token: %w", err)
	}

	// Fetch user info from Google
	client := s.oauth2Config.Client(nil, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("google returned status %d", resp.StatusCode)
	}

	var googleUser GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, "", fmt.Errorf("failed to decode google user info: %w", err)
	}

	// Check if user exists by Google ID
	user, err := s.userRepo.FindByGoogleID(googleUser.ID)
	if err == nil {
		// User exists, generate token
		jwtToken, err := s.GenerateToken(user)
		if err != nil {
			return nil, "", fmt.Errorf("failed to generate token")
		}
		return user, jwtToken, nil
	}

	// Check if email exists (link accounts)
	existingUser, err := s.userRepo.FindByEmail(googleUser.Email)
	if err == nil {
		// Link Google ID to existing account
		googleID := googleUser.ID
		existingUser.GoogleID = &googleID
		if err := s.userRepo.Update(existingUser); err != nil {
			return nil, "", fmt.Errorf("failed to link google account")
		}
		jwtToken, err := s.GenerateToken(existingUser)
		if err != nil {
			return nil, "", fmt.Errorf("failed to generate token")
		}
		return existingUser, jwtToken, nil
	}

	// Create new user
	googleID := googleUser.ID
	newUser := &models.User{
		RoleID:    3, // default to member
		Email:     googleUser.Email,
		Name:      googleUser.Name,
		GoogleID:  &googleID,
		AvatarURL: googleUser.Picture,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	jwtToken, err := s.GenerateToken(newUser)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token")
	}

	return newUser, jwtToken, nil
}

// ─── Admin Onboarding ──────────────────────────────────

func (s *AuthService) CreateAdminWithStore(name, email, password, storeName, address, phone, cardColorHex string) (*models.User, *models.Store, error) {
	// Create admin user
	adminRoleID := uint(2)
	
	// Use provided password or generate random one
	pw := password
	if pw == "" {
		pw = randomPassword()
	}
	hashedPass, _ := HashPassword(pw)

	admin := &models.User{
		RoleID:   adminRoleID,
		Email:    email,
		Name:     name,
		Password: hashedPass,
	}

	if err := s.userRepo.Create(admin); err != nil {
		return nil, nil, fmt.Errorf("failed to create admin: %w", err)
	}

	storeNameToUse := storeName
	if storeNameToUse == "" {
		storeNameToUse = name + "'s Store"
	}

	store := &models.Store{
		AdminID:      admin.UserID,
		Name:         storeNameToUse,
		Address:      address,
		Phone:        phone,
		CardColorHex: cardColorHex,
	}

	if store.CardColorHex == "" {
		store.CardColorHex = "#1E40AF"
	}

	if err := s.storeRepo.Create(store); err != nil {
		return nil, nil, fmt.Errorf("failed to create store: %w", err)
	}

	return admin, store, nil
}

// ─── Helpers ───────────────────────────────────────────

func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func randomPassword() string {
	b := make([]byte, 12)
	rand.Read(b)
	return hex.EncodeToString(b)
}
