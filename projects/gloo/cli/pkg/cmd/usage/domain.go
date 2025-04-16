package usage

import (
	"fmt"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/snapshot"
	"sort"
)

type FeatureType string
type Label string

type Category string

type API string

const (
	GEMENI_AI_BACKEND       FeatureType = "Gemini AI Backend"
	LLM_AI_BACKEND          FeatureType = "LLM AI Backend"
	AZURE_OPENAI_AI_BACKEND FeatureType = "Azure OpenAI AI Backend"
	ANTHROPIC_AI_BACKEND    FeatureType = "Anthropic AI Backend"
	MISTRAL_AI_BACKEND      FeatureType = "MistralAI AI Backend"
	OPENAI_AI_BACKEND       FeatureType = "OpenAI AI Backend"
	VERTEXAI_AI_BACKEND     FeatureType = "VertexAI AI Backend"

	CONSUL_BACKEND         FeatureType = "Consul Backend"
	GCP_BACKEND            FeatureType = "GCP Backend"
	STATIC_BACKEND         FeatureType = "Static Backend"
	AWS_EC2_BACKEND        FeatureType = "AWS EC2 Backend"
	AWS_LAMBDA_BACKEND     FeatureType = "AWS Lambda Backend"
	AZURE_FUNCTION_BACKEND FeatureType = "Azure Function Backend"

	PIPE_BACKEND           FeatureType = "Pipe Backend"
	KUBE_BACKEND           FeatureType = "Kube Backend"
	UPSTREAM_HEALTH_CHECKS FeatureType = "Upstream Health Checks"
	OUTLIER_DETECTION      FeatureType = "Outlier Detection"

	UPSTREAM_HTTP2      FeatureType = "Upstream HTTP2"
	UPSTREAM_HTTP_PROXY FeatureType = "Upstream HTTP Proxy"

	FAILOVER         FeatureType = "Failover"
	CIRCUIT_BREAKERS FeatureType = "Circuit Breakers"
	UPSTREAM_MTLS    FeatureType = "Upstream mTLS"
	UPSTREAM_TLS     FeatureType = "Upstream TLS"
	REDIRECT         FeatureType = "Redirect"
	GRAPHQL          FeatureType = "GraphQl"

	DELEGATION          FeatureType = "Delegation"
	ROUTE_TABLE_WEIGHTS FeatureType = "Route Table Weights"

	DIRECT_RESPONSE            FeatureType = "Direct Response"
	AI_SEMANTIC_CACHE          FeatureType = "AI Semantic Cache"
	AI_PROMPT_ENRICHMENT       FeatureType = "AI Prompt Enrichment"
	AI_RAG                     FeatureType = "AI Rag"
	AI_PROMPT_GUARD            FeatureType = "AI Prompt Guard"
	CONNECT_MATCHING           FeatureType = "Connect Matching"
	PREFIX_PATH_MATCHING       FeatureType = "Prefix Path Matching"
	EXACT_PATH_MATCHING        FeatureType = "Exact Path Matching"
	REGEX_PATH_MATCHING        FeatureType = "Regex Path Matching"
	METHOD_MATCHING            FeatureType = "Method Matching"
	QUERY_PARAMETER_MATCHING   FeatureType = "Query Parameter Matching"
	HEADER_MATCHING            FeatureType = "Header Matching"
	PORTAL_AUTH                FeatureType = "Portal Auth"
	HMAC_AUTH                  FeatureType = "HMAC Auth"
	BUFFER_PER_ROUTE           FeatureType = "Buffer Per Route"
	PASSTHROUGH_AUTH           FeatureType = "Passthrough Auth"
	OPA_AUTH                   FeatureType = "OPA Auth"
	PLUGIN_AUTH                FeatureType = "Plugin Auth"
	API_KEY_AUTH               FeatureType = "API Key Auth"
	BASIC_AUTH                 FeatureType = "Basic Auth"
	LDAP_AUTH                  FeatureType = "LDAP Auth"
	OAUTH2_AUTH                FeatureType = "OAuth2 Auth"
	TRANSFORMATIONS            FeatureType = "Transformations"
	CORS                       FeatureType = "CORS"
	HEADER_MANIPULATION        FeatureType = "Header Manipulation"
	ROUTE_STATS                FeatureType = "Route Stats"
	RETRIES                    FeatureType = "Retries"
	HTTPS                      FeatureType = "HTTPS"
	MTLS                       FeatureType = "MTLS"
	PROXY_PROTOCOL             FeatureType = "Proxy Protocol"
	TCP_ROUTING                FeatureType = "TCP Routing"
	LOCAL_RATE_LIMITING        FeatureType = "Local Rate Limiting"
	CONNECTION_LIMIT           FeatureType = "Connection Limit"
	HTTP_CONNECTION_MANAGEMENT FeatureType = "HTTP Connection Management"
	INBOUND_HEALTH_CHECK       FeatureType = "Inbound Health Check"
	WEB_APPLICATION_FIREWALL   FeatureType = "Web Application Firewall"
	HEADER_VALIDATION          FeatureType = "Header Validation"
	TAP_FILTER                 FeatureType = "Tap Filter"
	STATEFUL_SESSION           FeatureType = "Stateful Session"
	WASM                       FeatureType = "WASM"
	DATA_LOSS_PREVENTION       FeatureType = "Data Loss Prevention"
	DYNAMIC_FORWAD_PROXY       FeatureType = "Dynamic Forward Proxy"
	CLUSTER_HEADER             FeatureType = "Cluster Header"
	MULTI_UPSTREAM_ROUTE       FeatureType = "Multi Upstream Route"
	UPSTREAM_GROUP_ROUTE       FeatureType = "Upstream Group Route"

	GRPC_JSON_TRANSCODING FeatureType = "GRPC JSON Transcoder"
	GZIP_COMPRESSION      FeatureType = "gzip Compression"
	FAULT_INJECTION       FeatureType = "Fault Injection"
	PREFIX_REWRITE        FeatureType = "Prefix Rewrite"
	UPSTREAM_TIMEOUT      FeatureType = "Request Timeout"
	ROUTE_TRACING         FeatureType = "Route Tracing"
	REQUEST_SHADOWING     FeatureType = "Request Shadowing"
	HOST_REWRITE          FeatureType = "Host Rewrite"
	AUTO_HOST_REWRITE     FeatureType = "Auto Host Rewrite"
	APPEND_XFF_HEADER     FeatureType = "Append XFF Header"
	CSRF                  FeatureType = "CSRF"
	MAX_STREAM_DURATION   FeatureType = "Max Stream Duration"
	IDLE_TIMEOUT          FeatureType = "Idle Timeout"

	EXTERNAL_PROCESSING          FeatureType = "External Processing"
	HASH_LOAD_BALANCING          FeatureType = "Hash Load Balancing"
	LOCALITY_LOAD_BALANCING      FeatureType = "Locality Load Balancing"
	LEAST_REQUEST_LOAD_BALANCING FeatureType = "Least Request Load Balancing"
	MAGLEV_LOAD_BALANCING        FeatureType = "Maglev Load Balancing"
	ROUND_ROBIN_LOAD_BALANCING   FeatureType = "Round Robin Load Balancing"
	RANDOM_LOAD_BALANCING        FeatureType = "Random Load Balancing"
	RING_HASH_LOAD_BALANCING     FeatureType = "Ring Hash Load Balancing"

	REGEX_REWRITE                      FeatureType = "Regex Rewrite"
	RESPONSE_CACHING                   FeatureType = "Response Caching"
	BUFFER_LIMIT                       FeatureType = "Buffer Limit"
	TLS_ROUTING                        FeatureType = "TLS Routing"
	SOCKET_OPTIONS                     FeatureType = "Socket Options"
	MOST_SPECIAL_HEADER_MUTATIONS_WINS FeatureType = "Most Specific Header Mutations Wins"
	HYBRID_GATEWAY                     FeatureType = "Hybrid Gateway"
	MAX_DIRECT_RESPONSE_BODY_SIZE      FeatureType = "Max Direct Response Body Size"
	ACCESS_LOGGING                     FeatureType = "Access Logging"
	EARLY_ACCESS_LOGGING               FeatureType = "Early Access Logging"
	CONNECTION_BALANCING               FeatureType = "Connection Balancing"
	RATE_LIMITING                      FeatureType = "Rate Limiting"
	JWT                                FeatureType = "JWT"
	RBAC                               FeatureType = "RBAC"

	NameLabel Label = "name"
	HostLabel Label = "host"
	KindLabel Label = "kind"

	listenerCatagory Category = "Listener"
	routingCatagory  Category = "Routing"
	aiCategory       Category = "AI"
	upstreamCategory Category = "Upstream"

	GlooEdgeAPI API = "Gloo Edge API"
	GatewayAPI  API = "Gateway API"
	KGatewayAPI API = "kGateway API"
)

