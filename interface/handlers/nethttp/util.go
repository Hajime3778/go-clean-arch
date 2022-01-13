package nethttp

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Hajime3778/go-clean-arch/domain"
	"github.com/form3tech-oss/jwt-go"
	"github.com/form3tech-oss/jwt-go/request"
)

// WriteJSONResponse JSON形式でレスポンスを出力します
func WriteJSONResponse(w http.ResponseWriter, status int, body interface{}) {
	json, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Println(err.Error())
		return
	}
	w.WriteHeader(status)
	w.Write(json)
	log.Println(string(json))
}

// GetStatusCode エラー内容からHttpStatusCodeを返却します
func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrRecordNotFound:
		return http.StatusNotFound
	case domain.ErrBadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// VerifyAccessToken アクセストークン署名を検証し、トークンとUserIDを返却します。
func VerifyAccessToken(r *http.Request) (string, int64, error) {
	token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		b := []byte(os.Getenv("SECRET_KEY"))
		return b, nil
	})

	if err != nil {
		return "", 0, err
	}

	claims := token.Claims.(jwt.MapClaims)
	strUserID := claims["user_id"].(string)
	userID, _ := strconv.ParseInt(strUserID, 10, 64)

	return token.Raw, userID, nil
}
