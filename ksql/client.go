package ksql

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	context.Context
	client   *http.Client
	url      string
	username string
	password string
}

type CommandStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Response []struct {
	Streams       []Stream      `json:"streams"`
	CommandStatus CommandStatus `json:"commandStatus"`
	ErrorCode     int           `json:"error_code"`
	Message       string        `json:"message"`
}

type Stream struct {
	Name  string `json:"name"`
	Topic string `json:"topic"`
}

type Payload struct {
	Ksql string `json:"ksql"`
}

func NewClient(url string, username string, password string) *Client {
	return NewClientContext(context.Background(), url, username, password)
}

func NewClientContext(ctx context.Context, url string, username string, password string) *Client {
	client := &Client{
		Context:  ctx,
		url:      url,
		username: username,
		password: password,
		client:   &http.Client{},
	}
	return client
}

func (ksql *Client) ListStreams() ([]Stream, error) {
	payload := Payload{
		Ksql: "LIST STREAMS;",
	}

	response, _ := ksql.makePostRequest(payload)

	return response[0].Streams, nil
}

func (ksql *Client) GetStreamByName(streamName string) (*Stream, error) {
	listStreams, err := ksql.ListStreams()
	if err != nil {
		return nil, err
	}

	for _, stream := range listStreams {
		if stream.Name == streamName {
			return &stream, nil
		}
	}

	return nil, fmt.Errorf("there is no stream named %s", streamName)
}

func (ksql *Client) GetStreamsByTopic(topicName string) (*[]Stream, error) {
	listStreams, err := ksql.ListStreams()
	if err != nil {
		return nil, err
	}

	streams := []Stream{}
	for _, stream := range listStreams {
		if stream.Topic == topicName {
			streams = append(streams, stream)
		}
	}

	return &streams, nil
}

func (ksql *Client) GetStreamsByTag(tag string) (*[]Stream, error) {
	listStreams, err := ksql.ListStreams()
	if err != nil {
		return nil, err
	}

	streams := []Stream{}
	for _, stream := range listStreams {
		if strings.Contains(stream.Name, tag) {
			streams = append(streams, stream)
		}
	}

	return &streams, nil
}

func (ksql *Client) CreateStream(streamName string, query string) (Response, error) {
	payload := Payload{
		Ksql: fmt.Sprintf("CREATE STREAM %s %s;", streamName, query),
	}

	response, _ := ksql.makePostRequest(payload)

	if response[0].ErrorCode != 0 {
		return nil, errors.New(response[0].Message)
	}

	return response, nil
}

func (ksql *Client) DropStream(streamName string) (Response, error) {
	payload := Payload{
		Ksql: fmt.Sprintf("DROP STREAM %s;", streamName),
	}

	response, _ := ksql.makePostRequest(payload)

	if response[0].ErrorCode != 0 {
		return nil, errors.New(response[0].Message)
	}

	return response, nil
}

func (ksql *Client) makePostRequest(payload Payload) (Response, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/ksql", ksql.url), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ksql)

	req.SetBasicAuth(ksql.username, ksql.password)

	resp, err := ksql.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := Response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
