package testutil

import (
	"context"
	"fmt"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupTestDB(t *testing.T) (string, func()) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		t.Fatal(err)
	}

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	dsn := fmt.Sprintf("postgres://testuser:testpass@%s:%d/testdb?sslmode=disable", host, port.Int())

	return dsn, func() {
		container.Terminate(ctx)
	}
}

func Connect() {

}
