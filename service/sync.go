package service

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	gateway_api "github.com/synapse-service/gateway/transport/grpc"
	"github.com/synapse-service/gateway/transport/grpc/client"

	"github.com/synapse-service/node-vs/pkg/settings"
)

const ProcessNameSync = "sync"

func (s *Service) syncProcess(ctx context.Context) error {
	for address := range gatewayAddress(ctx, s.settings.Get().GatewayAPI.Address, s.cfg.GatewayAPI.FallbackAddress) {
		if err := s.syncTry(ctx, address); err != nil {
			s.log.Error().Err(err).Msg("sync try")
		}
	}
	return nil
}

func (s *Service) syncTry(ctx context.Context, address string) error {
	gatewayClient, err := client.NewClient(
		client.WithClientAddress(address),
		client.WithTransportCredentials(s.cfg.GatewayAPI.Certificate),
	)
	if err != nil {
		return errors.Wrap(err, "new client")
	}

	defer gatewayClient.Stop(ctx) //nolint:errcheck

	if err := gatewayClient.Start(ctx); err != nil {
		return errors.Wrap(err, "start gateway client")
	}

	req := gateway_api.SyncRequest{
		Id: s.cfg.ID.String(),
	}
	stream, err := gatewayClient.Sync(ctx, &req)
	if err != nil {
		return errors.Wrap(err, "sync")
	}

	return s.syncLoop(ctx, stream)
}

func (s *Service) syncLoop(ctx context.Context, stream gateway_api.GatewayAPI_SyncClient) error {
loop:
	for {
		select {
		case <-ctx.Done():
			if err := stream.CloseSend(); err != nil {
				return errors.Wrap(err, "close send")
			}
			break loop

		default:

			res, err := stream.Recv()
			if err == io.EOF {
				s.log.Debug().Err(err).Msg("eof")
				break loop
			}
			if err != nil {
				return errors.Wrap(err, "receive")
			}

			var value settings.Value
			value.Cameras = make([]settings.CameraInfo, len(res.Cameras))
			for i, info := range res.Cameras {
				value.Cameras[i].ID = uuid.MustParse(info.Id)
				value.Cameras[i].Name = info.Name
				value.Cameras[i].URL = info.Url
			}
			if err := s.settings.Update(value); err != nil {
				return errors.Wrap(err, "update settings")
			}

		}
	}

	return nil
}

func gatewayAddress(ctx context.Context, main, fallback string) <-chan string {
	addressCh := make(chan string)
	go func() {
		distribution := []struct {
			address string
			timeout time.Duration
		}{
			{main, time.Second},
			{main, 2 * time.Second},
			{main, 3 * time.Second},
			{main, 5 * time.Second},
			{main, 10 * time.Second},
			{main, 15 * time.Second},
			{main, 30 * time.Second},
			{main, time.Minute},
			{main, 2 * time.Minute},
			{main, 3 * time.Minute},
			{main, 5 * time.Minute},
			{main, 10 * time.Minute},
			{fallback, time.Minute},
			{fallback, time.Minute},
			{fallback, time.Minute},
			{fallback, time.Minute},
			{fallback, time.Minute},
			{fallback, 2 * time.Minute},
			{fallback, 3 * time.Minute},
			{fallback, 5 * time.Minute},
			{fallback, 10 * time.Minute},
		}
		for _, item := range distribution {
			item := item
			time.Sleep(item.timeout)
			addressCh <- item.address
		}
		last := distribution[len(distribution)-1]
		for {
			time.Sleep(last.timeout)
			addressCh <- last.address
		}
	}()
	return addressCh
}
