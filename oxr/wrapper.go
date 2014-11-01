package oxr

import (
//    "math/big"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "fmt"
)

var ApiUrl = "openexchangerates.org/api"
var ErrNotFound = "not_found"
var ErrNotAvailable = "not_available"
var ErrMissingAppId = "missing_app_id"
var ErrInvalidAppId = "invalid_app_id"
var ErrNotAllowed = "not_allowed"
var ErrAccessRestricted = "access_restricted"
var ErrInvalidBase = "invalid_base"
var ProtoHttp = "http"
var ProtoHttps = "https"

type ApiError struct {
    IsError     bool    `json:"error"`
    Status      int     `json:"status"`
    Message     string  `json:"message"`
    Description string  `json:"description"`
}

func (e ApiError) Error() string {
    return fmt.Sprintf("%v: %v", e.Message, e.Description)
}

type Rates struct {
    Disclaimer  string                  `json:"disclaimer"`
    License     string                  `json:"license"`
    Timestamp   int                     `json:"timestamp"`
    Base        string                  `json:"base"`
    Rates       map[string]json.Number  `json:"rates"`
}

type ApiClient struct {
    appId       string
    proto       string
    url         string
}

func New(appId string) *ApiClient {
    return NewWithOptions(appId, ProtoHttps, ApiUrl)
}

func NewWithOptions(appId string, proto string, url string) *ApiClient {
    return &ApiClient{appId: appId, proto: proto, url: url}
}

func (c *ApiClient) Currencies() (map[string]string, error) {
    data, err := c.apiCall("currencies", nil)
    if err != nil { return nil, err }

    var curr map[string]string
    err = json.Unmarshal(data, &curr)
    if err != nil { return nil, err }

    return curr, nil
}

func (c *ApiClient) Latest() (*Rates, error) {
    return c.rates("latest")
}

func (c *ApiClient) Historical(date string) (*Rates, error) {
    return c.rates("historical/" + date)
}

func (c *ApiClient) rates(endpoint string) (*Rates, error) {
    data, err := c.apiCall(endpoint, nil)
    if err != nil { return nil, err }

    var r Rates
    err = json.Unmarshal(data, &r)
    if err != nil { return nil, err }

    return &r, nil
}

func (c *ApiClient) apiCall(endpoint string, args map[string]string) ([]byte, error) {
    // Build URL
    url := fmt.Sprintf("%v://%v/%v.json?app_id=%v", c.proto, c.url, endpoint, c.appId)
    for k := range args {
        url = fmt.Sprintf("%v&%v=%v", url, k, args[k])
    }

    // Make request
    resp, err := http.Get(url)
    if err != nil { return nil, err }

    // Retrieve raw JSON response
    var body []byte
    body, err = ioutil.ReadAll(resp.Body)
    if err != nil { return nil, err }
    defer resp.Body.Close()

    // Process API-level error conditions
    if resp.StatusCode != 200 {
        var e ApiError
        err = json.Unmarshal(body, &e)
        if err != nil { return nil, err }
        return nil, e
    }

    return body, nil
}
