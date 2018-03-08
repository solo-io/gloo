package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"

	"github.com/solo-io/gloo-testing/helpers/local"
	"github.com/solo-io/gloo/pkg/coreplugins/service"
)

func main() {
	os.Setenv("HELPER_TMP", os.Getenv("PWD"))

	envoy, err := localhelpers.NewEnvoyInstance()
	if err != nil {
		panic(err)
	}
	defer envoy.Clean()

	gloo, err := localhelpers.NewGlooInstance()
	if err != nil {
		panic(err)
	}
	defer gloo.Clean()

	err = envoy.Run()

	if err != nil {
		panic(err)
	}
	err = gloo.Run()
	if err != nil {
		panic(err)
	}
	envoyPort := gloo.EnvoyPort()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	backendport, response := runServer(ctx)

	serviceSpec := service.UpstreamSpec{
		Hosts: []service.Host{{
			Addr: "localhost",
			Port: backendport,
		}},
	}

	u := &v1.Upstream{
		Name: "local",
		Type: "service",
		Spec: service.EncodeUpstreamSpec(serviceSpec),
	}
	err = gloo.AddUpstream(u)

	if err != nil {
		panic(err)
	}

	v := &v1.VirtualHost{
		Name: "default",
		Routes: []*v1.Route{{
			Matcher: &v1.Route_RequestMatcher{
				RequestMatcher: &v1.RequestMatcher{
					Path: &v1.RequestMatcher_PathPrefix{PathPrefix: "/"},
				},
			},
			SingleDestination: &v1.Destination{
				DestinationType: &v1.Destination_Upstream{
					Upstream: &v1.UpstreamDestination{
						Name: u.Name,
					},
				},
			},
		}},
	}

	err = gloo.AddVhost(v)
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	body := "solo.io test"
	timeout := time.After(10 * time.Minute)
	for {
		select {
		case gotbody := <-response:
			if body == gotbody {
				fmt.Println("yes " + body)
			} else {
				fmt.Println("no  " + body)
			}
			return
		case <-time.After(time.Second):
			// call the server again is it might not have initialized
			var buf bytes.Buffer
			buf.Write([]byte(body))
			_, err := http.Post(fmt.Sprintf("http://%s:%d", "localhost", envoyPort), "application/octet-stream", &buf)
			if err != nil {
				//	fmt.Println("post err " + err.Error())
			}
		case <-timeout:
			panic("timeout")
		case <-c:
			fmt.Println("existing")
			return
		}
	}

}

func runServer(ctx context.Context) (uint32, <-chan string) {
	bodychan := make(chan string)
	handlerfunc := func(rw http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			body, _ := ioutil.ReadAll(r.Body)
			r.Body.Close()
			bodychan <- string(body)
		}
	}
	port := uint32(1334)
	handler := http.HandlerFunc(handlerfunc)
	go func() {
		h := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handler}
		go func() {
			if err := h.ListenAndServe(); err != nil {
				if err != http.ErrServerClosed {
					panic(err)
				}
			}
		}()

		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		h.Shutdown(ctx)
		cancel()

	}()
	return port, bodychan
}
