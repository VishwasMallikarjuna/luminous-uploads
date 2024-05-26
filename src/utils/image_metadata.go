package utils

import (
	"fmt"
	"image"
	_ "image/jpeg" // Import for JPEG decoding
	_ "image/png"  // Import for PNG decoding
	"mime/multipart"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/VishwasMallikarjuna/luminous-uploads/model"
)

func ExtractMetadata(file multipart.File) (*model.ImageMetadata, error) {
	img, format, err := image.DecodeConfig(file)
	if err != nil {
		return nil, fmt.Errorf("format error: %v", err)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("seek error: %v", err)
	}

	x, err := exif.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode EXIF metadata: %v", err)
	}

	metadata := &model.ImageMetadata{
		Width:  img.Width,
		Height: img.Height,
		Format: format,
	}

	if x != nil {
		cameraModel, err := x.Get(exif.Model)
		if err == nil {
			metadata.CameraModel, _ = cameraModel.StringVal()
		}

		lat, long, err := x.LatLong()
		if err == nil {
			metadata.Location = fmt.Sprintf("Lat: %f, Long: %f", lat, long)
		}
	}

	return metadata, nil
}
