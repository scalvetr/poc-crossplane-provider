package poc

import (
	"github.com/scalvetr/poc-crossplane-provider/service-client/pkg/poc"
	"net/http"
)

type POCService struct {
	POCClient poc.Client
}

var topicNotFoundError *poc.TopicNotFoundError

func NewService(tenant string, baseUrl string, _ []byte) (POCService, error) {
	return POCService{
		POCClient: poc.Client{
			Tenant:  tenant,
			BaseURL: baseUrl,
			Client:  http.Client{},
		},
	}, nil
}

func (c *POCService) GetTopics() ([]poc.Topic, error) {
	return c.POCClient.GetTopics()
}
func (c *POCService) GetTopic(topicName string) (*poc.Topic, error) {
	topic, err := c.POCClient.GetTopic(topicName)
	if _, isTopicNotFoundError := err.(*poc.TopicNotFoundError); isTopicNotFoundError {
		return nil, nil
	}
	return topic, err
}

func (c *POCService) CreateTopic(t poc.Topic) error {
	return c.POCClient.CreateTopic(t)
}

func (c *POCService) UpdateTopic(t poc.Topic) error {
	err := c.POCClient.UpdateTopic(t)
	if _, isTopicNotFoundError := err.(*poc.TopicNotFoundError); isTopicNotFoundError {
		return nil
	}
	return err
}

func (c *POCService) DeleteTopic(topicName string) error {
	err := c.POCClient.DeleteTopic(topicName)
	if _, isTopicNotFoundError := err.(*poc.TopicNotFoundError); isTopicNotFoundError {
		return nil
	}
	return err
}
