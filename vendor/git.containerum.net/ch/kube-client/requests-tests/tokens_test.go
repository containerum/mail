package reqtests

import (
	"testing"

	"git.containerum.net/ch/kube-client/pkg/model"
)

func TestTokenMethods(test *testing.T) {
	client := newClient(test)
	client.UserManagerURL = "http://192.168.88.200:8111"
	client.AuthURL = "http://192.168.88.200:1111"
	client.SetHeaders(map[string]string{
		"X-User-Agent":  "kube-client",
		"X-User-Client": "315d3143bab041b3656e4666355adb15",
		"X-Client-IP":   "192.168.0.1",
	})
	username := "helpik94@yandex.ru"
	password := "12345678"
	recaptcha := "03AHhf_52156hcrzZpAgJse24k1JVDN4nGjujmnlYW7KTjV-JuxmNE13SUfJNfxEC1Rj4"
	login := model.Login{
		Username:  username,
		Password:  password,
		Recaptcha: &recaptcha,
	}
	tokens, err := client.Login(login)
	if err != nil {
		test.Fatalf("error while login: %v", err)
	}

	_, err = client.CheckToken(tokens.AccessToken)
	if err != nil {
		test.Fatalf("error while checking token: %v", err)
	}
	_, err = client.ExtendToken(tokens.RefreshToken)
	if err != nil {
		test.Fatalf("error while refreshing token: %v", err)
	}
}
