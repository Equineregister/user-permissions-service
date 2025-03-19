package testdatabase

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"time"

	"github.com/Equineregister/user-permissions-service/pkg/migrations"
	"github.com/Equineregister/user-permissions-service/pkg/postgresutil"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	DbName = "test_db"
	DbUser = "test_user"     // nolint
	DbPass = "test_password" //nolint
)

type TestDatabase struct {
	DB        *pgxpool.Pool
	DBAddress string
	container testcontainers.Container
}

func NewTestDatabase(ctx context.Context, migrationsFS fs.FS) (*TestDatabase, error) {
	// setup db container
	container, dbConn, err := createContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("createContainer: %w", err)
	}

	// migrate db schema
	if err := migrations.Migrate(ctx, dbConn, migrationsFS); err != nil {
		return nil, fmt.Errorf("migrate fs: %w", err)
	}

	return &TestDatabase{
		container: container,
		DB:        dbConn,
	}, nil
}

func (tdb *TestDatabase) TearDown() {
	tdb.DB.Close()

	// remove test container
	_ = tdb.container.Terminate(context.Background())
}

// type containerLogger struct {
// 	enabled bool
// }

// func (cl *containerLogger) Printf(format string, args ...interface{}) {
// 	if cl.enabled {
// 		log.Printf(format, args...)
// 	}
// }

func createContainer(ctx context.Context) (testcontainers.Container, *pgxpool.Pool, error) {

	env := map[string]string{
		"POSTGRES_PASSWORD": DbPass,
		"POSTGRES_USER":     DbUser,
		"POSTGRES_DB":       DbName,
	}
	port := "5432/tcp"

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:         "postgres-test-" + uuid.New().String()[0:6],
			Image:        "postgres:14-alpine",
			ExposedPorts: []string{port},
			Env:          env,
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	})
	if err != nil {
		return container, nil, fmt.Errorf("failed to start container: %v", err)
	}
	time.Sleep(1 * time.Second)

	p, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return container, nil, fmt.Errorf("failed to get container external port: %v", err)
	}

	log.Println("postgres container ready and running at port: ", p.Port())
	// TODO: Look into removing this sleep, and look for hooks to wait for container to be ready

	dbHost := "localhost"
	db, _, err := postgresutil.Connect(ctx, p.Int(), dbHost, DbUser, DbPass, DbName)
	if err != nil {
		return container, nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return container, db, nil
}
