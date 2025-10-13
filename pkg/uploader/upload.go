package uploader

import (
	"context"
	"encoding/json"
	"fmt"
	"portal/pkg/constants"
	"portal/pkg/consul"
	"net/http"
	"os"

	"github.com/hashicorp/consul/api"
)

type Avatar struct {
	Url string `json:"url"`
}

type ImageKey struct {
	Key string `json:"key"`
}

type ImageService interface {
	GetImageKey(ctx context.Context, key string) (*Avatar, error)
	DeleteImageKey(ctx context.Context, key string) error
}

type imageService struct {
	client *callAPI
}

type callAPI struct {
	client       consul.ServiceDiscovery
	clientServer *api.CatalogService
}

var (
	imageServiceStr = "go-main-service"
)

func NewImageService(client *api.Client) ImageService {
	mainServiceAPI := NewServiceAPI(client, imageServiceStr)
	return &imageService{
		client: mainServiceAPI,
	}
}

func NewServiceAPI(client *api.Client, serviceName string) *callAPI {
	sd, err := consul.NewServiceDiscovery(client, serviceName)
	if err != nil {
		fmt.Printf("Error creating service discovery: %v\n", err)
		return nil
	}

	service, err := sd.DiscoverService()
	if err != nil {
		fmt.Printf("Error discovering service: %v\n", err)
		return nil
	}

	if os.Getenv("LOCAL_TEST") == "true" {
		fmt.Println("Running in LOCAL_TEST mode — overriding service address to localhost")
		service.ServiceAddress = "localhost"
	}

	return &callAPI{
		client:       sd,
		clientServer: service,
	}
}

func (s *imageService) GetImageKey(ctx context.Context, key string) (*Avatar, error) {
	
    token, ok := ctx.Value(constants.TokenKey).(string)
    if !ok || token == "" {
        return nil, fmt.Errorf("token not found in context")
    }

    image, err := s.client.getImageKey(key, token)
    if err != nil {
        return nil, err
    }

    if image != nil {
        if sc, ok := image["status_code"].(float64); ok && int(sc) == 500 {
            return nil, nil
        }
    }

    if image == nil {
        return nil, nil
    }
    innerData, ok := image["data"].(string)
    if !ok || innerData == "" {
        return nil, nil
    }

    return &Avatar{Url: innerData}, nil
}


func (s *imageService) DeleteImageKey(ctx context.Context, key string) error {

	token, ok := ctx.Value(constants.TokenKey).(string)
	if !ok {
		return fmt.Errorf("token not found in context")
	}

	err := s.client.deleleImage(key, token)

	if err != nil {
		return err
	}

	return nil

}

func (c *callAPI) deleleImage(key string, token string) (error) {

	endpoint := "/v1/images/delete"

	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	body := map[string]string{
		"key":  key,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshalling body: %v", err)
	}

	_, err = c.client.CallAPI(c.clientServer, endpoint, http.MethodPost, jsonBody, header)
	if err != nil {
		fmt.Printf("Error calling API: %v\n", err)
		return err
	}

	return nil
}

func (c *callAPI) getImageKey(key string, token string) (map[string]interface{}, error) {

    endpoint := "/v1/images"
    header := map[string]string{
        "Content-Type":  "application/json",
        "Authorization": "Bearer " + token,
    }
    body := map[string]string{"key": key, "mode": "public"}

    jsonBody, err := json.Marshal(body)
    if err != nil {
        return nil, fmt.Errorf("error marshalling body: %v", err)
    }

    res, err := c.client.CallAPI(c.clientServer, endpoint, http.MethodPost, jsonBody, header)
    if err != nil {
        return nil, err
    }

    var imageData interface{}
    if err := json.Unmarshal([]byte(res), &imageData); err != nil {
        return nil, fmt.Errorf("error unmarshalling response: %v", err)
    }

    myMap, ok := imageData.(map[string]interface{})
    if !ok {
        return nil, fmt.Errorf("unexpected response format")
    }
    return myMap, nil

}

