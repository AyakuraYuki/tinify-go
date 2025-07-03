package tinify_test

import (
	"os"
	"testing"

	"github.com/AyakuraYuki/tinify-go/tinify"
)

const testKey = "H2yRky1qwH8xTRzPx6lK0g1n7LThh364"

var client *tinify.Client

func TestMain(m *testing.M) {
	client = tinify.NewClient(testKey)
	_ = os.MkdirAll("./testoutput", os.ModePerm)
	m.Run()
}

func TestCompressFromFile(t *testing.T) {
	source, err := client.FromFile("testdata/sunflower.jpg")
	if err != nil {
		t.Fatal(err)
	}

	err = client.ToFile(source, "./testoutput/compressed-from-file.jpg")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCompressFromBuffer(t *testing.T) {
	buffer, err := os.ReadFile("testdata/sunflower.jpg")
	if err != nil {
		t.Fatal(err)
	}

	source, err := client.FromBuffer(buffer)
	if err != nil {
		t.Fatal(err)
	}

	err = client.ToFile(source, "./testoutput/compressed-from-buffer.jpg")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCompressFromURL(t *testing.T) {
	source, err := client.FromURL("https://pixabay.com/get/gcbc64f6fc4d1aee5f677f8f641516dc522756e6e9b1831414388527555a427d8c52aad142cd7906666be45f71b9601c7a124a9384788c8049c6185f797edf47d345cc74a76c9ba9bb92618f09fd3476e_1920.jpg?attachment=")
	if err != nil {
		t.Fatal(err)
	}

	err = client.ToFile(source, "./testoutput/compressed-from-url.jpg")
	if err != nil {
		t.Fatal(err)
	}
}

func TestResize(t *testing.T) {
	source, err := client.FromFile("testdata/sunflower.jpg")
	if err != nil {
		t.Fatal(err)
	}

	err = client.Resize(source, &tinify.ResizeOption{
		Method: tinify.ResizeMethodFit,
		Width:  128,
		Height: 128,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = client.ToFile(source, "./testoutput/resize-from-file.jpg")
	if err != nil {
		t.Fatal(err)
	}
}
