package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"anuario/interno/model"
)

type Claims struct {
	Subject   int    `json:"sub"`
	Usuario   string `json:"usuario"`
	Role      string `json:"role"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

type JWTManager struct {
	secret    []byte
	accessTTL time.Duration
	now       func() time.Time
}

func NewJWTManager(secret string, accessTTL time.Duration) (*JWTManager, error) {
	if len(secret) < 16 {
		return nil, errors.New("JWT_SECRET deve ter pelo menos 16 caracteres")
	}
	if accessTTL <= 0 {
		return nil, errors.New("accessTTL deve ser positivo")
	}
	return &JWTManager{
		secret:    []byte(secret),
		accessTTL: accessTTL,
		now:       time.Now,
	}, nil
}

func MustJWTManager(secret string, accessTTL time.Duration) *JWTManager {
	manager, err := NewJWTManager(secret, accessTTL)
	if err != nil {
		panic(err)
	}
	return manager
}

func (m *JWTManager) Generate(admin model.Admin) (string, error) {
	now := m.now().UTC()
	claims := Claims{
		Subject:   admin.ID,
		Usuario:   admin.Usuario,
		Role:      admin.Role,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(m.accessTTL).Unix(),
	}
	return m.sign(claims)
}

func (m *JWTManager) AccessTTL() time.Duration {
	return m.accessTTL
}

func (m *JWTManager) Validate(token string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, errors.New("token JWT malformado")
	}

	signingInput := parts[0] + "." + parts[1]
	expected := m.signature(signingInput)
	got, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return Claims{}, errors.New("assinatura JWT invalida")
	}
	if subtle.ConstantTimeCompare(got, expected) != 1 {
		return Claims{}, errors.New("assinatura JWT invalida")
	}

	var header struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}
	if err := decodeSegment(parts[0], &header); err != nil {
		return Claims{}, err
	}
	if header.Alg != "HS256" || header.Typ != "JWT" {
		return Claims{}, errors.New("cabecalho JWT invalido")
	}

	var claims Claims
	if err := decodeSegment(parts[1], &claims); err != nil {
		return Claims{}, err
	}
	if claims.Subject == 0 || claims.Role == "" {
		return Claims{}, errors.New("claims JWT incompletas")
	}
	if m.now().UTC().Unix() >= claims.ExpiresAt {
		return Claims{}, errors.New("token JWT expirado")
	}
	return claims, nil
}

func (m *JWTManager) sign(claims Claims) (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	headerPart := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsPart := base64.RawURLEncoding.EncodeToString(claimsJSON)
	signingInput := headerPart + "." + claimsPart
	signature := base64.RawURLEncoding.EncodeToString(m.signature(signingInput))
	return signingInput + "." + signature, nil
}

func (m *JWTManager) signature(signingInput string) []byte {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(signingInput))
	return mac.Sum(nil)
}

func decodeSegment(segment string, dst any) error {
	data, err := base64.RawURLEncoding.DecodeString(segment)
	if err != nil {
		return fmt.Errorf("segmento JWT invalido: %w", err)
	}
	if err := json.Unmarshal(data, dst); err != nil {
		return fmt.Errorf("JSON JWT invalido: %w", err)
	}
	return nil
}
