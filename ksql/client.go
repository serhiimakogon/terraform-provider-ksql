package ksql

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	context.Context
	cli       *http.Client
	url       string
	apiKey    string
	apiSecret string
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

func NewEmptyClient(ctx context.Context) *Client {
	return NewClientContext(ctx, "", "", "")
}

func NewClientContext(ctx context.Context, url, apiKey, apiSecret string) *Client {
	client := &Client{
		Context:   ctx,
		url:       url,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		cli:       &http.Client{},
	}
	return client
}

func (ksql *Client) ListStreams() ([]Stream, error) {
	payload := Payload{
		Ksql: "LIST STREAMS;",
	}

	response, err := ksql.makePostRequest(payload)
	if err != nil {
		return nil, err
	}

	return response[0].Streams, nil
}

func (ksql *Client) GetStreamByName(streamName string) (*Stream, error) {
	listStreams, err := ksql.ListStreams()
	if err != nil {
		return nil, err
	}

	for _, stream := range listStreams {
		if strings.EqualFold(stream.Name, streamName) {
			return &stream, nil
		}
	}

	return nil, fmt.Errorf("there is no stream named %s", streamName)
}

func (ksql *Client) GetStreamsByTopic(topicName string) ([]Stream, error) {
	listStreams, err := ksql.ListStreams()
	if err != nil {
		return nil, err
	}

	streams := []Stream{}
	for _, stream := range listStreams {
		if strings.EqualFold(stream.Topic, topicName) {
			streams = append(streams, stream)
		}
	}

	return streams, nil
}

func (ksql *Client) GetStreamsByTag(tag string) ([]Stream, error) {
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

	return streams, nil
}

func (ksql *Client) CreateStream(streamName string, query string) (Response, error) {
	payload := Payload{
		Ksql: fmt.Sprintf("CREATE STREAM %s %s", streamName, query),
	}

	err := validateStreamName(streamName)
	if err != nil {
		return nil, err
	}

	err = validateQuery(query)
	if err != nil {
		return nil, err
	}

	_, err = ksql.GetStreamByName(streamName)
	if err == nil {
		return nil, fmt.Errorf("there is already a stream named %s", streamName)
	}

	response, err := ksql.makePostRequest(payload)
	if err != nil {
		return nil, err
	}

	if response[0].ErrorCode != 0 {
		return nil, errors.New(response[0].Message)
	}

	return response, nil
}

func (ksql *Client) DropStream(streamName string) (Response, error) {
	payload := Payload{
		Ksql: fmt.Sprintf("DROP STREAM %s;", streamName),
	}

	err := validateStreamName(streamName)
	if err != nil {
		return nil, err
	}

	_, err = ksql.GetStreamByName(streamName)
	if err != nil {
		return nil, fmt.Errorf("there is no stream named %s", streamName)
	}

	response, err := ksql.makePostRequest(payload)
	if err != nil {
		return nil, err
	}

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

	req = req.WithContext(context.Background())

	req.SetBasicAuth(ksql.apiKey, ksql.apiSecret)

	resp, err := ksql.cli.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
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

func validateStreamName(streamName string) error {
	if strings.Contains(streamName, "-") {
		return fmt.Errorf("stream name should not contain '-' character ")
	}

	return nil
}

func validateQuery(query string) error {
	lastCharacter := query[len(query)-1:]
	if lastCharacter != ";" {
		return fmt.Errorf("query missing ';' at the end ")
	}

	return nil
}
