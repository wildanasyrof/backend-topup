package handler

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/config" // <-- TAMBAHKAN
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity" // <-- TAMBAHKAN
	"github.com/wildanasyrof/backend-topup/internal/service"
	apperror "github.com/wildanasyrof/backend-topup/pkg/apperr"
	"github.com/wildanasyrof/backend-topup/pkg/oauth"
	"github.com/wildanasyrof/backend-topup/pkg/response"
	"github.com/wildanasyrof/backend-topup/pkg/validator"
	"golang.org/x/oauth2"
)

const refreshTokenCookieName = "refresh_token"

type AuthHandler struct {
	authService service.AuthService
	userService service.UserService
	oauthPkg    oauth.GoogleOauthPkg
	devStore    oauth.DevStore
	validator   validator.Validator
	cfg         *config.Config // <--- TAMBAHKAN
}

// Modifikasi NewAuthHandler
func NewAuthHandler(
	authService service.AuthService,
	userService service.UserService,
	oauthPkg oauth.GoogleOauthPkg,
	devStore oauth.DevStore,
	validator validator.Validator,
	cfg *config.Config, // <--- TAMBAHKAN
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
		oauthPkg:    oauthPkg,
		devStore:    devStore,
		validator:   validator,
		cfg:         cfg, // <--- TAMBAHKAN
	}
}

// Helper untuk set cookie
func (h *AuthHandler) setRefreshTokenCookie(c *fiber.Ctx, session *entity.UserSession) {
	cookie := &fiber.Cookie{
		Name:     refreshTokenCookieName,
		Value:    session.ID.String(),
		Expires:  session.ExpiresAt,
		HTTPOnly: true,
		SameSite: "Strict", // Melindungi dari CSRF
		// 'Secure: true' WAJIB di produksi (HTTPS)
		Secure: h.cfg.Server.Env != "development",
		Path:   "/auth", // Hanya kirim cookie ke endpoint /auth/*
		Domain: h.cfg.Server.CookieDomain,
	}
	c.Cookie(cookie)
}

// Helper untuk clear cookie
func (h *AuthHandler) clearRefreshTokenCookie(c *fiber.Ctx) {
	cookie := &fiber.Cookie{
		Name:     refreshTokenCookieName,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Set kedaluwarsa di masa lalu
		HTTPOnly: true,
		SameSite: "Strict",
		Secure:   h.cfg.Server.Env != "development",
		Path:     "/auth",
		Domain:   h.cfg.Server.CookieDomain,
	}
	c.Cookie(cookie)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterUserRequest

	if err := c.BodyParser(&req); err != nil {
		return apperror.New(apperror.CodeBadRequest, "Invalid JSON", err)
	}

	if err := h.validator.ValidateBody(req); err != nil {
		return apperror.Validation(err)
	}

	user, err := h.authService.Register(c.UserContext(), &req)
	if err != nil {
		return err

	}
	// Tidak auto-login saat register, kembalikan user data saja
	return response.OK(c, user)
}

// Modifikasi Login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginUserRequest

	if err := c.BodyParser(&req); err != nil {
		return apperror.New(apperror.CodeBadRequest, "Invalid JSON", err)
	}
	if err := h.validator.ValidateBody(req); err != nil {
		return apperror.Validation(err)
	}

	userAgent := c.Get(fiber.HeaderUserAgent)
	clientIP := c.IP()

	// Service mengembalikan AT dan Sesi (RT)
	user, accessToken, session, err := h.authService.Login(c.UserContext(), &req, userAgent, clientIP)
	if err != nil {
		return err
	}

	// 1. Set RT sebagai HttpOnly cookie
	h.setRefreshTokenCookie(c, session)

	// 2. Kirim AT dan data user di body JSON
	return response.OK(c, fiber.Map{
		"user":         user,
		"access_token": accessToken,
	})
}

// Handler Baru: Refresh
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	// 1. Baca RT dari cookie, bukan body
	oldRefreshToken := c.Cookies(refreshTokenCookieName)
	if oldRefreshToken == "" {
		return apperror.New(apperror.CodeUnauthorized, "missing refresh token", nil)
	}

	userAgent := c.Get(fiber.HeaderUserAgent)
	clientIP := c.IP()

	// 2. Panggil service untuk rotasi
	newAccessToken, newSession, err := h.authService.Refresh(c.UserContext(), oldRefreshToken, userAgent, clientIP)
	if err != nil {
		// Jika token invalid/expired, service akan return error.
		// Kita harus clear cookie lama yg mungkin invalid.
		h.clearRefreshTokenCookie(c)
		return err
	}

	// 3. Set RT baru sebagai HttpOnly cookie
	h.setRefreshTokenCookie(c, newSession)

	// 4. Kirim AT baru di body JSON
	return response.OK(c, &dto.TokenResponse{
		AccessToken: newAccessToken,
	})
}

