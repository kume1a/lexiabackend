package helpers

import (
	"context"
	"fmt"
	"lexia/ent"
	"lexia/internal/modules"
	"lexia/internal/shared"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type E2ETestSuite struct {
	suite.Suite
	ctx               context.Context
	postgresContainer *postgres.PostgresContainer
	dbClient          *ent.Client
	server            *gin.Engine
	testServer        *httptest.Server
	originalDbConnStr string
}

func (suite *E2ETestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	suite.ctx = context.Background()

	suite.originalDbConnStr = os.Getenv(shared.EnvDbConnectionString)

	var err error
	suite.postgresContainer, err = postgres.Run(suite.ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Minute),
		),
	)
	suite.Require().NoError(err)

	connStr, err := suite.postgresContainer.ConnectionString(suite.ctx, "sslmode=disable")
	suite.Require().NoError(err)

	os.Setenv(shared.EnvEnvironment, "test")
	os.Setenv(shared.EnvDbConnectionString, connStr)
	os.Setenv(shared.EnvPort, "8080")
	os.Setenv(shared.EnvAccessTokenSecret, "test-jwt-secret-key-for-testing")
	os.Setenv(shared.EnvAccessTokenExpSeconds, "3600")

	suite.dbClient, err = ent.Open("postgres", connStr)
	suite.Require().NoError(err)

	err = suite.dbClient.Schema.Create(suite.ctx)
	suite.Require().NoError(err)

	resouceConfig := &shared.ResourceConfig{
		DB: suite.dbClient,
	}

	apiCfg := shared.ApiConfig{
		ResourceConfig: resouceConfig,
	}

	suite.server, err = modules.CreateWebserver(&apiCfg)
	suite.Require().NoError(err)

	suite.testServer = httptest.NewServer(suite.server)
}

func (suite *E2ETestSuite) TearDownSuite() {
	if suite.testServer != nil {
		suite.testServer.Close()
	}

	if suite.dbClient != nil {
		suite.dbClient.Close()
	}

	if suite.postgresContainer != nil {
		err := suite.postgresContainer.Terminate(suite.ctx)
		suite.Require().NoError(err)
	}

	if suite.originalDbConnStr != "" {
		os.Setenv(shared.EnvDbConnectionString, suite.originalDbConnStr)
	} else {
		os.Unsetenv(shared.EnvDbConnectionString)
	}
}

func (suite *E2ETestSuite) SetupTest() {
	suite.cleanupDatabase()
}

func (suite *E2ETestSuite) TearDownTest() {
	suite.cleanupDatabase()
}

func (suite *E2ETestSuite) cleanupDatabase() {
	_, err := suite.dbClient.Word.Delete().Exec(suite.ctx)
	suite.Require().NoError(err)

	_, err = suite.dbClient.Folder.Delete().Exec(suite.ctx)
	suite.Require().NoError(err)

	_, err = suite.dbClient.User.Delete().Exec(suite.ctx)
	suite.Require().NoError(err)
}

func (suite *E2ETestSuite) GetTestServerURL() string {
	return suite.testServer.URL
}

func (suite *E2ETestSuite) GetAPIURL(endpoint string) string {
	return fmt.Sprintf("%s/api/v1%s", suite.testServer.URL, endpoint)
}

func (suite *E2ETestSuite) GetDBClient() *ent.Client {
	return suite.dbClient
}

func (suite *E2ETestSuite) GetContext() context.Context {
	return suite.ctx
}

func RunE2ETestSuite(t *testing.T, testSuite suite.TestingSuite) {
	suite.Run(t, testSuite)
}
