package httpsvr

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
)

func headers(c *gin.Context) {
	c.Header("was_server", "was")
}

func Ping(ctx context.Context, c *gin.Context) {
	headers(c)
	span, _ := opentracing.StartSpanFromContext(ctx, "Ping")
	defer span.Finish()

	c.String(http.StatusOK, "pong")
}

func ProxyAir(ctx context.Context, c *gin.Context) {
	headers(c)
	span, _ := opentracing.StartSpanFromContext(ctx, "ProxyAir")
	defer span.Finish()

	endpoint := os.Getenv("AIR_SERVICE_ENDPOINT")
	if url, err := url.Parse("http://" + endpoint); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		r := httputil.NewSingleHostReverseProxy(url)
		allowCORS(ctx, c)
		r.ServeHTTP(c.Writer, c.Request)
	}

}

func allowCORS(ctx context.Context, c *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "allowCORS")
	defer span.Finish()

	// allow CORS
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

}
