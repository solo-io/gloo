package nginx

const httpContextTemplateText = `
http {
}
`

func GenerateHttpContext() ([]byte, error) {
	httpContext := []byte(httpContextTemplateText)
	return httpContext, nil
}
