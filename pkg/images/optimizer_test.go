package images_test

import (
	"testing"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/memnix/memnix-rest/pkg/images"
)

func jpgToTest() string {
	return "test/1.jpg"
}

func pngToTest() string {
	return "test/2.png"
}

func TestConvertToWebp(t *testing.T) {
	images.StartVips()

	type args struct {
		buffer string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ConvertJpgToWebp",
			args: args{
				buffer: jpgToTest(),
			},
			wantErr: false,
		},
		{
			name: "ConvertPngToWebp",
			args: args{
				buffer: pngToTest(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got1, err := images.ConvertToWebp(tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToWebp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Check the metadata forrmats
			if !tt.wantErr {
				if got1.Format.FileExt() != ".webp" {
					t.Errorf("ConvertToWebp() got1 = %v, want %v", got1.Format, ".webp")
				}
			}
		})
	}
}

func TestStoreImage(t *testing.T) {
	images.StartVips()

	type args struct {
		name   string
		buffer string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "StoreJpgImage",
			args: args{
				name:   "test/1_out.jpg",
				buffer: jpgToTest(),
			},
			wantErr: false,
		},
		{
			name: "StorePngImage",
			args: args{
				name:   "test/2_out.jpg",
				buffer: pngToTest(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newImage, err := vips.NewImageFromFile(tt.args.buffer)
			if err != nil {
				t.Errorf("StoreImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			defer newImage.Close()

			newJpg, _, _ := newImage.ExportJpeg(vips.NewJpegExportParams())

			if err := images.StoreImage(tt.args.name, newJpg); (err != nil) != tt.wantErr {
				t.Errorf("StoreImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConvertToWebpWithQuality(t *testing.T) {
	images.StartVips()

	type args struct {
		name    string
		quality int
	}

	tests := []struct {
		name    string
		args    args
		want    []byte
		want1   vips.ImageMetadata
		wantErr bool
	}{
		{
			name: "ConvertToWebpWithQuality",
			args: args{
				name:    "test/1.jpg",
				quality: 80,
			},
			wantErr: false,
		},
		{
			name: "ConvertToWebpWithQuality",
			args: args{
				name:    "test/2.png",
				quality: 80,
			},
			wantErr: false,
		},
		{
			name: "ConvertToWebpWithWrongQuality",
			args: args{
				name:    "test/1.jpg",
				quality: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got1, err := images.ConvertToWebpWithQuality(tt.args.name, tt.args.quality)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToWebp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Check the metadata formats
			if !tt.wantErr {
				if got1.Format.FileExt() != ".webp" {
					t.Errorf("ConvertToWebp() got1 = %v, want %v", got1.Format.FileExt(), ".webp")
				}
			}
		})
	}
}
