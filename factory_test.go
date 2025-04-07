package redisfactory

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

var testContainers = []struct {
	name  string
	image string
}{
	{"redis7", "docker.io/redis:7"},
	{"valkey7.2", "docker.io/valkey/valkey:7.2"},
	{"valkey8.0", "docker.io/valkey/valkey:8.0"},
	{"valkey8.1", "docker.io/valkey/valkey:8.1"},
}

func TestFromURL(t *testing.T) {
	for _, tt := range testContainers {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			container, err := redis.RunContainer(ctx,
				testcontainers.WithImage(tt.image),
				redis.WithLogLevel(redis.LogLevelVerbose),
			)

			if err != nil {
				t.Fatalf("Error starting test container: %v", err)
			}
			t.Cleanup(func() {
				if err := container.Terminate(ctx); err != nil {
					t.Fatalf("Error terminating test container: %v", err)
				}
			})

			connStr, err := container.ConnectionString(ctx)
			if err != nil {
				t.Fatalf("Error getting test container connection string: %v", err)
			}

			client, err := FromURL(ctx, 2, connStr)
			if err != nil {
				t.Fatalf("Error getting client: %v", err)
			}

			_ = client.Close()
		})
	}
}

func TestBadURL(t *testing.T) {
	ctx := context.Background()
	_, err := FromURL(ctx, 3, "")
	if err == nil {
		t.Fatalf("Should have timedout")
	}
}

func TestFromAddress(t *testing.T) {
	for _, tt := range testContainers {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()
			container, err := redis.RunContainer(ctx,
				testcontainers.WithImage(tt.image),
				redis.WithLogLevel(redis.LogLevelVerbose),
			)

			if err != nil {
				t.Fatalf("Error starting test container: %v", err)
			}
			t.Cleanup(func() {
				if err := container.Terminate(ctx); err != nil {
					t.Fatalf("Error terminating test container: %v", err)
				}
			})

			addr, err := container.Endpoint(ctx, "")
			if err != nil {
				t.Fatalf("Error getting test container connection string: %v", err)
			}

			client, err := FromAddress(ctx, 2, addr, "")
			if err != nil {
				t.Fatalf("Error getting client: %v", err)
			}

			_ = client.Close()
		})
	}
}

func TestClientNoConnect(t *testing.T) {
	ctx := context.Background()
	_, err := FromAddress(ctx, 3, "", "")
	if err == nil {
		t.Fatalf("Should have timedout")
	}
}

func TestClientNoConnectDelay(t *testing.T) {
	t.Setenv(initialDelayEnvKey, "25")
	t.Setenv(maxDelayEnvKey, "50")
	ctx := context.Background()
	_, err := FromAddress(ctx, 5, "", "")
	if err == nil {
		t.Fatalf("Should have timedout")
	}
}

func Test_envIntOrDefault(t *testing.T) {
	t.Setenv("A", "25")
	t.Setenv("B", "dog")

	type args struct {
		key          string
		defaultValue int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"simple", args{key: "A", defaultValue: 5}, 25},
		{"bad value", args{key: "B", defaultValue: 5}, 5},
		{"no value", args{key: "C", defaultValue: 5}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := envIntOrDefault(tt.args.key, tt.args.defaultValue); got != tt.want {
				t.Errorf("envIntOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}
