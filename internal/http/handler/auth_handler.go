package handler

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/service"
	"github.com/wildanasyrof/backend-topup/pkg/oauth"
	"github.com/wildanasyrof/backend-topup/pkg/response"
	"github.com/wildanasyrof/backend-topup/pkg/validator"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	authService service.AuthService
	userService service.UserService
	oauthPkg    oauth.GoogleOauthPkg
	devStore    oauth.DevStore
	validator   validator.Validator
}

func NewAuthHandler(authService service.AuthService, userService service.UserService, oauthPkg oauth.GoogleOauthPkg, devStore oauth.DevStore, validator validator.Validator) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
		oauthPkg:    oauthPkg,
		devStore:    devStore,
		validator:   validator,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterUserRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	if err := h.validator.ValidateBody(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "validation error", err)
	}

	user, err := h.authService.Register(c.UserContext(), &req)
	if err != nil {
		errMsg := err.Error()

		if strings.Contains(errMsg, "uni_users_email") {
			return response.Error(c, fiber.StatusConflict, "Registration failed", fiber.Map{
				"message": "Email already exists",
				"field":   "email",
			})
		}

		return response.Error(c, fiber.StatusInternalServerError, "Registration failed", err.Error())

	}

	return response.Success(c, "user registered succesfully", user)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginUserRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if err := h.validator.ValidateBody(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "validation error", err)
	}

	user, token, err := h.authService.Login(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "login failed", err.Error())
	}

	return response.Success(c, "login success", fiber.Map{
		"user":  user,
		"token": token,
	})
}

func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	state, verifier, challenge, err := h.devStore.NewStateAndPKCE()
	if err != nil || verifier == "" {
		return response.Error(c, 500, "OAUTH_STATE_GEN_FAILED", err)
	}
	url := h.oauthPkg.Config.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)
	return c.Redirect(url)
}

func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	code, state := c.Query("code"), c.Query("state")
	if code == "" || state == "" {
		return response.Error(c, 400, "OAUTH_INVALID_REQUEST", "missing code/state")
	}
	codeVerifier, ok := h.devStore.Consume(state)
	if !ok {
		return response.Error(c, 400, "OAUTH_STATE_MISMATCH", "invalid state")
	}

	// Exchange WITH secret (from config) AND PKCE verifier
	tok, err := h.oauthPkg.Config.Exchange(c.UserContext(), code,
		oauth2.SetAuthURLParam("code_verifier", codeVerifier),
	)
	if err != nil {
		return response.Error(c, 502, "OAUTH_EXCHANGE_FAILED", err.Error())
	}

	client := h.oauthPkg.Config.Client(c.UserContext(), tok)
	uinfoResp, err := client.Get("https://openidconnect.googleapis.com/v1/userinfo")
	if err != nil {
		return response.Error(c, 502, "OAUTH_USERINFO_FAILED", err.Error())
	}
	defer uinfoResp.Body.Close()
	if uinfoResp.StatusCode >= 400 {
		return response.Error(c, 502, "OAUTH_USERINFO_BAD_STATUS", uinfoResp.Status)
	}

	var userInfo dto.GoogleLoginResponse
	if err := json.NewDecoder(uinfoResp.Body).Decode(&userInfo); err != nil {
		return response.Error(c, 502, "OAUTH_USERINFO_PARSE_FAILED", err.Error())
	}

	user, _ := h.userService.FindUserByGoogleID(c.UserContext(), userInfo.Sub)

	if user == nil && userInfo.Email != "" {
		user, _ = h.userService.FindUserByEmail(c.UserContext(), userInfo.Email)
	}

	if user == nil {
		user, err = h.authService.RegisterByGoogle(c.UserContext(), &userInfo)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "USER_CREATE_FAILED", err.Error())
		}
	}

	token, err := h.authService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "JWT_CREATE_FAILED", err.Error())
	}
	// TODO: find/create user, issue JWT, return {token, user}
	return response.Success(c, "Login success", fiber.Map{
		"access_token": token,
		"user":         user,
	})
}
