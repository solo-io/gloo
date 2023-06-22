package assertions

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	. "github.com/onsi/gomega"
)

/*
  This file was created by copying the contents of e2e/ratelimit_test.go.
  There are improvements to these assertions that we can make, to leverage our
  existing http response matchers, but to reduce the footprint of changes, we will
  keep the existing assertions as-is for now.
*/

func ConsistentlyNotRateLimited(hostname string, port uint32) {
	// There are occasional failures in CI around the 'consistently' portion of this.
	// As a result, we wrap this logic in an 'eventually' assertion, in hopes that will resolve the flakes
	// There may be a better solution to this, but we have not explored this further
	Eventually(func(g Gomega) {
		testStatus(g, hostname, port, nil, http.StatusOK, 2, false)
		testStatus(g, hostname, port, nil, http.StatusOK, 2, true)
	}, "30s", "1s").Should(Succeed())
}

func EventuallyRateLimited(hostname string, port uint32) {
	testStatus(Default, hostname, port, nil, http.StatusTooManyRequests, 2, false)
}

func EventuallyRateLimitedWithHeaders(hostname string, port uint32, headers http.Header) {
	testStatus(Default, hostname, port, headers, http.StatusTooManyRequests, 2, false)
}

func EventuallyRateLimitedWithExpectedHeaders(hostname string, port uint32, expectedHeaders http.Header) {
	testHeaders(Default, hostname, port, nil, http.StatusTooManyRequests, 2, false, expectedHeaders)
}

func testHeaders(g Gomega, hostname string, port uint32, requestHeaders http.Header, expectedStatus int,
	offset int, consistently bool, expectedHeaders http.Header) {
	parts := strings.SplitN(hostname, "/", 2)
	hostname = parts[0]
	path := "1"
	if len(parts) > 1 {
		path = parts[1]
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/%s", "localhost", port, path), nil)
	g.Expect(err).NotTo(HaveOccurred())
	if len(requestHeaders) > 0 {
		req.Header = requestHeaders
	}

	// remove password part if exists
	parts = strings.SplitN(hostname, "@", 2)
	if len(parts) > 1 {
		hostname = parts[1]
		auth := strings.Split(parts[0], ":")
		req.SetBasicAuth(auth[0], auth[1])
	}

	req.Host = hostname

	if consistently {
		g.ConsistentlyWithOffset(offset, func() (bool, error) {
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return false, err
			}
			defer resp.Body.Close()
			_, _ = io.ReadAll(resp.Body)
			isok := resp.StatusCode == expectedStatus
			for hk := range expectedHeaders {
				if resp.Header.Get(hk) != expectedHeaders.Get(hk) {
					fmt.Println("got ", hk, resp.Header.Get(hk))
					fmt.Println("wanted ", hk, expectedHeaders.Get(hk))
					isok = false
				}
			}
			return isok, nil
		}, "3s", ".1s").Should(BeTrue())
	} else {
		g.EventuallyWithOffset(offset, func() (bool, error) {
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return false, err
			}
			defer resp.Body.Close()
			_, _ = io.ReadAll(resp.Body)
			isok := resp.StatusCode == expectedStatus
			for hk := range expectedHeaders {
				if resp.Header.Get(hk) != expectedHeaders.Get(hk) {
					fmt.Println("got ", hk, resp.Header.Get(hk))
					fmt.Println("wanted ", hk, expectedHeaders.Get(hk))
					isok = false
				}
			}
			return isok, nil
		}, "5s", ".1s").Should(BeTrue())
	}
}

func testStatus(g Gomega, hostname string, port uint32, headers http.Header, expectedStatus int,
	offset int, consistently bool) {
	parts := strings.SplitN(hostname, "/", 2)
	hostname = parts[0]
	path := "1"
	if len(parts) > 1 {
		path = parts[1]
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%d/"+path, "localhost", port), nil)
	g.Expect(err).NotTo(HaveOccurred())
	if len(headers) > 0 {
		req.Header = headers
	}

	// remove password part if exists
	parts = strings.SplitN(hostname, "@", 2)
	if len(parts) > 1 {
		hostname = parts[1]
		auth := strings.Split(parts[0], ":")
		req.SetBasicAuth(auth[0], auth[1])
	}

	req.Host = hostname

	if consistently {
		g.ConsistentlyWithOffset(offset, func() (int, error) {
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return 0, err
			}
			defer resp.Body.Close()
			_, _ = io.ReadAll(resp.Body)
			return resp.StatusCode, nil
		}, "3s", ".1s").Should(Equal(expectedStatus))
	} else {
		g.EventuallyWithOffset(offset, func() (int, error) {
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return 0, err
			}
			defer resp.Body.Close()
			_, _ = io.ReadAll(resp.Body)
			return resp.StatusCode, nil
		}, "5s", ".1s").Should(Equal(expectedStatus))
	}
}
