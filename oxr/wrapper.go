package oxr

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "strings"
    "fmt"
)

var (
    ErrNotFound = "not_found"
    ErrNotAvailable = "not_available"
    ErrMissingAppId = "missing_app_id"
    ErrInvalidAppId = "invalid_app_id"
    ErrNotAllowed = "not_allowed"
    ErrAccessRestricted = "access_restricted"
    ErrInvalidBase = "invalid_base"
    ErrInvalidCurrency = "invalid_currency"
    ErrInvalidDateRange = "invalid_date_range"
    ProtoHttp = "http"
    ProtoHttps = "https"
    ApiUrl = "openexchangerates.org/api"
)

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

type RatesSeries struct {
    Disclaimer  string                              `json:"disclaimer"`
    License     string                              `json:"license"`
    StartDate   string                              `json:"start_date"`
    EndDate     string                              `json:"end_date"`
    Base        string                              `json:"base"`
    Rates       map[string]map[string]json.Number   `json:"rates"`
}

type Conversion struct {
    Disclaimer  string                  `json:"disclaimer"`
    License     string                  `json:"license"`
    Request     ConversionRequest       `json:"request"`
    Meta        ConversionMeta          `json:"meta"`
    Response    json.Number             `json:"response"`
}

type ConversionRequest struct {
    Query       string                  `json:"query"`
    Amount      json.Number             `json:"amount"`
    From        string                  `json:"from"`
    To          string                  `json:"to"`
}

type ConversionMeta struct {
    Timestamp   int                     `json:"timestamp"`
    Rate        json.Number             `json:"rate"`
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
    data, err := c.apiCall("currencies.json", nil)
    if err != nil { return nil, err }

    var curr map[string]string
    err = json.Unmarshal(data, &curr)
    if err != nil { return nil, err }

    return curr, nil
}

func (c *ApiClient) Latest() (*Rates, error) {
    return c.LatestWithOptions("USD", nil)
}

func (c *ApiClient) LatestWithOptions(base string, symbols []string) (*Rates, error) {
    return c.rates("latest.json", base, symbols)
}

func (c *ApiClient) Historical(date string) (*Rates, error) {
    return c.HistoricalWithOptions(date, "USD", nil)
}

func (c *ApiClient) HistoricalWithOptions(date string, base string, symbols []string) (*Rates, error) {
    return c.rates("historical/" + date + ".json", base, symbols)
}

func (c *ApiClient) TimeSeries(start string, end string) (*RatesSeries, error) {
    return c.TimeSeriesWithOptions(start, end, "USD", nil)
}

func (c *ApiClient) TimeSeriesWithOptions(start string, end string, base string, symbols []string) (*RatesSeries, error) {
    args := make(map[string]string)

    args["start"] = start
    args["end"] = end

    if base != "USD" && base != "" {
        args["base"] = base
    }

    if symbols != nil {
        args["symbols"] = strings.Join(symbols, ",")
    }

    data, err := c.apiCall("time-series.json", args)
    if err != nil { return nil, err }

    var r RatesSeries
    err = json.Unmarshal(data, &r)
    if err != nil { return nil, err }

    return &r, nil
}

func (c *ApiClient) Convert(amount string, from string, to string) (*Conversion, error) {
    ep := fmt.Sprintf("convert/%v/%v/%v", amount, from, to)
    data, err := c.apiCall(ep, nil)
    if err != nil { return nil, err }

    var conv Conversion
    err = json.Unmarshal(data, &conv)
    if err != nil { return nil, err }

    return &conv, nil
}

func (c *ApiClient) rates(endpoint string, base string, symbols []string) (*Rates, error) {
    args := make(map[string]string)

    if base != "USD" && base != "" {
        args["base"] = base
    }

    if symbols != nil {
        args["symbols"] = strings.Join(symbols, ",")
    }

    data, err := c.apiCall(endpoint, args)
    if err != nil { return nil, err }

    var r Rates
    err = json.Unmarshal(data, &r)
    if err != nil { return nil, err }

    return &r, nil
}

func (c *ApiClient) apiCall(endpoint string, args map[string]string) ([]byte, error) {
    // Build URL
    url := fmt.Sprintf("%v://%v/%v?app_id=%v", c.proto, c.url, endpoint, c.appId)
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
