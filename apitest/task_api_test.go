package apitest_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Hajime3778/go-clean-arch/domain"
	"github.com/stretchr/testify/assert"
)

const taskURL = "http://localhost:8080/tasks"

func TestGetByID(t *testing.T) {
	t.Run("正常系 存在するIDで1件取得", func(t *testing.T) {
		req, _ := http.NewRequest("GET", taskURL+"/2", nil)

		expectedTask := domain.Task{
			ID:      2,
			UserID:  1,
			Title:   "買い出しに行く",
			Content: "スーパーで、卵と鶏肉と三葉を買う",
			DueDate: time.Now(),
		}

		client := new(http.Client)
		response, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}
		defer response.Body.Close()

		var resTask domain.Task
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&resTask)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectedTask.ID, resTask.ID)
		assert.Equal(t, expectedTask.Title, resTask.Title)
		assert.Equal(t, expectedTask.Content, resTask.Content)
	})
}
