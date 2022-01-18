package apitest_test

import "testing"

func TestSignUp(t *testing.T) {
	t.Run("正常系 ユーザーの新規登録", func(t *testing.T) {})
	t.Run("準正常系 リクエストパラメータが足りていない場合、400エラーとなること", func(t *testing.T) {})
	t.Run("準正常系 リクエスト形式が間違っている場合、400エラーとなること", func(t *testing.T) {})
}

func TestSignIn(t *testing.T) {
	t.Run("正常系 ユーザーの新規登録後、サインインできること", func(t *testing.T) {})
	t.Run("準正常系 存在しないEmailの場合、401エラーとなること", func(t *testing.T) {})
	t.Run("準正常系 パスワードが間違っている場合、401エラーとなること", func(t *testing.T) {})
	t.Run("準正常系 リクエストパラメータが足りていない場合、400エラーとなること", func(t *testing.T) {})
	t.Run("準正常系 リクエスト形式が間違っている場合、400エラーとなること", func(t *testing.T) {})
}