// Handler Baru: Logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// 1. Baca RT dari cookie
	refreshToken := c.Cookies(refreshTokenCookieName)

	// 2. Hapus sesi dari database
	if err := h.authService.Logout(c.UserContext(), refreshToken); err != nil {
		// Gagal hapus di DB, tapi kita tetap harus hapus cookie di client
		h.clearRefreshTokenCookie(c)
		return err
	}

	// 3. Hapus cookie di client
	h.clearRefreshTokenCookie(c)

	return response.OK(c, fiber.Map{"message": "logged out"})
}

// Modifikasi GoogleCallback
func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	code, state := c.Query("code"), c.Query("state")
	if code == "" || state == "" {
		return apperror.New(apperror.CodeBadRequest, "OAUTH_INVALID_REQUEST", errors.New("missing code/state"))
	}
	codeVerifier, ok := h.devStore.Consume(state)
	if !ok {
		return apperror.New(apperror.CodeBadRequest, "OAUTH_STATE_MISMATCH", errors.New("invalid state"))
	}

	tok, err := h.oauthPkg.Config.Exchange(c.UserContext(), code,
		oauth2.SetAuthURLParam("code_verifier", codeVerifier),
	)
	if err != nil {
		return apperror.New(apperror.CodeUnavailable, "OAUTH_EXCHANGE_FAILED", err)
	}

	client := h.oauthPkg.Config.Client(c.UserContext(), tok)
	uinfoResp, err := client.Get("https://openidconnect.googleapis.com/v1/userinfo")
	if err != nil {
		return apperror.New(apperror.CodeUnavailable, "OAUTH_USERINFO_FAILED", err)
	}
	defer uinfoResp.Body.Close()
	if uinfoResp.StatusCode >= 400 {
		return apperror.New(apperror.CodeUnavailable, "OAUTH_USERINFO_BAD_STATUS", err)
	}

	var userInfo dto.GoogleLoginResponse
	if err := json.NewDecoder(uinfoResp.Body).Decode(&userInfo); err != nil {
		return apperror.New(apperror.CodeUnavailable, "OAUTH_USERINFO_PARSE_FAILED", err)
	}

	// Cari user berdasarkan Google ID
	user, _ := h.userService.FindUserByGoogleID(c.UserContext(), userInfo.Sub)

	// Jika tidak ada, cari berdasarkan Email
	if user == nil && userInfo.Email != "" {
		user, _ = h.userService.FindUserByEmail(c.UserContext(), userInfo.Email)
		if user != nil {
			// TODO: Tautkan akun (Update user dengan GoogleID) - (Opsional)
		}
	}

	// Jika tetap tidak ada, buat user baru
	if user == nil {
		user, err = h.authService.RegisterByGoogle(c.UserContext(), &userInfo)
		if err != nil {
			return apperror.New(apperror.CodeInternal, "USER_CREATE_FAILED", err)
		}
	}

	// Buat Sesi (AT & RT)
	userAgent := c.Get(fiber.HeaderUserAgent)
	clientIP := c.IP()

	accessToken, session, err := h.authService.CreateSession(c.UserContext(), user, userAgent, clientIP)
	if err != nil {
		return apperror.New(apperror.CodeInternal, "JWT_CREATE_FAILED", err)
	}

	// 1. Set RT sebagai HttpOnly cookie
	h.setRefreshTokenCookie(c, session)

	// 2. Kirim AT dan data user di body JSON
	// PENTING: Untuk callback, seringkali lebih baik me-redirect kembali ke
	// frontend dengan token di query param. Tapi untuk API, ini OK.
	// Mari kita asumsikan client API-based.
	return response.OK(c, fiber.Map{
		"access_token": accessToken,
		"user":         user,
	})
}

func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	state, verifier, challenge, err := h.devStore.NewStateAndPKCE()
	if err != nil || verifier == "" {
		return apperror.New(apperror.CodeInternal, "OAUTH_STATE_GEN_FAILED", err)
	}
	url := h.oauthPkg.Config.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)
	return c.Redirect(url)
}
