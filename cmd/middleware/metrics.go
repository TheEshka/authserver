package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//InitPrometheus : add custom metrics from standardMetrics
func InitPrometheus() *Prometheus {
	p := &Prometheus{
		standardMetrics: defaultStandardMetrics,
		metricsPath:     defaultMetricPath,
	}

	for _, metric := range p.standardMetrics {
		prometheus.MustRegister(metric)
	}
	return p
}

//Use adds the middleware to a gin engine.
func (p *Prometheus) Use(e *gin.Engine) {
	e.Use(p.HandlerFunc())
	e.GET(p.metricsPath, prometheusHandler())

}

//HandlerFunc : function for middleware
func (p *Prometheus) HandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		troute := c.Request.URL.String()
		if troute == p.metricsPath {
			c.Next()
			return
		}
		c.Next()
		switch troute {
		case "/user":
			routeUser.WithLabelValues(strconv.Itoa(c.Writer.Status()), c.Request.Method).Inc()
			break
		case "/auth":
			routeAuth.WithLabelValues(strconv.Itoa(c.Writer.Status()), c.Request.Method).Inc()
			break
		}
	}
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

//Prometheus :
type Prometheus struct {
	standardMetrics []prometheus.Collector

	metricsPath string
}

var defaultMetricPath = "/metrics"

var defaultStandardMetrics = []prometheus.Collector{
	routeUser,
	routeAuth,
}

var routeUser = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "requests_to_route_USER",
		Help: "How many HTTP requests processed on rouse /user, partitioned by status code and HTTP method.",
	},
	[]string{"code", "method"},
)

var routeAuth = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "requests_to_route_AUTH",
		Help: "How many HTTP requests processed on rouse /auth, partitioned by status code and HTTP method.",
	},
	[]string{"code", "method"},
)
