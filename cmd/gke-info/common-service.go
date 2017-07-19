package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"os"

	"cloud.google.com/go/compute/metadata"
)

type CommonService interface {
	MetaData(r *http.Request) *Instance
	Version(r *http.Request) string
	Error(r *http.Request) error
	Health(r *http.Request) string
	Home(r *http.Request) string
}

type commonService struct {
	backendURL string
	sdc        *stackDriverClient
}

func (cs commonService) Version(r *http.Request) string {
	return version
}

func (cs commonService) MetaData(r *http.Request) *Instance {
	instance := newInstance(r.Context(), cs)
	raw, _ := httputil.DumpRequest(r, true)
	instance.LBRequest = string(raw)
	instance.ClientIP = r.RemoteAddr
	instance.Version = version
	instance.Color = "orange"
	instance.PodName = os.Getenv("HOSTNAME")
	return instance
}

func (cs commonService) Health(r *http.Request) string {
	return "ok"
}

func (cs commonService) Error(r *http.Request) error {
	message := "Unable to perform your request: " + r.URL.Query().Get("message")
	panic(message)
	return fmt.Errorf(message)
}

func (cs commonService) Home(r *http.Request) string {
	tpl := template.Must(template.New("out").Parse(html))
	req, _ := http.NewRequest(
		"GET",
		cs.backendURL,
		nil,
	)

	req = req.WithContext(r.Context())
	body := makeRequest(req, cs)
	i := &Instance{}
	err := json.Unmarshal([]byte(body), i)
	if err != nil {
		message := "Unable to unmarshall response: " + err.Error()
		panic(message)
	}
	var buf bytes.Buffer
	tpl.Execute(&buf, i)
	return buf.String()
}

type Instance struct {
	Name       string
	Color      string
	Version    string
	Zone       string
	Project    string
	InternalIP string
	ExternalIP string
	LBRequest  string
	ClientIP   string
	PodName    string
	Error      string
}

func getMetaData(ctx context.Context, cs commonService, path string) string {
	metaDataURL := "http://metadata/computeMetadata/v1/"
	req, _ := http.NewRequest(
		"GET",
		metaDataURL+path,
		nil,
	)
	req.Header.Add("Metadata-Flavor", "Google")
	req = req.WithContext(ctx)
	return string(makeRequest(req, cs))
}

func makeRequest(r *http.Request, cs commonService) []byte {
	transport := http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: &transport}
	// traceClient := cs.sdc.traceClient.NewHTTPClient(client)
	resp, err := client.Do(r)
	if err != nil {
		message := "Unable to call backend: " + err.Error()
		panic(message)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		message := "Unable to read response body: " + err.Error()
		panic(message)
	}
	return body
}

func newInstance(ctx context.Context, cs commonService) *Instance {
	var i = new(Instance)
	if !metadata.OnGCE() {
		i.Error = "Not running on GCE"
		return i
	}

	i.Error = "None"
	i.Zone = getMetaData(ctx, cs, "instance/zone")
	i.Name = getMetaData(ctx, cs, "instance/name")
	i.Project = getMetaData(ctx, cs, "project/project-id")
	i.InternalIP = getMetaData(ctx, cs, "instance/network-interfaces/0/ip")
	i.ExternalIP = getMetaData(ctx, cs, "instance/network-interfaces/0/access-configs/0/external-ip")

	return i
}
