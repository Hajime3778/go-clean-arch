package auth

import (
	"encoding/json"
	"net/http"

	"github.com/Hajime3778/go-clean-arch/domain"
	httpUtil "github.com/Hajime3778/go-clean-arch/interface/handlers/nethttp"
	usecase "github.com/Hajime3778/go-clean-arch/usecase/auth"
)

const SignUpPath string = "/auth/sign_up"
const SignInPath string = "/auth/sign_in"

type authHandler struct {
	authUsecase usecase.AuthUsecase
}

// NewAuthHandler 認証機能のHandlerオブジェクトを作成します
func NewAuthHandler(u usecase.AuthUsecase) *authHandler {
	return &authHandler{u}
}

// SignUpHandler
func (t *authHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var request SignUpRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&request)
	if err != nil {
		httpUtil.WriteJSONResponse(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	var ok bool
	if ok, err = request.IsSignUpRequestValid(); !ok {
		httpUtil.WriteJSONResponse(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	user := domain.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}

	token, err := t.authUsecase.SignUp(ctx, user)
	if err != nil {
		httpUtil.WriteJSONResponse(w, httpUtil.GetStatusCode(err), domain.ErrorResponse{Message: err.Error()})
		return
	}

	httpUtil.WriteJSONResponse(w, http.StatusOK, SignUpResponse{Token: token})
}

// SignInHandler
func (t *authHandler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var request SignInRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&request)
	if err != nil {
		httpUtil.WriteJSONResponse(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	var ok bool
	if ok, err = request.IsSignInRequestValid(); !ok {
		httpUtil.WriteJSONResponse(w, http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	token, err := t.authUsecase.SignIn(ctx, request.Email, request.Password)
	if err != nil {
		httpUtil.WriteJSONResponse(w, httpUtil.GetStatusCode(err), domain.ErrorResponse{Message: err.Error()})
		return
	}

	httpUtil.WriteJSONResponse(w, http.StatusOK, SignInResponse{Token: token})
}
