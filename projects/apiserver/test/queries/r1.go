package queries

// key: 2019-03-28 16:43:09.960448 -0400 EDT m=+4.836362095
const query_8765159107952418687_0 = `{"operationName":null,"variables":{},"query":"{\n  allNamespaces {\n    virtualServices {\n      metadata {\n        guid\n        resourceVersion\n      }\n    }\n    upstreams {\n      metadata {\n        name\n      }\n    }\n  }\n}\n"}`

// key: 2019-03-28 16:43:12.454645 -0400 EDT m=+7.330560589
const query_8765159107952418687_1 = `{"operationName":null,"variables":{},"query":"{\n  namespace(name: \"gloo-system\") {\n    virtualServices {\n      rateLimitConfig {\n        anonymousLimits {\n          unit\n          requestsPerUnit\n        }\n      }\n    }\n  }\n}\n"}`

var queries_8765159107952418687 = []string{
	query_8765159107952418687_0,
	query_8765159107952418687_1,
}
