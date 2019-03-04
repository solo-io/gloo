package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/solo-io/solo-projects/pkg/license/db"
	"github.com/solo-io/solo-projects/pkg/license/keys"
	"github.com/solo-io/solo-projects/pkg/license/model"
	"github.com/solo-io/solo-projects/pkg/license/notify"
)

type LicenseGenServer struct {
	KeyGenerator keys.KeyGenerator
	KeyDb        db.KeyDb
	Notifier     notify.Notifier
	KeyDuration  time.Duration
}

const MaxBody = 10 * 1 << 10

func (l *LicenseGenServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := extractRequest(r)
	userInfo, err := l.getUserInfo(&request)
	if err != nil {
		l.handleError(w, err)
		return
	}

	licenseKey, err := l.genKey(ctx)
	if err != nil {
		l.handleError(w, err)
		return
	}

	err = l.saveRequest(ctx, &request, userInfo, licenseKey)
	if err != nil {
		l.handleError(w, err)
		return
	}

	err = l.emailKeyToUser(ctx, userInfo, licenseKey)
	if err != nil {
		l.saveEmail(ctx, userInfo, licenseKey, false)
		l.handleError(w, err)
		return
	} else {
		l.saveEmail(ctx, userInfo, licenseKey, true)
	}

	l.handleSuccess(w)

}
func (l *LicenseGenServer) saveEmail(ctx context.Context, ua model.UserInfo, k string, emailSuccess bool) {
	// TODO: make this work
}

func (l *LicenseGenServer) handleSuccess(w http.ResponseWriter) {
	// do nothing for now; potentially redirect somewhere here..
	w.WriteHeader(http.StatusNoContent)
}

func (l *LicenseGenServer) saveRequest(ctx context.Context, r *model.Request, ua model.UserInfo, k string) error {
	// save the whole request to a db
	return l.KeyDb.Save(ctx, r, ua, k)
}

func (l *LicenseGenServer) genKey(ctx context.Context) (string, error) {
	return l.KeyGenerator.GenerateKey(ctx, time.Now().Add(l.KeyDuration))
}

func (l *LicenseGenServer) emailKeyToUser(ctx context.Context, ua model.UserInfo, key string) error {
	return l.Notifier.Notify(ctx, ua, key)
}

func (l *LicenseGenServer) handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
}

func (l *LicenseGenServer) getUserInfo(r *model.Request) (model.UserInfo, error) {
	contentType := r.Headers.Get("content-type")
	var ui model.UserInfo
	var err error
	switch {
	case strings.Contains(contentType, "application/json"):
		err = json.Unmarshal(r.Body, &ui)
	case strings.Contains(contentType, "application/x-www-form-urlencoded"):
		var vs url.Values
		vs, err = url.ParseQuery(string(r.Body))
		ui.Email = vs.Get("email")
	default:
		err = fmt.Errorf("unknown content-type")
	}

	return ui, err
}

func extractRequest(r *http.Request) model.Request {
	defer r.Body.Close()
	var buf bytes.Buffer
	io.Copy(&buf, io.LimitReader(r.Body, MaxBody))
	return model.Request{
		Body:    buf.Bytes(),
		Headers: r.Header,
	}

}
