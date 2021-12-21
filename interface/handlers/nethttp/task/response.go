package task

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Hajime3778/go-clean-arch/domain"
)

type ResponseError struct {
	Message string `json:"message"`
}

// writeJSONResponse JSON形式でレスポンスを出力します
func writeJSONResponse(w http.ResponseWriter, status int, body interface{}) {
	json, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(status)
		w.Write([]byte(err.Error()))
		log.Println(err.Error())
		return
	}
	w.WriteHeader(status)
	w.Write(json)
	log.Println(string(json))
}

// getStatusCode エラー内容からHttpStatusCodeを返却します
func getStatusCode(err error) int {
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
