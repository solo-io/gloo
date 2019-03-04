package main

import (
	"net/http"
	"time"

	debugdb "github.com/solo-io/solo-projects/pkg/license/db/debug"
	"github.com/solo-io/solo-projects/pkg/license/keys/jwt"
	debugnotify "github.com/solo-io/solo-projects/pkg/license/notify/debug"
	"github.com/solo-io/solo-projects/pkg/license/server"
)

func main() {

	secret := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0}
	kg := &jwt.KeyGenHMAC{
		Secret: secret,
	}

	lgs := &server.LicenseGenServer{
		KeyGenerator: kg,
		KeyDb:        new(debugdb.DebugKeyDb),
		Notifier:     new(debugnotify.DebugNotifier),
		KeyDuration:  30 * 24 * time.Hour,
	}
	s := server.NewDebugServer(lgs)

	http.ListenAndServe(":8080", s)

}
