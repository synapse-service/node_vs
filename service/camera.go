package service

import (
	"context"

	"github.com/synapse-service/node-vs/pkg/process"
	"github.com/synapse-service/node-vs/pkg/settings"
)

func cameraProcess(info settings.CameraInfo) process.Process {
	return func(ctx context.Context) error {
		// if err := camera.Process(ctx, info); err != nil {
		// 	return errors.Wrap(err, "camera process")
		// }
		return nil
	}
}

// Start processes with new cameras and stop with absent cameras
func (s *Service) updateCameraProcesses() {
	// Cameras that we need to be processed
	cameras := s.settings.Get().Cameras
	target := make([]string, len(cameras))
	infos := make(map[string]settings.CameraInfo, len(cameras))
	for i, camera := range cameras {
		target[i] = camera.Name
		infos[camera.Name] = camera
	}
	// Currently processed cameras
	current := s.processes.Filter(ProcessClassCamera)
	// Calculate cameras ids to stop and start, leaving intersection alone
	toStop, _, toStart := Separate(current, target)
	s.processes.Stop(toStop...)
	for _, name := range toStart {
		s.processes.Start(name, ProcessClassCamera, cameraProcess(infos[name]))
	}
}
