package echocustoms

import (
	"errors"
	"net/http"

	"github.com/dinno7/ride-sharing/shared/util"
	"github.com/labstack/echo/v5"
)

func CustomHTTPErrorHandler(c *echo.Context, err error) {
	if resp, uErr := echo.UnwrapResponse(c.Response()); uErr == nil {
		if resp.Committed {
			return // response has been already sent to the client by handler or some middleware
		}
	}

	code := http.StatusInternalServerError
	var sc echo.HTTPStatusCoder
	if errors.As(err, &sc) { // find error in an error chain that implements HTTPStatusCoder
		if tmp := sc.StatusCode(); tmp != 0 {
			code = tmp
		}
	}

	var cErr error
	if c.Request().Method == http.MethodHead {
		cErr = c.NoContent(code)
	} else {
		// errorPage := fmt.Sprintf("%d.html", code)
		// cErr = c.File(errorPage)
		cErr = c.JSON(code, util.NewErrorPayload(http.StatusText(code), err.Error()))
	}

	if cErr != nil {
		c.Logger().Error("failed to send error page to client", "error", errors.Join(err, cErr))
	}
}
