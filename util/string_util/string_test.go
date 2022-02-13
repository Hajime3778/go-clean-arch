package string_util_test

import (
	"testing"

	"github.com/Hajime3778/go-clean-arch/util/string_util"
	"github.com/stretchr/testify/assert"
)

func TestGenerateRundomString(t *testing.T) {
	t.Run("正常系 文字列の桁数が正しいこと", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			str := string_util.GenerateRundomString(i)
			assert.Equal(t, i, len(str))
		}
	})
}
