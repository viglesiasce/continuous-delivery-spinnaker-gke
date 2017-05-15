package main

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/compute/metadata"
)

type CommonService interface {
	MetaData() *Instance
	Version() string
	Error() error
	Health() int
}

type commonService struct{}

func (commonService) Version() string {
	return version
}

func (commonService) MetaData() *Instance {
	return newInstance()
}

func (commonService) Health() int {
	return http.StatusOK
}

func (commonService) Error() error {
	message := "Unable to perform your request because of reasons"
	panic(message)
	return fmt.Errorf(message)
}

type Instance struct {
	Id         string
	Name       string
	Version    string
	Hostname   string
	Zone       string
	Project    string
	InternalIP string
	ExternalIP string
	LBRequest  string
	ClientIP   string
	Error      string
}

type assigner struct {
	err error
}

func (a *assigner) assign(getVal func() (string, error)) string {
	if a.err != nil {
		return ""
	}
	s, err := getVal()
	if err != nil {
		a.err = err
	}
	return s
}

func newInstance() *Instance {
	var i = new(Instance)
	if !metadata.OnGCE() {
		i.Error = "Not running on GCE"
		return i
	}

	a := &assigner{}
	i.Id = a.assign(metadata.InstanceID)
	i.Zone = a.assign(metadata.Zone)
	i.Name = a.assign(metadata.InstanceName)
	i.Hostname = a.assign(metadata.Hostname)
	i.Project = a.assign(metadata.ProjectID)
	i.InternalIP = a.assign(metadata.InternalIP)
	i.ExternalIP = a.assign(metadata.ExternalIP)
	// TODO Add back the version
	//i.Version = version

	if a.err != nil {
		i.Error = a.err.Error()
	}
	return i
}
