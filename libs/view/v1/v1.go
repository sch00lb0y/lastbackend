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

package v1

import (
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/libs/view/v1/project"
	"github.com/lastbackend/lastbackend/libs/view/v1/service"
	"github.com/lastbackend/lastbackend/libs/view/v1/user"
)

func NewUser(obj *model.User) *user.User {
	return user.New(obj)
}

func NewProject(obj *model.Project) *project.Project {
	return project.New(obj)
}

func NewProjectList(obj *model.ProjectList) *project.ProjectList {
	return project.NewList(obj)
}

func NewService(obj *model.Service) *service.Service {
	return service.New(obj)
}

func NewServiceList(obj *model.ServiceList) *service.ServiceList {
	return service.NewList(obj)
}
