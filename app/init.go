package app

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/ederavilaprado/golang-web-architecture-template/errors"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/access"
	"github.com/go-ozzo/ozzo-routing/fault"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/ogier/pflag"
)

// ConfigFile is the config file when specified, otherwhise the config will be loaded from the default place.
// See "app/config.go" for more information.
var ConfigFile *string

// LoadAppFlags load all the app flags
func LoadAppFlags() {
	configFile := pflag.StringP("config", "c", "", "Config file for the server")
	showVersion := pflag.BoolP("version", "v", false, "Version of the server")
	pflag.Parse()
	// print version with the flag: --version
	if *showVersion {
		fmt.Printf("v%s\n", Version)
		os.Exit(0)
	}
	// export configFile if present
	if *configFile != "" {
		ConfigFile = configFile
	}
}

// Init returns a middleware that prepares the request context and processing environment.
// The middleware will populate RequestContext, handle possible panics and errors from the processing
// handlers, and add an access log entry.
func Init(logger *logrus.Logger) routing.Handler {
	return func(rc *routing.Context) error {
		now := time.Now()

		rc.Response = &access.LogResponseWriter{
			ResponseWriter: rc.Response,
			Status:         http.StatusOK,
			BytesWritten:   0,
		}

		ac := newRequestScope(now, logger, rc.Request)
		rc.Set("Context", ac)

		fault.Recovery(ac.Errorf, convertError)(rc)
		logAccess(rc, ac.Infof, ac.Now())

		return nil
	}
}

// GetRequestScope returns the RequestScope of the current request.
func GetRequestScope(c *routing.Context) RequestScope {
	return c.Get("Context").(RequestScope)
}

// logAccess logs a message describing the current request.
func logAccess(c *routing.Context, logFunc access.LogFunc, start time.Time) {
	rw := c.Response.(*access.LogResponseWriter)
	elapsed := float64(time.Now().Sub(start).Nanoseconds()) / 1e6
	requestLine := fmt.Sprintf("%s %s %s", c.Request.Method, c.Request.URL.Path, c.Request.Proto)
	logFunc(`[%.3fms] %s %d %d`, elapsed, requestLine, rw.Status, rw.BytesWritten)
}

// convertError converts an error into an APIError so that it can be properly sent to the response.
// You may need to customize this method by adding conversion logic for more error types.
func convertError(c *routing.Context, err error) error {
	if err == sql.ErrNoRows {
		return errors.NotFound("the requested resource")
	}
	switch err.(type) {
	case *errors.APIError:
		return err
	case validation.Errors:
		return errors.InvalidData(err.(validation.Errors))
	case routing.HTTPError:
		switch err.(routing.HTTPError).StatusCode() {
		case http.StatusUnauthorized:
			return errors.Unauthorized(err.Error())
		case http.StatusNotFound:
			return errors.NotFound("the requested resource")
		}
	}
	return errors.InternalServerError(err)
}