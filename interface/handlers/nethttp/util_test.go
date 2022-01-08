package nethttp_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Hajime3778/go-clean-arch/domain"
	"github.com/Hajime3778/go-clean-arch/interface/handlers/nethttp"
	"github.com/stretchr/testify/assert"
)

type TestResponse struct {
	Number json.Number `json:"number"`
}

func TestWriteJSONResponse(t *testing.T) {
	t.Run("正常系 JSON文字列となり、リクエストされたStatusCodeが返却されること", func(t *testing.T) {
		w := httptest.NewRecorder()
		nethttp.WriteJSONResponse(w, http.StatusOK, TestResponse{Number: "123456"})
		res := w.Result()
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "{\"number\":123456}", string(data))
	})

	t.Run("異常系 JSON変換に失敗した場合、500エラーとなること", func(t *testing.T) {
		w := httptest.NewRecorder()
		nethttp.WriteJSONResponse(w, http.StatusOK, TestResponse{Number: "fooo"})
		res := w.Result()
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.NotNil(t, data)
	})
}

func TestGetStatusCode(t *testing.T) {
	t.Run("正常系 エラーがnilの場合、200が返却されること", func(t *testing.T) {
		status := nethttp.GetStatusCode(nil)
		assert.Equal(t, http.StatusOK, status)
	})

	t.Run("正常系 既定のエラーでない場合、500が返却されること", func(t *testing.T) {
		status := nethttp.GetStatusCode(errors.New("test error"))
		assert.Equal(t, http.StatusInternalServerError, status)
	})

	t.Run("正常系 ErrBadRequestの場合、400が返却されること", func(t *testing.T) {
		status := nethttp.GetStatusCode(domain.ErrBadRequest)
		assert.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("正常系 ErrRecordNotFoundの場合、404が返却されること", func(t *testing.T) {
		status := nethttp.GetStatusCode(domain.ErrRecordNotFound)
		assert.Equal(t, http.StatusNotFound, status)
	})

	t.Run("正常系 ErrInternalServerErrorの場合、500が返却されること", func(t *testing.T) {
		status := nethttp.GetStatusCode(domain.ErrInternalServerError)
		assert.Equal(t, http.StatusInternalServerError, status)
	})
}
