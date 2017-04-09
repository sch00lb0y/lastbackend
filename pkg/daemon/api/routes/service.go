//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package routes

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/daemon/api/views/v1"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/errors"
	"github.com/lastbackend/lastbackend/pkg/util/converter"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type serviceCreateS struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Template    string               `json:"template"`
	Image       string               `json:"image"`
	Url         string               `json:"url"`
	Region      string               `json:"region"`
	Config      *types.ServiceConfig `json:"config,omitempty"`
	source      *types.ServiceSource
}

type resources struct {
	Region string `json:"region"`
	Memory int    `json:"memory"`
}

func (s *serviceCreateS) decodeAndValidate(reader io.Reader) *errors.Err {

	var ctx = c.Get()

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		ctx.Log.Error(err)
		return errors.New("user").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("service").IncorrectJSON(err)
	}

	if s.Template == "" && s.Image == "" && s.Url == "" {
		return errors.New("service").BadParameter("template,image,url")
	}

	if s.Template != "" {
		if s.Name == "" {
			s.Name = s.Template
		}
	}

	if s.Image != "" && s.Url == "" {
		source, err := converter.DockerNamespaceParse(s.Image)
		if err != nil {
			return errors.New("service").BadParameter("image")
		}

		if s.Name == "" {
			s.Name = source.Repo
		}

		s.source = &types.ServiceSource{
			Type:   types.SourceDockerType,
			Hub:    source.Hub,
			Owner:  source.Owner,
			Repo:   source.Repo,
			Branch: source.Branch,
		}
	}

	if s.Url != "" {
		if !validator.IsGitUrl(s.Url) {
			return errors.New("service").BadParameter("url")
		}

		source, err := converter.GitUrlParse(s.Url)
		if err != nil {
			return errors.New("service").BadParameter("url")
		}

		if s.Name == "" {
			s.Name = source.Repo
		}

		s.source = &types.ServiceSource{
			Type:   types.SourceGitType,
			Hub:    source.Hub,
			Owner:  source.Owner,
			Repo:   source.Repo,
			Branch: "master",
		}
	}

	s.Name = strings.ToLower(s.Name)

	if s.Name == "" {
		return errors.New("service").BadParameter("name")
	}

	s.Name = strings.ToLower(s.Name)

	if len(s.Name) < 4 && len(s.Name) > 64 && !validator.IsServiceName(s.Name) {
		return errors.New("service").BadParameter("name")
	}

	return nil
}

