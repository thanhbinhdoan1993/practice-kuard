package memqclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/thanhbinhdoan1993/practice-kuard/pkg/memq"
)

type Client struct {
	BaseServerURL string
}

func errorFromResponse(resp *http.Response) error {
	if resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP Error: %v", resp.Status)
	}
	return nil
}

func (c *Client) queueURL(queue string, s ...string) string {
	s = append([]string{"queues", queue}, s...)
	tail := path.Join(s...)
	return fmt.Sprintf("%s/%s", c.BaseServerURL, tail)
}

func (c *Client) CreateQueue(queue string) error {
	req, err := http.NewRequest(http.MethodPut, c.queueURL(queue), nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return errorFromResponse(resp)
}

func (c *Client) DeleteQueue(queue string) error {
	req, err := http.NewRequest(http.MethodDelete, c.queueURL(queue), nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return errorFromResponse(resp)
}

func (c *Client) DrainQueue(queue string) error {
	req, err := http.NewRequest(http.MethodPost, c.queueURL(queue, "drain"), nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return errorFromResponse(resp)
}

func (c *Client) Enqueue(queue, data string) (*memq.Message, error) {
	req, err := http.NewRequest(
		http.MethodPost, c.queueURL(queue, "enqueue"),
		bytes.NewBufferString(data))
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	err = errorFromResponse(resp)
	if err != nil {
		return nil, err
	}

	m := &memq.Message{}
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Dequeue takes an item off of queue from the server. If a nil message is
// returned with no error then the queue is empty.
func (c *Client) Dequeue(queue string) (*memq.Message, error) {
	req, err := http.NewRequest(http.MethodPost, c.queueURL(queue, "dequeue"), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	err = errorFromResponse(resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	m := &memq.Message{}
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return nil, err
	}
	return m, err
}

func (c *Client) Stats() (*memq.Stats, error) {
	req, err := http.NewRequest(http.MethodGet, c.BaseServerURL+"/stats", nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	err = errorFromResponse(resp)
	if err != nil {
		return nil, err
	}

	s := &memq.Stats{}
	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		return nil, err
	}
	return s, nil
}
