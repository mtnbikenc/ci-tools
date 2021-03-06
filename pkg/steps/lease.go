package steps

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	"github.com/openshift/ci-tools/pkg/api"
	"github.com/openshift/ci-tools/pkg/junit"
	"github.com/openshift/ci-tools/pkg/lease"
	"github.com/openshift/ci-tools/pkg/results"
)

const DefaultLeaseEnv = "LEASED_RESOURCE"

var NoLeaseClientErr = errors.New("step needs a lease but no lease client provided")

type stepLease struct {
	api.StepLease
	resource string
}

// leaseStep wraps another step and acquires/releases one or more leases.
type leaseStep struct {
	client  *lease.Client
	leases  []stepLease
	wrapped api.Step

	// for sending heartbeats during lease acquisition
	namespace func() string
}

func LeaseStep(client *lease.Client, leases []api.StepLease, wrapped api.Step, namespace func() string) api.Step {
	ret := leaseStep{
		client:    client,
		wrapped:   wrapped,
		namespace: namespace,
	}
	for _, l := range leases {
		ret.leases = append(ret.leases, stepLease{StepLease: l})
	}
	return &ret
}

func (s *leaseStep) Inputs() (api.InputDefinition, error) {
	return s.wrapped.Inputs()
}

func (s *leaseStep) Validate() error {
	if s.client == nil {
		return NoLeaseClientErr
	}
	return nil
}

func (s *leaseStep) Name() string             { return s.wrapped.Name() }
func (s *leaseStep) Description() string      { return s.wrapped.Description() }
func (s *leaseStep) Requires() []api.StepLink { return s.wrapped.Requires() }
func (s *leaseStep) Creates() []api.StepLink  { return s.wrapped.Creates() }
func (s *leaseStep) Provides() api.ParameterMap {
	parameters := s.wrapped.Provides()
	if parameters == nil {
		parameters = api.ParameterMap{}
	}
	for i := range s.leases {
		l := &s.leases[i]
		parameters[l.Env] = func() (string, error) {
			chunks := strings.SplitN(l.resource, "--", 2)
			return chunks[0], nil
		}
	}
	return parameters
}

func (s *leaseStep) SubTests() []*junit.TestCase {
	if subTests, ok := s.wrapped.(subtestReporter); ok {
		return subTests.SubTests()
	}
	return nil
}

func (s *leaseStep) Run(ctx context.Context) error {
	return results.ForReason("utilizing_lease").ForError(s.run(ctx))
}

func (s *leaseStep) run(ctx context.Context) error {
	log.Printf("Acquiring leases for %q", s.Name())
	client := *s.client
	ctx, cancel := context.WithCancel(ctx)
	if err := acquireLeases(client, ctx, cancel, s.leases); err != nil {
		return err
	}
	wrappedErr := results.ForReason("executing_test").ForError(s.wrapped.Run(ctx))
	log.Printf("Releasing leases for %q", s.Name())
	releaseErr := results.ForReason("releasing_lease").ForError(releaseLeases(client, s.leases))

	// we want a sensible output error for reporting, so we bubble up these individually
	//if we can, as this is the only step that can have multiple errors
	if wrappedErr != nil && releaseErr == nil {
		return wrappedErr
	} else if wrappedErr == nil && releaseErr != nil {
		return releaseErr
	} else {
		return utilerrors.NewAggregate([]error{wrappedErr, releaseErr})
	}
}

func acquireLeases(
	client lease.Client,
	ctx context.Context,
	cancel context.CancelFunc,
	leases []stepLease,
) error {
	// Sort by resource type to avoid a(n unlikely and temporary) deadlock.
	var sorted []int
	for i := range leases {
		sorted = append(sorted, i)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return leases[i].ResourceType < leases[j].ResourceType
	})
	var errs []error
	for _, i := range sorted {
		l := &leases[i]
		log.Printf("Acquiring lease for %q", l.ResourceType)
		name, err := client.Acquire(l.ResourceType, ctx, cancel)
		if err != nil {
			if err == lease.ErrNotFound {
				printResourceMetrics(client, l.ResourceType)
			}
			errs = append(errs, results.ForReason(results.Reason("acquiring_lease:"+l.ResourceType)).WithError(err).Errorf("failed to acquire lease: %v", err))
			break
		}
		log.Printf("Acquired lease %q for %q", name, l.ResourceType)
		l.resource = name
	}
	if errs != nil {
		if err := releaseLeases(client, leases); err != nil {
			errs = append(errs, fmt.Errorf("failed to release leases after acquisition failure: %w", err))
		}
	}
	return utilerrors.NewAggregate(errs)
}

func releaseLeases(client lease.Client, leases []stepLease) error {
	var errs []error
	for _, l := range leases {
		if l.resource == "" {
			continue
		}
		log.Printf("Releasing lease %q for %q", l.resource, l.ResourceType)
		if err := client.Release(l.resource); err != nil {
			errs = append(errs, err)
		}
	}
	return utilerrors.NewAggregate(errs)
}

func printResourceMetrics(client lease.Client, rtype string) {
	m, err := client.Metrics(rtype)
	if err != nil {
		log.Printf("warning: Could not get resource metrics: %v", err)
		return
	}
	log.Printf("error: Failed to acquire resource, current capacity: %d free, %d leased", m.Free, m.Leased)
}
