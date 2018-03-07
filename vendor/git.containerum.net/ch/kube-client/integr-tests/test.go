package test

import (
	"log"
	"testing"

	"git.containerum.net/ch/kube-client/pkg/model"
)

func TestDeployment(test *testing.T) {
	client := newClient(test)
	client.SetFingerprint("514c67239bcd3f2b7837eb9a3edc30bc")
	tokens, err := client.Login(model.Login{
		Username: "helpik94@yandex.ru",
		Password: "12345678",
	})
	if err != nil {
		test.Fatalf("error while login: %v", err)
	}
	client.SetToken(tokens.AccessToken)

	ns, err := client.GetNamespaceList(nil)
	if err != nil {
		log.Fatalf("error while getting namespace: %v", err)
	}
	test.Logf("%v", ns)
}