func ServiceCreateH(w http.ResponseWriter, r *http.Request) {

	var (
		session      *types.Session
		ctx          = c.Get()
		params       = utils.Vars(r)
		projectParam = params["project"]
	)

	ctx.Log.Debug("Create service handler")

	session = utils.Session(r)
	if session == nil {
		ctx.Log.Error(http.StatusText(http.StatusUnauthorized))
		errors.HTTP.Unauthorized(w)
		return
	}

	// request body struct
	rq := new(serviceCreateS)
	if err := rq.decodeAndValidate(r.Body); err != nil {
		ctx.Log.Error("Error: validation incomming data", err)
		errors.New("Invalid incomming data").Unknown().Http(w)
		return
	}

	service, err := ctx.Storage.Service().GetByName(session.Username, projectParam, rq.Name)
	if err != nil {
		ctx.Log.Error("Error: check exists by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if service != nil {
		errors.New("service").NotUnique("name").Http(w)
		return
	}

	// Load template from registry
	if rq.Template != "" {
		// TODO: send request for get template config from registry
		// TODO: Set service source with types.SourceTemplateType type field
	}

	// If you are not using a template, then create a standard configuration template
	//if tpl == nil {
	// TODO: Generate default template for service
	//}

	// Patch config if exists custom configurations
	if rq.Config != nil {
		// TODO: If have custom config, then need patch this config
	}

	rq.Config = types.ServiceConfig{}.GetDefault()

	service, err = ctx.Storage.Service().Insert(session.Username, projectParam, rq.Name, rq.Description, rq.source, rq.Config)
	if err != nil {
		ctx.Log.Error("Error: insert service to db", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewService(service).ToJson()
	if err != nil {
		ctx.Log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}
}

type serviceUpdateS struct {
	Description string               `json:"description"`
	Config      *types.ServiceConfig `json:"config,omitempty"`
	Domains     *[]string            `json:"domains,omitempty"`
}

func (s *serviceUpdateS) decodeAndValidate(reader io.Reader) *errors.Err {

	var ctx = c.Get()

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		ctx.Log.Error(err)
		return errors.New("user").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("service").IncorrectJSON(err)
	}

	return nil
}

func ServiceUpdateH(w http.ResponseWriter, r *http.Request) {

	var (
		ctx          = c.Get()
		session      *types.Session
		params       = utils.Vars(r)
		projectParam = params["project"]
		serviceParam = params["service"]
	)

	ctx.Log.Debug("Update service handler")

	session = utils.Session(r)
	if session == nil {
		ctx.Log.Error(http.StatusText(http.StatusUnauthorized))
		errors.HTTP.Unauthorized(w)
		return
	}

	// request body struct
	rq := new(serviceUpdateS)
	if err := rq.decodeAndValidate(r.Body); err != nil {
		ctx.Log.Error("Error: validation incomming data", err)
		errors.New("Invalid incomming data").Unknown().Http(w)
		return
	}

	service, err := ctx.Storage.Service().GetByName(session.Username, projectParam, serviceParam)
	if err != nil {
		ctx.Log.Error("Error: Get service by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	service.Description = rq.Description

	if rq.Config != nil {
		if err := service.Config.Update(rq.Config); err != nil {
			ctx.Log.Error("Error: update service config", err.Error())
			errors.New("service").BadParameter("config", err)
			return
		}
	}

	if rq.Domains != nil {
		service.Domains = *rq.Domains
	}

	service, err = ctx.Storage.Service().Update(session.Username, projectParam, service)
	if err != nil {
		ctx.Log.Error("Error: insert service to db", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewService(service).ToJson()
	if err != nil {
		ctx.Log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}
}

func ServiceListH(w http.ResponseWriter, r *http.Request) {

	var (
		ctx          = c.Get()
		session      *types.Session
		params       = utils.Vars(r)
		projectParam = params["project"]
	)

	ctx.Log.Debug("List service handler")

	session = utils.Session(r)
	if session == nil {
		ctx.Log.Error(http.StatusText(http.StatusUnauthorized))
		errors.HTTP.Unauthorized(w)
		return
	}

	serviceList, err := ctx.Storage.Service().ListByProject(session.Username, projectParam)
	if err != nil {
		ctx.Log.Error("Error: find service list by user", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewServiceList(serviceList).ToJson()
	if err != nil {
		ctx.Log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}
}

func ServiceInfoH(w http.ResponseWriter, r *http.Request) {
	var (
		session      *types.Session
		service      *types.Service
		ctx          = c.Get()
		params       = utils.Vars(r)
		projectParam = params["project"]
		serviceParam = params["service"]
	)

	ctx.Log.Debug("Get service handler")

	session = utils.Session(r)
	if session == nil {
		ctx.Log.Error(http.StatusText(http.StatusUnauthorized))
		errors.HTTP.Unauthorized(w)
		return
	}

	service, err := ctx.Storage.Service().GetByName(session.Username, projectParam, serviceParam)
	if err != nil {
		ctx.Log.Error("Error: find service by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if service == nil {
		errors.New("service").NotFound().Http(w)
		return
	}

	response, err := v1.NewService(service).ToJson()
	if err != nil {
		ctx.Log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}
}

func ServiceRemoveH(w http.ResponseWriter, r *http.Request) {
	var (
		ctx          = c.Get()
		session      *types.Session
		params       = utils.Vars(r)
		projectParam = params["project"]
		serviceParam = params["service"]
	)

	ctx.Log.Info("Remove service")

	session = utils.Session(r)
	if session == nil {
		ctx.Log.Error(http.StatusText(http.StatusUnauthorized))
		errors.HTTP.Unauthorized(w)
		return
	}

	// Todo: remove all activity by service name

	if err := ctx.Storage.Service().Remove(session.Username, projectParam, serviceParam); err != nil {
		ctx.Log.Error("Error: remove service from db", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}
}
