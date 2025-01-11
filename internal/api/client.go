// internal/api/client.go
package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"

	"merlion/internal/auth"
)

const (
	baseURL = "https://api.merlion.dev"
)

type Client struct {
	httpClient  *http.Client
	credentials *auth.Credentials
	token       string         // For Bearer auth
	cookies     []*http.Cookie // For Cookie auth
}

func NewClient(credentials *auth.Credentials) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("creating cookie jar: %w", err)
	}
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			Jar:     jar, // Use the created jar
		},
		credentials: credentials,
	}, nil
}

func (c *Client) setAuthHeaders(req *http.Request) {
	// Try Bearer token first
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
		return
	}

	// Fall back to Basic auth if no token
	if c.credentials != nil {
		auth := base64.StdEncoding.EncodeToString(
			[]byte(fmt.Sprintf("%s:%s", c.credentials.Email, c.credentials.Password)),
		)
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))
	}

	// Cookies are automatically handled by http.Client's cookie jar
}

func (c *Client) doRequest(method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	url := fmt.Sprintf("%s/%s", baseURL, path)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}

	// Add authentication headers
	c.setAuthHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %s (status: %d)", string(respBody), resp.StatusCode)
	}

	return respBody, nil
}

// Note operations
func (c *Client) ListNotes() ([]Note, error) {
	respBody, err := c.doRequest(http.MethodGet, "notes", nil)
	if err != nil {
		return nil, err
	}

	var notes []Note
	if err := json.Unmarshal(respBody, &notes); err != nil {
		return nil, fmt.Errorf("unmarshaling response: %w", err)
	}

	return notes, nil
}

func (c *Client) GetNote(noteID string) (*Note, error) {
	respBody, err := c.doRequest(http.MethodGet, fmt.Sprintf("notes/%s", noteID), nil)
	if err != nil {
		return nil, err
	}

	var note Note
	if err := json.Unmarshal(respBody, &note); err != nil {
		return nil, fmt.Errorf("unmarshaling response: %w", err)
	}

	return &note, nil
}

func (c *Client) CreateNote(req CreateNoteRequest) (*Note, error) {
	respBody, err := c.doRequest(http.MethodPost, "notes", req)
	if err != nil {
		return nil, err
	}

	var note Note
	if err := json.Unmarshal(respBody, &note); err != nil {
		return nil, fmt.Errorf("unmarshaling response: %w", err)
	}

	return &note, nil
}

func (c *Client) UpdateNote(noteID string, req CreateNoteRequest) (*Note, error) {
	respBody, err := c.doRequest(http.MethodPut, fmt.Sprintf("notes/%s", noteID), req)
	if err != nil {
		return nil, err
	}

	var note Note
	if err := json.Unmarshal(respBody, &note); err != nil {
		return nil, fmt.Errorf("unmarshaling response: %w", err)
	}

	return &note, nil
}

func (c *Client) DeleteNote(noteID string) error {
	_, err := c.doRequest(http.MethodDelete, fmt.Sprintf("notes/%s", noteID), nil)
	return err
}

// Login handles authentication and stores the session
func (c *Client) Login() error {
	if c.credentials == nil {
		return fmt.Errorf("no credentials provided")
	}

	// Attempt login
	_, err := c.doRequest(http.MethodPost, "users/login", map[string]string{
		"email":    c.credentials.Email,
		"password": c.credentials.Password,
	})

	return err
}
