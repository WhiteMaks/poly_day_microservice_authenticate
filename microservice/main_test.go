package microservice

import (
	"github.com/gin-gonic/gin"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.ReleaseMode)

	os.Exit(m.Run())
}
