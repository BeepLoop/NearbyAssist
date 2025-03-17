package fs

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"nearbyassist/internal/service/core"
	"testing"
)

func generateMockImage() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	red := color.RGBA{255, 0, 0, 255}

	for x := 0; x < 200; x++ {
		for y := 0; y < 200; y++ {
			img.Set(x, y, red)
		}
	}

	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}

func TestSaveFile(t *testing.T) {
	input := generateMockImage()
	testOutDir := "../../../test/back_id/"
	directory := map[Category]string{
		ID_BACK: testOutDir,
	}
	hash := core.NewSha256()
	disk := NewDiskStorage(directory, hash)

	location, err := disk.SaveFile(File{
		Data:     input,
		Category: ID_BACK,
	})

	if err != nil {
		t.Errorf("SaveFile() failed, expected %v, got %v", nil, err)
	}

	t.Logf("File saved at %s", location)
}
