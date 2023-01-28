package plugin

import (
	"bytes"
	"image"
	"image/jpeg"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/synapse-service/node-vs/proto/plugin"
)

func (p *Plugin) SendSnapshot(image image.Image) error {
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, image, &jpeg.Options{Quality: 60}); err != nil {
		return errors.Wrap(err, "encode image")
	}

	message := plugin.Output{
		Data: &plugin.Output_Image{
			Image: &plugin.Image{
				Bytes: buf.Bytes(),
			},
		},
	}
	b, err := proto.Marshal(&message)
	if err != nil {
		return errors.Wrap(err, "marshal image message")
	}

	p.stdinCh <- b

	return nil
}
