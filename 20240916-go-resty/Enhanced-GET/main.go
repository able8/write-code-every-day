package main

// Import resty into your code and refer it as `resty`.
import (
	"errors"
	"net/url"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

func main() {
	// Create a Resty Client
	client := resty.New()

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"page_no": "1",
			"limit":   "20",
			"sort":    "name",
			"order":   "asc",
			"random":  strconv.FormatInt(time.Now().Unix(), 10),
		}).
		SetHeader("Accept", "application/json").
		SetAuthToken("BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F").
		Get("/search_result")

	// Sample of using Request.SetQueryString method
	resp, err = client.R().
		SetQueryString("productId=232&template=fresh-sample&cat=resty&source=google&kw=buy a lot more").
		SetHeader("Accept", "application/json").
		SetAuthToken("BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F").
		Get("/show_product")

	// If necessary, you can force response content type to tell Resty to parse a JSON response into your struct
	resp, err = client.R().
		SetResult(result).
		ForceContentType("application/json").
		Get("v2/alpine/manifests/latest")

		// Create a Resty Client
	client := resty.New()

	// Setting output directory path, If directory not exists then resty creates one!
	// This is optional one, if you're planning using absolute path in
	// `Request.SetOutput` and can used together.
	client.SetOutputDirectory("/Users/jeeva/Downloads")

	// HTTP response gets saved into file, similar to curl -o flag
	_, err := client.R().
		SetOutput("plugin/ReplyWithHeader-v5.1-beta.zip").
		Get("http://bit.ly/1LouEKr")

	// OR using absolute path
	// Note: output directory path is not used for absolute path
	_, err := client.R().
		SetOutput("/MyDownloads/plugin/ReplyWithHeader-v5.1-beta.zip").
		Get("http://bit.ly/1LouEKr")

		// Create a Resty Client
	client := resty.New()

	// just mentioning about POST as an example with simple flow
	// User Login
	resp, err := client.R().
		SetFormData(map[string]string{
			"username": "jeeva",
			"password": "mypass",
		}).
		Post("http://myapp.com/login")

	// Followed by profile update
	resp, err := client.R().
		SetFormData(map[string]string{
			"first_name": "Jeevanandam",
			"last_name":  "M",
			"zip_code":   "00001",
			"city":       "new city update",
		}).
		Post("http://myapp.com/profile")

	// Multi value form data
	criteria := url.Values{
		"search_criteria": []string{"book", "glass", "pencil"},
	}
	resp, err := client.R().
		SetFormDataFromValues(criteria).
		Post("http://myapp.com/search")

		// Create a Resty Client
	client := resty.New()

	// Registering Request Middleware
	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		// Now you have access to Client and current Request object
		// manipulate it as per your need

		return nil // if its success otherwise return error
	})

	// Registering Response Middleware
	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		// Now you have access to Client and current Response object
		// manipulate it as per your need

		return nil // if its success otherwise return error
	})

	// Create a Resty Client
	client := resty.New()

	// Retries are configured per client
	client.
		// Set retry count to non zero to enable retries
		SetRetryCount(3).
		// You can override initial retry wait time.
		// Default is 100 milliseconds.
		SetRetryWaitTime(5 * time.Second).
		// MaxWaitTime can be overridden as well.
		// Default is 2 seconds.
		SetRetryMaxWaitTime(20 * time.Second).
		// SetRetryAfter sets callback to calculate wait time between retries.
		// Default (nil) implies exponential backoff with jitter
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			return 0, errors.New("quota exceeded")
		})

}
