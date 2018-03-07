package reqtests

import (
	"testing"

	"git.containerum.net/ch/kube-client/pkg/model"
)

func TestTokenMethods(test *testing.T) {
	client := newMockClient(test)
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
