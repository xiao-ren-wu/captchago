package capembed

import (
	"testing"
)

func TestFileResources_RandSliderResource(t *testing.T) {
	sliderResource, err := DefaultFileResource.RandSliderResource()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", sliderResource)
}
