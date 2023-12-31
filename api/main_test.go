package api

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	db "github.com/naderSameh/billing_system/db/sqlc"
	"github.com/naderSameh/billing_system/util"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	server, err := NewServer(store)
	require.NoError(t, err)

	server1 := &http.Server{
		Addr:    ":8080",
		Handler: server.router,
	}

	go server.Start(":8080")
	err = server1.Close()
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	err := util.Loadconfig("../.")
	if err != nil {
		fmt.Print("cannot load config")
	}

	os.Exit(m.Run())
}
