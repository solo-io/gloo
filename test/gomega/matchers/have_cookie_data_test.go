package matchers_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-projects/test/gomega/matchers"
	"github.com/solo-io/solo-projects/test/gomega/transforms"
)

var _ = Describe("CookieDataMapper", func() {
	var testData = []*http.Cookie{
		{
			Name:     "cookie1",
			Value:    "value1",
			HttpOnly: true,
		},
		{
			Name:     "cookie2",
			Value:    "value2",
			HttpOnly: false,
		},
	}

	It("matches", func() {
		Expect(testData).To(WithTransform(transforms.CookieDataMapper(), And(
			HaveKeyWithValue("cookie1", matchers.MatchCookieData(&transforms.CookieData{
				Value:    "value1",
				HttpOnly: true,
			})),
			HaveKeyWithValue("cookie2", matchers.MatchCookieData(&transforms.CookieData{
				Value:    "value2",
				HttpOnly: false,
			})),
		)))
	})

	It("matches substring in value", func() {
		Expect(testData).To(WithTransform(transforms.CookieDataMapper(), And(
			HaveKeyWithValue("cookie1", matchers.MatchCookieData(&transforms.CookieData{
				Value:    ContainSubstring("value"),
				HttpOnly: true,
			})),
			HaveKeyWithValue("cookie2", matchers.MatchCookieData(&transforms.CookieData{
				Value:    ContainSubstring("2"),
				HttpOnly: false,
			})),
		)))
	})

	It("partially matches value", func() {
		Expect(testData).To(WithTransform(transforms.CookieDataMapper(), And(
			HaveKeyWithValue("cookie1", matchers.MatchCookieData(&transforms.CookieData{
				Value:    "value1",
				HttpOnly: true,
			})),
			HaveKeyWithValue("cookie2", matchers.MatchCookieData(&transforms.CookieData{
				HttpOnly: false,
			})),
		)))
	})

	It("fails to match when value is wrong", func() {
		Expect(testData).To(WithTransform(transforms.CookieDataMapper(), Not(
			HaveKeyWithValue("cookie1", matchers.MatchCookieData(&transforms.CookieData{
				Value:    "value",
				HttpOnly: true,
			})),
		)))
	})

	It("fails to match when http only is wrong", func() {
		Expect(testData).To(WithTransform(transforms.CookieDataMapper(), Not(
			HaveKeyWithValue("cookie1", matchers.MatchCookieData(&transforms.CookieData{
				Value:    "value1",
				HttpOnly: false,
			})),
		)))
	})
})
