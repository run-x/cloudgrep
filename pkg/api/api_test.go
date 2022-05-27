package api

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore"
	"github.com/run-x/cloudgrep/pkg/datastore/testdata"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/version"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func PrepareApiUnitTest(t *testing.T) (*zap.Logger, datastore.Datastore, *gin.Engine) {
	ctx := context.Background()
	logger := zaptest.NewLogger(t)

	datastoreConfigs := config.Datastore{
		Type:           "sqlite",
		DataSourceName: "file::memory:",
	}
	cfg, err := config.GetDefault()
	assert.NoError(t, err)
	cfg.Datastore = datastoreConfigs

	ds, err := datastore.NewDatastore(ctx, cfg, zaptest.NewLogger(t))
	assert.NoError(t, err)

	router := gin.Default()
	SetupRoutes(router, cfg, logger, ds)
	return logger, ds, router
}

func TestHealthRoute(t *testing.T) {
	_, _, router := PrepareApiUnitTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"status\":\"All good!\"}", w.Body.String())
}

func TestHomeRoute(t *testing.T) {
	_, _, router := PrepareApiUnitTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))

	assert.Equal(t, 200, w.Code)
	assert.True(t, w.Body.Len() > 0)
}

func TestInfoRoute(t *testing.T) {
	_, _, router := PrepareApiUnitTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/info", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, 200, w.Code)
	assert.True(t, body["version"] == version.Version)
	assert.True(t, body["go_version"] == version.GoVersion)
	assert.True(t, body["git_sha"] == version.GitCommit)
	assert.True(t, body["build_time"] == version.BuildTime)
}

func TestStatsRoute(t *testing.T) {
	_, _, router := PrepareApiUnitTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/stats", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, 200, w.Code)
}

func TestResourcesRoute(t *testing.T) {
	ctx := context.Background()
	_, ds, router := PrepareApiUnitTest(t)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/resources", nil)
	assert.NoError(t, err)
	router.ServeHTTP(w, req)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	var body model.Resources
	err = json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, len(body), 0)

	resources := testdata.GetResources(t)
	assert.NotZero(t, len(resources))
	sort.Sort(model.ResourcesById(resources))

	//write the resources
	assert.NoError(t, ds.WriteResources(ctx, resources))
	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/api/resources", nil)
	assert.NoError(t, err)
	router.ServeHTTP(w, req)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	body = nil
	err = json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, len(body), 3)
	sort.Sort(model.ResourcesById(body))
	model.AssertEqualsResources(t, body, resources)
}
