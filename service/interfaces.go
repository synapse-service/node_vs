package service

import "github.com/synapse-service/node-vs/pkg/settings"

type Settings interface {
	Get() settings.Value
	Update(settings.Value) error
	OnUpdate(callback func())
}
