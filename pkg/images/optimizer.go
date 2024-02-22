package images

import (
	"os"
	"sync"

	"github.com/davidbyttow/govips/v2/vips"
)

var vipsOnce sync.Once //nolint:gochecknoglobals //Singleton

func StartVips() {
	vipsOnce.Do(func() {
		vips.Startup(nil)
	})
}

func StopVips() {
	vips.Shutdown()
}

func ConvertToWebp(name string) ([]byte, vips.ImageMetadata, error) {
	image1, err := vips.NewImageFromFile(name)
	if err != nil {
		return nil, vips.ImageMetadata{}, err
	}
	defer image1.Close()

	newBuffer, metadata, err := image1.ExportWebp(vips.NewWebpExportParams())
	if err != nil {
		return nil, vips.ImageMetadata{}, err
	}

	return newBuffer, *metadata, nil
}

func StoreImage(name string, buffer []byte) error {
	image1, err := vips.NewImageFromBuffer(buffer)
	if err != nil {
		return err
	}

	defer image1.Close()

	err = os.WriteFile(name, buffer, 0o600) // Fix: Change permissions to 0o600
	if err != nil {
		return err
	}

	return nil
}
