package etcd

import (
	"context"
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	etcd3 "go.etcd.io/etcd/client/v3"
	"math/rand"
	"time"
)

var _ g.ServiceRegistry = (*Registry)(nil)
var _ g.ServiceDiscovery = (*Registry)(nil)

type Registry struct {
	gone.Flag

	logger gone.Logger   `gone:"*"`
	client *etcd3.Client `gone:"*"`

	dialTimeout  time.Duration `gone:"config,etcd.dial-timeout=5s"`
	keepaliveTTL time.Duration `gone:"config,etcd.keepalive-ttl=10s"`

	lease etcd3.Lease
}

func (r *Registry) Register(instance g.Service) error {
	return gone.ToError(r.doRegisterLease(context.Background(), instance))
}

// doKeepAlive continuously keeps alive the lease from ETCD.
func (r *Registry) doKeepAlive(
	service g.Service, leaseID etcd3.LeaseID, keepAliceCh <-chan *etcd3.LeaseKeepAliveResponse,
) {
	var ctx = context.Background()
	for {
		select {
		case <-r.client.Ctx().Done():
			r.logger.Infof("keepalive done for lease id: %d", leaseID)
			return

		case res, ok := <-keepAliceCh:
			if res != nil {
				// r.logger.Debugf(ctx, `keepalive loop: %v, %s`, ok, res.String())
			}
			if !ok {
				r.logger.Warnf(`keepalive exit, lease id: %d, retry register`, leaseID)
				// Re-register the service.
				for {
					if err := r.doRegisterLease(ctx, service); err != nil {
						retryDuration := time.Duration(rand.Intn(3*1000*1000))*time.Microsecond + time.Second
						r.logger.Errorf(
							`keepalive retry register failed, will retry in %s: %+v`,
							retryDuration, err,
						)
						time.Sleep(retryDuration)
						continue
					}
					break
				}
				return
			}
		}
	}
}

func (r *Registry) doRegisterLease(ctx context.Context, service g.Service) error {
	r.lease = etcd3.NewLease(r.client)

	ctx, cancel := context.WithTimeout(context.Background(), r.dialTimeout)
	defer cancel()

	grant, err := r.lease.Grant(ctx, int64(r.keepaliveTTL.Seconds()))
	if err != nil {
		return gone.ToErrorWithMsg(err, fmt.Sprintf(`etcd grant failed with keepalive ttl "%s"`, r.keepaliveTTL))
	}
	var (
		key   = g.GetServiceId(service)
		value = g.GetServerValue(service)
	)
	_, err = r.client.Put(ctx, key, value, etcd3.WithLease(grant.ID))
	if err != nil {
		return gone.ToErrorWithMsg(err, fmt.Sprintf(`etcd put failed with key "%s", value "%s", lease "%d"`, key, value, grant.ID))
	}
	r.logger.Debugf(
		`etcd put success with key "%s", value "%s", lease "%d"`,
		key, value, grant.ID,
	)
	keepAliceCh, err := r.client.KeepAlive(context.Background(), grant.ID)
	if err != nil {
		return err
	}
	go r.doKeepAlive(service, grant.ID, keepAliceCh)
	return nil
}

func (r *Registry) Deregister(instance g.Service) error {
	_, err := r.client.Delete(context.Background(), g.GetServiceId(instance))
	if r.lease != nil {
		_ = r.lease.Close()
	}
	return gone.ToErrorWithMsg(err, "Deregister delete key from etcd err")
}
