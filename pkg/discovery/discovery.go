package discovery

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Registry interface {
	Register(ctx context.Context, serviceName, instanceID, hostPort string) error
	Deregister(ctx context.Context, serviceName, instanceID string) error
	ServiceAddresses(ctx context.Context, serviceID string) ([]string, error)
	ReportHealthyState(instanceID, serviceName string) error
}

var ErrNotFound = errors.New("no service address found")

func GenerateInstanceID(serviceName string) string {
	return fmt.Sprintf("%s-%d", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}
