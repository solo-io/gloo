package server

import (
	"net/http"
)

type DebugServer struct {
}

func (*DebugServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte(`
	<html>
	<body>
	<form action="/post" method="POST">

	<label for="email">email</label><input type="email" id="email" />
	<input type="submit" />
	</form>
	</body>
	</html>
	
	`))

}

func NewDebugServer(l *LicenseGenServer) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/post", l)
	mux.Handle("/", new(DebugServer))
	return mux
}
