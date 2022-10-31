package hclient

import (
	"encoding/json"
	"net/http"
)

// Header returns the response headers
func (r *Response) Header() http.Header {
	return r.headers
}

// Body returns the response body
func (r *Response) Body() []byte {
	return r.body
}

// BodyAsString returns the response body as a string
func (r *Response) BodyAsString() string {
	return string(r.body)
}

// BodyAsJSON unmarshalls the current response body to the specified data structure
func (r *Response) BodyAsJSON(data interface{}) error {
	return json.Unmarshal(r.body, data)
}

// Status returns the response status
func (r *Response) Status() string {
	return r.status
}

// StatusCode returns the response status code
func (r *Response) StatusCode() int {
	return r.statusCode
}

func (r *Response) GetCookieByName(name string) *http.Cookie {
	for _, cookie := range r.cookies {
		if cookie.Name == name {
			return cookie
		}
	}

	return nil
}
