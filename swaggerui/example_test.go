package swaggerui_test

import (
	"net/http"
	"strings"

	"github.com/fiorix/go-swagger-ui/swaggerui"
)

func Example() {
	http.Handle("/", swaggerui.Handler("/", strings.NewReader("{}")))
}
