package main

import (
	"fmt"
	"log"

	"github.com/solo-io/solo-projects/projects/apiserver/test/harness"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	apiServer := harness.ApiServer{
		Origin: "http://localhost:8082",
		// Token:  "any_string",
	}
	if err := callQ(apiServer, from_query_2207477075972191016_2); err != nil {
		return err
	}
	if err := callQ(apiServer, from_query_9013262992668048334_27); err != nil {
		return err
	}
	return nil

}

func callQ(apiServer harness.ApiServer, q string) error {
	result, err := apiServer.CallQuery(q)
	if err != nil {
		return err
	}
	if result.Errors != nil {
		return result.Errors
	}
	fmt.Println(result)
	return nil
}

const from_query_2207477075972191016_2 = `{"operationName":"CreateVirtualService","variables":{"vs":{"metadata":{"name":"m","namespace":"gloo-system","resourceVersion":"one"},"displayName":"the display name","extAuthConfig":{"oAuth":{"clientId":"one","clientSecretRef":{"name":"aname","namespace":"gloo-system"},"issuerUrl":"web.com","appUrl":"hey.com","callbackPath":"/the/callback"}}}},"query":"mutation CreateVirtualService($vs: InputVirtualService!) {\n  virtualServices {\n    create(virtualService: $vs) {\n      metadata {\n        name\n        guid\n      }\n      extAuthConfig {\n        authType {\n          ... on OAuthConfig {\n            clientId\n            clientSecret\n            issuerUrl\n            appUrl\n            callbackPath\n          }\n        }\n      }\n    }\n  }\n}\n"}`

const from_query_9013262992668048334_27 = `{"operationName":"AddRoute","variables":{"vsid":"*v1.VirtualService gloo-system m","rv":"23268","ind":1,"route":{"matcher":{"pathMatch":"/home2","pathMatchType":"EXACT"},"redirectAction":{"prefixRewrite":"/new/name"},"destination":{"singleDestination":{"upstream":{"name":"extauth","namespace":"gloo-system"}}}}},"query":"mutation AddRoute($vsid: ID!, $rv: String!, $ind: Int!, $route: InputRoute!) {\n  virtualServices {\n    addRoute(virtualServiceId: $vsid, resourceVersion: $rv, index: $ind, route: $route) {\n      routes {\n        matcher {\n          pathMatch\n        }\n        redirectAction {\n          prefixRewrite\n        }\n      }\n    }\n  }\n}\n"}`
