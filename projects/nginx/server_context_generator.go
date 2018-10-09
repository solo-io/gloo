package nginx

const serverContextTemplateText = `
server {
}
`

func GenerateServerContext() ([]byte, error) {
	serverContext := []byte(serverContextTemplateText)
	return serverContext, nil
}
