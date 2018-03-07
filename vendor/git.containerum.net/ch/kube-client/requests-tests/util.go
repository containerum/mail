package reqtests

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"

	kubeClient "git.containerum.net/ch/kube-client/pkg/client"
	"git.containerum.net/ch/kube-client/pkg/model"
	"git.containerum.net/ch/kube-client/pkg/rest/re"
	"git.containerum.net/ch/kube-client/pkg/rest/remock"
)

const (
	testAPIurl = "http://192.168.88.200"
)

func newMockClient(test *testing.T) *kubeClient.Client {
	client, err := kubeClient.NewClient(kubeClient.Config{
		APIurl:  "http://192.168.88.200",
		RestAPI: remock.NewMock(),
		User: kubeClient.User{
			Role: "user",
		},
	})
	if err != nil {
		test.Fatalf("error while client initialisation: %v", err)
	}
	return client
}
func newClient(test *testing.T) *kubeClient.Client {
	client, err := kubeClient.NewClient(
		kubeClient.Config{
			RestAPI: re.NewResty(),
			APIurl:  "http://192.168.88.200",
			User: kubeClient.User{
				Role: "user",
			},
		})
	if err != nil {
		test.Fatalf("error while creating client: %v", err)
	}
	return client
}

func newFakeNamespaces(test *testing.T) []model.Namespace {
	return []model.Namespace{
		{
			TariffID: "4563e8c1-fb41-416a-9798-e949a2616260",
		},
	}
}

func newFakeDeployment(test *testing.T) model.Deployment {
	deployment := model.Deployment{
		Name:     "gateway",
		Replicas: 4,
		Labels:   map[string]string{},
		Containers: []model.Container{
			{
				Name: "proxy", Image: "nginx",
				Limits: model.Resource{CPU: "1", Memory: "256"},
				Ports: []model.ContainerPort{
					{Name: "Gate", Port: 1080, Protocol: model.TCP},
				},
				Env: []model.Env{
					{Name: "TEAPOT", Value: "TRUE"},
				},
			},
		},
	}
	return deployment
}

func newFakeUpdateImage(test *testing.T) model.UpdateImage {
	var updateImage model.UpdateImage
	loadTestJSONdata(test, "test_data/update_image.json", &updateImage)
	return updateImage
}

func newFakeKubeAPInamespace(test *testing.T) model.Namespace {
	var namespace model.Namespace
	loadTestJSONdata(test, "test_data/kube_api_namespace.json", &namespace)
	return namespace
}

func newFakeVolume(test *testing.T) []model.Volume {
	var volume []model.Volume
	loadTestJSONdata(test, "test_data/fake_volumes.json", &volume)
	return volume
}

func loadTestJSONdata(test *testing.T, file string, data interface{}) {
	jsonData, err := ioutil.ReadFile(file)
	if err != nil {
		test.Fatalf("error wgile reading from %q: %v", file, err)
	}
	err = json.Unmarshal(jsonData, data)
	if err != nil {
		test.Fatalf("error while unmarshalling data: %v", err)
	}
}

func newRandomName(size int64) string {
	buf := &bytes.Buffer{}
	encoder := base64.NewEncoder(base64.RawURLEncoding, buf)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	io.CopyN(encoder, rnd, (3*size)/4)
	return buf.String()
}

func newUpdateDeployment(test *testing.T) model.Deployment {
	var deployment model.Deployment
	loadTestJSONdata(test, "test_data/update_deployment.json", &deployment)
	return deployment
}
