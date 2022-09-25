package rest

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

type RestfulAction int

// Return the http method string for the restful action
// to aid in registering the handler in the router group
func (a RestfulAction) String() string {
    switch a {
    case RestfulActionGetById:
        return http.MethodGet
    case RestfulActionList:
        return http.MethodGet
    case RestfulActionCreate:
        return http.MethodPost
    case RestfulActionUpdateById:
        return http.MethodPut
    case RestfulActionDeleteById:
        return http.MethodDelete
    }
    return "not-found"
}

const (
    RestfulActionGetById RestfulAction = iota
    RestfulActionList
    RestfulActionCreate
    RestfulActionUpdateById
    RestfulActionDeleteById
    RestfulActionExistsById
)

var (
    identifierRequiredSet = map[RestfulAction]struct{}{
        RestfulActionGetById:    {},
        RestfulActionUpdateById: {},
        RestfulActionDeleteById: {},
        RestfulActionExistsById: {},
    }
    identifierNotRequiredSet = map[RestfulAction]struct{}{
        RestfulActionList:   {},
        RestfulActionCreate: {},
    }
)

// RegisterRestfulHandlers registers one handler per restful action.
func RegisterRestfulHandlers(
    routerGroup *gin.RouterGroup,
    model string,
    handlers map[RestfulAction]gin.HandlerFunc,
    modelNameSanitizers ...ModelNameSanitizer) {

    if modelNameSanitizers != nil {
        for _, sanitizeModelNameFunc := range modelNameSanitizers {
            model = sanitizeModelNameFunc(model)
        }
    }

    var (
        baseUri       = "/" + model
        baseUriWithId = baseUri + "/:id"
    )

    for restfulAction, requestHandler := range handlers {
        httpAction := restfulAction.String()
        if _, ok := identifierNotRequiredSet[restfulAction]; ok {
            routerGroup.Handle(httpAction, baseUri, requestHandler)
        }
        if _, ok := identifierRequiredSet[restfulAction]; ok {
            routerGroup.Handle(httpAction, baseUriWithId, requestHandler)
        }
    }
}