type UsageStats struct {
	stats map[API][]*UsageStat
	cache *snapshot.Instance
}

func (u *UsageStats) Print() {
	//featureCount := map[FeatureType]int{}

	fmt.Printf("Gloo Edge APIs: \n")
	fmt.Printf("\tGloo Gateways: %d\n", len(u.cache.GlooGateways()))
	fmt.Printf("\tVirtualService: %d\n", len(u.cache.VirtualServices()))
	fmt.Printf("\tRouteTables: %d\n", len(u.cache.RouteTables()))
	fmt.Printf("\tUpstreams: %d\n", len(u.cache.Upstreams()))
	fmt.Printf("\tVirtualHostOptions: %d\n", len(u.cache.VirtualHostOptions()))
	fmt.Printf("\tListenerOptions: %d\n", len(u.cache.ListenerOptions()))
	fmt.Printf("\tHTTPListenerOptions: %d\n", len(u.cache.HTTPListenerOptions()))
	fmt.Printf("\tAuthConfigs: %d\n", len(u.cache.AuthConfigs()))
	fmt.Printf("\tSettings: %d\n", len(u.cache.Settings()))

	fmt.Printf("\nGateway API APIs: \n")
	fmt.Printf("\tGateways: %d\n", len(u.cache.Gateways()))
	fmt.Printf("\tListenerSets: %d\n", len(u.cache.ListenerSets()))
	fmt.Printf("\tHTTPRoutes: %d\n", len(u.cache.HTTPRoutes()))
	fmt.Printf("\tDirectResponses: %d\n", len(u.cache.DirectResponses()))
	fmt.Printf("\tGatewayParameters: %d\n", len(u.cache.GatewayParameters()))

	fmt.Printf("\nTotal Features Used Per API\n")
	fmt.Printf("\tGloo Edge API: %d\n", len(u.stats[GlooEdgeAPI]))
	fmt.Printf("\tGateway API: %d\n", len(u.stats[GatewayAPI]))
	fmt.Printf("\tkGateway API: %d\n\n", len(u.stats[KGatewayAPI]))

	for api, stats := range u.stats {
		fmt.Printf("API: %s", api)

		// organize by category
		categories := map[Category]map[FeatureType]int{}

		// group all the stats by their codes
		for _, stat := range stats {
			//featureCount[stat.Type]++
			if categories[stat.Metadata.Category] == nil {
				categories[stat.Metadata.Category] = map[FeatureType]int{}
			}
			categories[stat.Metadata.Category][stat.Type]++

		}
		for category, features := range categories {
			fmt.Printf("\n\tCategory: %s", category)
			//sort the features
			//add all the features to a list for naming
			featureList := []FeatureType{}
			for key := range features {
				featureList = append(featureList, key)
			}

			sort.Slice(featureList, func(i, j int) bool {
				return featureList[i] < featureList[j]
			})
			for _, feature := range featureList {
				fmt.Printf("\n\t\t%s: %d", feature, features[feature])
			}
		}
	}
}

type UsageMetadata struct {
	Name      string
	Namespace string
	Kind      string
	Category  Category
	API       API
}

type UsageStat struct {
	Type     FeatureType
	Metadata UsageMetadata
}

func (u *UsageStats) AddUsageStat(stat *UsageStat) {
	if stat.Metadata.API == "" {
		stat.Metadata.API = GlooEdgeAPI
	}
	if u.stats == nil {
		u.stats = make(map[API][]*UsageStat)
	}
	if u.stats[stat.Metadata.API] == nil {
		u.stats[stat.Metadata.API] = []*UsageStat{}
	}
	u.stats[stat.Metadata.API] = append(u.stats[stat.Metadata.API], stat)
}


type OutputStats struct {
	Features *GlooFeatures
	ProxyStats    *ProxyStats
}

type GlooFeatures struct {
	
}

type ProxyStats struct {
	TotalRequests int
	Uptime int
	UpstreamStats []UpstreamStats
}

type UpstreamStats struct {
	UpstreamName string
	TotalRequests int
}
