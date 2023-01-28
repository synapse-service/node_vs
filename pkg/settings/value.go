package settings

import (
	"time"

	"github.com/google/uuid"
)

type Value struct {
	GatewayAPI struct {
		Address string `yaml:"address"`
	} `yaml:"gateway_api"`
	Cameras []CameraInfo `yaml:"cameras"`
}

type CameraInfo struct {
	ID    uuid.UUID     `yaml:"id"`
	Name  string        `yaml:"name"`
	URL   string        `yaml:"url"`
	Timer time.Duration `yaml:"timer"`
}
