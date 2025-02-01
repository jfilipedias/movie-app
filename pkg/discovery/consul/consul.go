package consul

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	capi "github.com/hashicorp/consul/api"
	"github.com/jfilipedias/movie-app/pkg/discovery"
)

type Registry struct {
	client *capi.Client
}

func NewRegistry(addr string) (*Registry, error) {
	cfg := capi.DefaultConfig()
	cfg.Address = addr

	client, err := capi.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Registry{client}, nil
}

func (r *Registry) Register(ctx context.Context, instanceID, serviceName, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("hostPort must be in a form of <host>:<port>, example: localhost:8081")
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}

	return r.client.Agent().ServiceRegister(&capi.AgentServiceRegistration{
		ID:      instanceID,
		Name:    serviceName,
		Address: parts[0],
		Port:    port,
		Check:   &capi.AgentServiceCheck{CheckID: instanceID, TTL: "5s"},
	})
}

func (r *Registry) Deregister(ctx context.Context, _ string, instanceID string) error {
	return r.client.Agent().ServiceDeregister(instanceID)
}

func (r *Registry) ServicesAddresses(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, discovery.ErrNotFound
	}

	var res []string
	for _, e := range entries {
		res = append(res, fmt.Sprintf("%s:%d", e.Service.Address, e.Service.Port))
	}
	return res, nil
}

func (r *Registry) ReportHealthyState(_, instanceID string) error {
	return r.client.Agent().PassTTL(instanceID, "")
}
