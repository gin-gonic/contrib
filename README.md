# gin-gonic/contrib 

[![Build Status](https://travis-ci.org/gin-gonic/contrib.svg)](https://travis-ci.org/gin-gonic/contrib)

Here you'll find middleware ready to  use with [Gin Framework](https://github.com/gin-gonic/gin). Submit your pull request, either with the package in a folder, or by adding a link to this `README.md`.

If adding a package directly, don't forget to create a `README.md` inside with author name.
If adding a link to your own repository, please follow this example:

```
+ nameOfMiddleware (https://github.com/yourusername/yourrepo)
```

Each author is responsible for maintaining their own code, although if you submit as a package, you allow the community to fix it. You can also submit a pull request to fix an existing package.
  
## List of external middleware

+ [goswag](https://github.com/diegoclair/goswag) - Wrapper to create swagger docs easily from your endpoints 
+ [RestGate](https://github.com/pjebs/restgate) - Secure authentication for REST API endpoints
+ [staticbin](https://github.com/olebedev/staticbin) - middleware/handler for serving static files from binary data
+ [gin-cachecontrol](https://github.com/joeig/gin-cachecontrol) - Cache-Control middleware
+ [gin-cors](https://github.com/gin-contrib/cors) - Official CORS gin's middleware
+ [gin-csrf](https://github.com/utrack/gin-csrf) - CSRF protection
+ [gin-health](https://github.com/utrack/gin-health) - middleware that simplifies stat reporting via [gocraft/health](https://github.com/gocraft/health)
+ [gin-merry](https://github.com/utrack/gin-merry) - middleware for pretty-printing [merry](https://github.com/ansel1/merry) errors with context
+ [gin-revision](https://github.com/appleboy/gin-revision-middleware) - Revision middleware for Gin framework
+ [gin-jwt](https://github.com/appleboy/gin-jwt) - JWT Middleware for Gin Framework
+ [gin-sessions](https://github.com/kimiazhu/ginweb-contrib/tree/master/sessions) - session middleware based on mongodb and mysql
+ [gin-location](https://github.com/drone/gin-location) - middleware for exposing the server's hostname and scheme
+ [gin-nice-recovery](https://github.com/ekyoung/gin-nice-recovery) - panic recovery middleware that lets you build a nicer user experience
+ [gin-limiter](https://github.com/davidleitw/gin-limiter) - A simple gin middleware for ip limiter based on redis.
+ [gin-limit](https://github.com/aviddiviner/gin-limit) - limits simultaneous requests; can help with high traffic load
+ [gin-limit-by-key](https://github.com/yangxikun/gin-limit-by-key) - An in-memory middleware to limit access rate by custom key and rate.
+ [ez-gin-template](https://github.com/michelloworld/ez-gin-template) - easy template wrap for gin
+ [gin-hydra](https://github.com/janekolszak/gin-hydra) - [Hydra](https://github.com/ory-am/hydra) middleware for Gin
+ [gin-glog](https://github.com/zalando/gin-glog) - meant as drop-in replacement for Gin's default logger
+ [gin-gomonitor](https://github.com/zalando/gin-gomonitor) - for exposing metrics with Go-Monitor
+ [gin-oauth2](https://github.com/zalando/gin-oauth2) - for working with OAuth2
+ [static](https://github.com/hyperboloide/static) An alternative static assets handler for the gin framework.
+ [xss-mw](https://github.com/dvwright/xss-mw) - XssMw is a middleware designed to "auto remove XSS" from user submitted input
+ [gin-helmet](https://github.com/danielkov/gin-helmet) - Collection of simple security middleware.
+ [gin-jwt-session](https://github.com/ScottHuangZL/gin-jwt-session) - middleware to provide JWT/Session/Flashes, easy to use while also provide options for adjust if necessary. Provide sample too.
+ [goview](https://github.com/foolin/goview) - a lightweight, minimalist and idiomatic template library
+ [ginvalidator](https://github.com/bube054/ginvalidator) - A robust and lightweight validation middleware for the Gin framework, simplifying input validation with customizable rules and error handling.
+ [gin-redis-ip-limiter](https://github.com/Salvatore-Giordano/gin-redis-ip-limiter) - Request limiter based on ip address. It works with redis and with a sliding-window mechanism.
+ [gin-method-override](https://github.com/bu/gin-method-override) - Method override by POST form param `_method`, inspired by Ruby's same name rack
+ [gin-access-limit](https://github.com/bu/gin-access-limit) - An access-control middleware by specifying allowed source CIDR notations.
+ [gin-session](https://github.com/go-session/gin-session) - Session middleware for Gin
+ [gin-stats](https://github.com/semihalev/gin-stats) - Lightweight and useful request metrics middleware
+ [gin-statsd](https://github.com/amalfra/gin-statsd) - A Gin middleware for reporting to statsd deamon
+ [gin-health-check](https://github.com/RaMin0/gin-health-check) - A health check middleware for Gin
+ [gin-session-middleware](https://github.com/go-session/gin-session) - A efficient, safely and easy-to-use session library for Go.
+ [ginception](https://github.com/kubastick/ginception) - Nice looking exception page
+ [gin-inspector](https://github.com/fatihkahveci/gin-inspector) - Gin middleware for investigating http request.
+ [go-gin-prometheus](https://github.com/zsais/go-gin-prometheus) - Gin Prometheus metrics exporter
+ [ginprom](https://github.com/chenjiandongx/ginprom) - Prometheus metrics exporter for Gin
+ [gin-go-metrics](https://github.com/bmc-toolbox/gin-go-metrics) - Gin middleware to gather and store metrics using [rcrowley/go-metrics](https://github.com/rcrowley/go-metrics)
+ [ginrpc](https://github.com/xxjwxc/ginrpc) - Gin middleware/handler auto binding tools. support object register by annotated route like beego
+ [logging](https://github.com/axiaoxin-com/logging#gin-middleware-ginlogger) - logging provide GinLogger uses zap to log detailed access logs in JSON or text format with trace id, supports flexible and rich configuration, and supports automatic reporting of log events above error level to sentry
+ [milogo](https://github.com/manuelarte/milogo) - Field selection patten for Gin
+ [ratelimiter](https://github.com/axiaoxin-com/ratelimiter) - Gin middleware for token bucket ratelimiter.
+ [servefiles](https://github.com/rickb777/servefiles) - serving static files with performance-enhancing cache control headers; also handles gzip & brotli compressed files 
+ [gin-brotli](https://github.com/anargu/gin-brotli) - Gin middleware to enable Brotli compression
+ [gin-jwt-cognito](https://github.com/akhettar/gin-jwt-cognito) - validate jwt tokens issued by [amazon cognito](https://aws.amazon.com/cognito/)
+ [gin-pagination](https://github.com/webstradev/gin-pagination) - A simple and customizable pagination middleware for Gin
+ [gin-opentelemetry](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation/github.com/gin-gonic/gin/otelgin) - Gin middleware for OpenTelemetry, you can also use OpenTelemetry via [Alibaba Compile-Time Auto Instrumentation](https://github.com/alibaba/opentelemetry-go-auto-instrumentation) or [OpenTelemetry Auto Instrumentation using eBPF](https://github.com/open-telemetry/opentelemetry-go-instrumentation) without changing any code
+ [scs_gin_adapter](https://github.com/39george/scs_gin_adapter) gin adapter for [SCS](https://github.com/alexedwards/scs) HTTP Session Management
+ [apitally](https://github.com/apitally/apitally-go) - Gin middleware for API monitoring, analytics and request logging with [Apitally](https://apitally.io/gin)
