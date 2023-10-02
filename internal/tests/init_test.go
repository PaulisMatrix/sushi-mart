package tests

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"sushi-mart/common"
	"sushi-mart/internal/database"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang/mock/gomock"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DatabaseSuite struct {
	suite.Suite
	*require.Assertions
	ctrl               *gomock.Controller
	config             *common.Config
	dockerTestPool     *dockertest.Pool
	dockerTestResource *dockertest.Resource
	postgresDB         *database.Postgres
	queries            *database.Queries
}

func TestDatabaseSuite(t *testing.T) {
	suite.Run(t, new(DatabaseSuite))
}

func (db *DatabaseSuite) SetupSuite() {
	db.Assertions = require.New(db.Suite.T())
	db.ctrl = gomock.NewController(db.T())
	db.config = common.GetConfig()
	db.dockerTestPool, db.dockerTestResource = initDockerTest(db.config)
	db.postgresDB, _ = database.NewPostgres(db.config.PgTestDbName, db.config.PgUser, db.config.PgPass)
	db.queries = database.New(db.postgresDB.DB)
}

func (db *DatabaseSuite) TearDownSuite() {
	db.ctrl = nil
	//remove the container and its resources if any.
	purgeResource(db.dockerTestPool, db.dockerTestResource)
	db.dockerTestPool = nil
	db.dockerTestResource = nil
	db.postgresDB.DB.Close()
	db.queries = nil
}

func initDockerTest(config *common.Config) (*dockertest.Pool, *dockertest.Resource) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.BuildAndRunWithBuildOptions(&dockertest.BuildOptions{
		Dockerfile: "Dockerfile",
		ContextDir: ".",
		BuildArgs: []docker.BuildArg{
			{
				Name:  "PG_USER",
				Value: fmt.Sprintf("POSTGRES_USER=%s", config.PgUser),
			},
			{
				Name:  "PG_PASSWORD",
				Value: fmt.Sprintf("POSTGRES_PASSWORD=%s", config.PgPass),
			},
			{
				Name:  "PGDB",
				Value: fmt.Sprintf("POSTGRES_DB=%s", config.PgTestDbName),
			},
		},
	}, &dockertest.RunOptions{Name: "testimage", ExposedPorts: []string{"3000"}}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	//currentcont, err := pool.CurrentContainer()
	//contname, ispresent := pool.ContainerByName()
	//fmt.Println("current container running err", err)
	//fmt.Println("current container running", currentcont)
	//if err != nil {
	//	os.Exit(1)
	//}
	hostAndPort := "localhost:3000"
	databaseUrl := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", config.PgUser, config.PgPass, hostAndPort, config.PgTestDbName)

	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err := sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	return pool, resource
}

func purgeResource(pool *dockertest.Pool, resource *dockertest.Resource) {
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func getHostAndPort(resource *dockertest.Resource, id string) string {
	dockerURL := os.Getenv("DOCKER_HOST")
	if dockerURL == "" {
		return resource.GetHostPort(id)
	}
	u, err := url.Parse(dockerURL)
	if err != nil {
		panic(err)
	}
	fmt.Println("host and post", u.Hostname(), resource.GetPort(id))
	return u.Hostname() + ":" + resource.GetPort(id)
}
