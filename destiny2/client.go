package destiny2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"

	"golang.org/x/oauth2"
)

const (

	// BaseURL is the base URL for the Bungie API.
	BaseURL = "https://www.bungie.net/Platform"
)

// RequestOption can be passed to functions that accept it to modify the http request
// made in the function before if is dispathced
type RequestOption = func(*http.Request) *http.Request

// OptionBody adds a body to a request
func OptionBody(r io.ReadCloser) RequestOption {
	return func(req *http.Request) *http.Request {
		req.Body = r
		return req
	}
}

// OptionContext adds a context to a request
func OptionContext(ctx context.Context) RequestOption {
	return func(req *http.Request) *http.Request {
		return req.WithContext(ctx)
	}
}

// OptionQuery adds a quesy param to a request
func OptionQuery(key string, value interface{}) RequestOption {
	return func(req *http.Request) *http.Request {
		q := req.URL.Query()
		q.Add(key, fmt.Sprint(value))
		req.URL.RawQuery = q.Encode()

		return req
	}
}

// OptionOAuthToken sets the authorization header on the request to the provided token
func OptionOAuthToken(t *oauth2.Token) RequestOption {
	return func(req *http.Request) *http.Request {
		t.SetAuthHeader(req)
		return req
	}
}

// APIResponse represents part of a response from returned from every Bungie API endpoint
type APIResponse struct {
	ErrorCode          int               `json:"ErrorCode"`
	ThrottleSeconds    int               `json:"ThrottleSeconds"`
	ErrorStatus        string            `json:"ErrorStatus"`
	Message            string            `json:"Message"`
	MessageData        map[string]string `json:"MessageData"`
	DetailedErrorTrace string            `json:"DetailedErrorTrace"`
}

// PagedQuery ...
// https://bungie-net.github.io/multi/schema_Queries-PagedQuery.html#schema_Queries-PagedQuery
type PagedQuery struct {
	ItemsPerPage             int    `json:"itemsPerPage"`
	CurrentPage              int    `json:"currentPage"`
	RequestContinuationToken string `json:"requestContinuationToken"`
}

// Client is used to communicate with the Desinty 2 API
type Client struct {
	httpClient   *http.Client
	apiKey       string
	oauth2Config *oauth2.Config

	GroupV2Service  *GroupV2Service
	Destiny2Service *Destiny2Service
	UserService     *UserService
}

// NewClient creates and returns a new client
func NewClient(apiKey string) *Client {
	c := &Client{
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
	}

	// Adding services to client
	c.GroupV2Service = &GroupV2Service{c}
	c.Destiny2Service = &Destiny2Service{c}
	c.UserService = &UserService{c}

	return c
}

// SetOAuthCredentials sets the OAuth2 credentials on the client so the client can generate authroize URLs
// and exchange codes for authorization tokens
func (c *Client) SetOAuthCredentials(clientID, clientSecret string) *Client {
	c.oauth2Config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.bungie.net/en/oauth/authorize",
			TokenURL: "https://www.bungie.net/Platform/App/OAuth/token/",
		},
	}

	return c
}

// GetOAuthConfig gets the OAuth2 config set on the client.
// If no config has been set yet using Client.SetOAuthCredentials an zero config will be returned
func (c *Client) GetOAuthConfig() oauth2.Config {
	if c.oauth2Config == nil {
		return oauth2.Config{}
	}

	return *c.oauth2Config
}

// WithOAuth2Token returns a copy of Client but with the http client being set to one that is authorized
// to make calls to authenticated endpoints
func (c Client) WithOAuth2Token(t *oauth2.Token) *Client {
	ctx := context.Background()
	c.SetHTTPClient(c.oauth2Config.Client(ctx, t))

	return &c
}

// SetHTTPClient sets the client used when communicating with the Destiny 2 API. Function returns
// self for ease of chaining
func (c *Client) SetHTTPClient(httpClient *http.Client) *Client {
	c.httpClient = httpClient
	return c
}

// GetAuthURL generates a auth URL to send to a user so they can authorize the app to access their account information.
// State is not nessesary but is strongly advised
func (c *Client) GetAuthURL(state string) string {
	return c.oauth2Config.AuthCodeURL(state)
}

// Exchange exchanges the code provided with Destiny 2's token servers for an access token that can be used to create a client
// for making authroized requests to endpoints that require authentication
func (c *Client) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return c.oauth2Config.Exchange(ctx, code)
}

func (c *Client) do(method, endpoint string, dst interface{}, opts ...RequestOption) error {
	u, _ := url.Parse(BaseURL)
	u.Path = path.Join(u.Path, endpoint) + "/"

	// Creating request
	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Add("X-Api-Key", c.apiKey)

	// Applying options to request
	for _, opt := range opts {
		req = opt(req)
	}

	// Sending request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Getting mime type to determine if we can json parse the body or if the request errored
	// and we should check for error on body
	contentType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return err
	}

	// Request errored and bungie did not return JSON in the body so we have to rely on the status
	// code to determine what went wrong
	if contentType != "application/json" {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return ErrUnautorized
		case http.StatusNotFound:
			return ErrNotFound
		default:
			return ErrUnknown
		}
	}

	// Consuming body into intermidiate state to read the error codes to determine if the request
	// has failed and in what way
	var respStruct struct {
		Response json.RawMessage `json:"Response"`
		APIResponse
	}
	if err = json.NewDecoder(resp.Body).Decode(&respStruct); err != nil {
		return err
	}

	// Request errored
	if respStruct.ErrorCode != 1 {
		switch respStruct.ErrorCode {
		case 21:
			return ErrNotFound
		case 99:
			return ErrWebAuthRequired
		default:
			return ErrUnknown
		}
	}

	return json.Unmarshal(respStruct.Response, dst)
}
