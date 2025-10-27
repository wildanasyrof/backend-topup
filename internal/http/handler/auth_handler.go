package handler

import (
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/service"
	apperror "github.com/wildanasyrof/backend-topup/pkg/apperr"
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
		return apperror.New(apperror.CodeBadRequest, "Invalid JSON", err)
	}

	if err := h.validator.ValidateBody(req); err != nil {
		return apperror.Validation(err)
	}

	user, err := h.authService.Register(c.UserContext(), &req)
	if err != nil {
		return err

	}

	return response.OK(c, user)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginUserRequest

	if err := c.BodyParser(&req); err != nil {
		return apperror.New(apperror.CodeBadRequest, "Invalid JSON", err)
	}

	if err := h.validator.ValidateBody(req); err != nil {
		return apperror.Validation(err)
	}

	user, token, err := h.authService.Login(c.UserContext(), &req)
	if err != nil {
		return err
	}

	return response.OK(c, fiber.Map{
		"user":  user,
		"token": token,
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

func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	code, state := c.Query("code"), c.Query("state")
	if code == "" || state == "" {
		return apperror.New(apperror.CodeBadRequest, "OAUTH_INVALID_REQUEST", errors.New("missing code/state"))

	}
	codeVerifier, ok := h.devStore.Consume(state)
	if !ok {
		return apperror.New(apperror.CodeBadRequest, "OAUTH_STATE_MISMATCH", errors.New("invalid state"))
	}

	// Exchange WITH secret (from config) AND PKCE verifier
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

	user, _ := h.userService.FindUserByGoogleID(c.UserContext(), userInfo.Sub)

	if user == nil && userInfo.Email != "" {
		user, _ = h.userService.FindUserByEmail(c.UserContext(), userInfo.Email)
	}

	if user == nil {
		user, err = h.authService.RegisterByGoogle(c.UserContext(), &userInfo)
		if err != nil {
			return apperror.New(apperror.CodeInternal, "USER_CREATE_FAILED", err)
		}
	}

	token, err := h.authService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return apperror.New(apperror.CodeInternal, "JWT_CREATE_FAILED", err)
	}
	// TODO: find/create user, issue JWT, return {token, user}
	return response.OK(c, fiber.Map{
		"access_token": token,
		"user":         user,
	})
}
