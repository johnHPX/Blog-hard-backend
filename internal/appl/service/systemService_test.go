package service

import "testing"

func TestSendEmail(t *testing.T) {
	t.Run("teste positivo", func(t *testing.T) {
		systemService := NewSystemService()
		err := systemService.SendEmail("<h1>isso é um teste</h1>", "kiramy763@gmail.com", "isso é um teste")
		if err != nil {
			t.Error(err)
		}
	})
}
