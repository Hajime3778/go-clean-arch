package nethttp

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Hajime3778/go-clean-arch/domain"
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
