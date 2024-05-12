package service

import "net/http"

func addRoutes(s *Service) *[]RouteHandler {
	return &[]RouteHandler{
		{
			Route:   "/api/v1/work",
			Method:  http.MethodGet,
			Handler: s.getWorkHandler,
		},
		{
			Route:   "/api/v1/work",
			Method:  http.MethodPost,
			Handler: s.postWorkHandler,
		},
		{
			Route:   "/api/v1/work",
			Method:  http.MethodPut,
			Handler: s.updateWorkHandler,
		},
		{
			Route:   "/api/v1/work",
			Method:  http.MethodDelete,
			Handler: s.deleteWorkHandler,
		},
		{
			Route:   "/api/v1/skillsTools",
			Method:  http.MethodGet,
			Handler: s.getSkillsToolsHandler,
		},
		{
			Route:   "/api/v1/skillsTools",
			Method:  http.MethodPost,
			Handler: s.postSkillsToolsHandler,
		},
		{
			Route:   "/api/v1/skillsTools",
			Method:  http.MethodPut,
			Handler: s.updateSkillsToolsHandler,
		},
		{
			Route:   "/api/v1/skillsTools",
			Method:  http.MethodDelete,
			Handler: s.deleteSkillsToolsHandler,
		},
		{
			Route:   "/api/v1/projects",
			Method:  http.MethodGet,
			Handler: s.getProjectsHandler,
		},
		{
			Route:   "/api/v1/projects",
			Method:  http.MethodPost,
			Handler: s.postProjectsHandler,
		},
		{
			Route:   "/api/v1/projects",
			Method:  http.MethodPut,
			Handler: s.updateProjectsHandler,
		},
	}
}
