package transforms

import "net/http"

type CookieData struct {
	// The value can be either a string or Gomega matcher
	Value    interface{}
	HttpOnly bool
}

// CookieMapper returns a Gomega Transform that maps an array of HTTP cookies the to the form <cookieName>: <CookieMapData>
func CookieDataMapper() func(c []*http.Cookie) map[string]*CookieData {
	return func(c []*http.Cookie) map[string]*CookieData {
		cookieMap := make(map[string]*CookieData)
		for _, cookie := range c {
			cookieMap[cookie.Name] = &CookieData{
				Value:    cookie.Value,
				HttpOnly: cookie.HttpOnly,
			}
		}

		return cookieMap
	}
}
