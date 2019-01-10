package config

var (
	APIVersion = "20181102"

	// the development server started by react-scripts defaults to ports 3000, 3001, etc. depending on what's available
	CorsAllowedOrigins = []string{
		"http://localhost:3000",
		"http://localhost:3001",
		"http://localhost:3002",
		"http://localhost:8000",
		"localhost/:1",
		"http://localhost:8082",
		"http://localhost",
	}
	CorsAllowedHeaders = []string{"Authorization", "Content-Type"}
)
