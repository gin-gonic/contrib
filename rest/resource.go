package rest

import (
	"github.com/gertd/go-pluralize"
	"github.com/gin-gonic/gin"
	"github.com/gobeam/stringy"
	"net/http"
	"regexp"
	"strings"
)

type ResourceHandlers map[string]gin.HandlerFunc

type ResourceOptions struct {
	Param   string
	Exclude []string
}

func Resource(group *gin.RouterGroup, handlers ResourceHandlers, options ResourceOptions) *gin.RouterGroup {
	name := getNameFromPath(group.BasePath())
	options.Param = getRouteParam(name, options.Param)

	handleResourceMethod(group, &handlers, options)

	return group
}

type routeMethod struct {
	Name       string
	HttpMethod string
	Param      string
}

type routeMethodList map[string]routeMethod

func allowedRouteMethods(param string) routeMethodList {
	return routeMethodList{
		"Index":  routeMethod{Name: "Index", HttpMethod: http.MethodGet},
		"Create": routeMethod{Name: "Create", HttpMethod: http.MethodPost},
		"Show":   routeMethod{Name: "Show", HttpMethod: http.MethodGet, Param: param},
		"Update": routeMethod{Name: "Update", HttpMethod: http.MethodPut, Param: param},
		"Delete": routeMethod{Name: "Delete", HttpMethod: http.MethodDelete, Param: param},
	}
}

func handleResourceMethod(group *gin.RouterGroup, handlers *ResourceHandlers, options ResourceOptions) {
	for _, method := range allowedRouteMethods(options.Param) {
		handler := (*handlers)[method.Name]

		if handler != nil && !isExcludedMethod(method.Name, options.Exclude) {
			param := ""
			if method.Param != "" {
				param += "/:" + method.Param
			}

			group.Handle(
				method.HttpMethod,
				param,
				handler,
			)
		}
	}
}

func isExcludedMethod(methodName string, excluded []string) bool {
	if len(excluded) < 1 {
		return false
	}

	for _, em := range excluded {
		if em == methodName {
			return true
		}
	}

	return false
}

func getNameFromPath(path string) (name string) {
	if !strings.Contains(path, "/") || path == "" {
		return
	}

	r := regexp.MustCompile(`\w+$`)
	name = r.FindStringSubmatch(path)[0]

	return
}

func getRouteParam(name string, param string) string {
	if param != "" {
		return param
	}

	return getParamFromName(name)
}

func getParamFromName(name string) string {
	if name == "" {
		return "id"
	}

	param := stringy.New(name).KebabCase().ToLower()
	paramSingular := pluralize.NewClient().Singular(param)

	return paramSingular
}
