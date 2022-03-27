package features_test

import (
	_ "embed"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/hedhyw/gherkingen/pkg/v1/bdd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const pathImageResize = "/api/image.jpg"

func TestPictureResizing(t *testing.T) {
	t.Parallel()

	f := bdd.NewFeature(t, "Picture resizing")

	f.Scenario("Duyên wants to receive a resized picture", func(_ *testing.T, f *bdd.Feature) {
		type testCase struct {
			Width  int `field:"<width>"`
			Height int `field:"<height>"`
		}

		testCases := map[string]testCase{
			"256_128": {256, 128},
			"128_128": {128, 128},
			"1_1":     {1, 1},
		}

		f.TestCases(testCases, func(t *testing.T, f *bdd.Feature, tc testCase) {
			t.Parallel()

			th := newTestHelper(t)
			q := make(url.Values)
			var resp *http.Response

			f.Given("Duyên selects the size '<width>x<height>'", func() {
				q.Set("size", fmt.Sprintf("%dx%d", tc.Width, tc.Height))
			})
			f.And("a link to existen picture of a bigger size of type 'image/jpeg'", func() {
				photoURL := newPhotoServer(t)
				q.Set("url", photoURL)
			})
			f.When("Duyên calls an endpoint", func() {
				resp = th.CallAPIGet(t, pathImageResize, q)
				require.Equal(t, http.StatusOK, resp.StatusCode)
			})
			f.Then("she receives an image of a content type 'image/jpeg'", func() {
				assert.Equal(t, resp.Header.Get("Content-Type"), "image/jpeg")
			})
			f.And("the size of the image is '<width>x<height>'", func() {
				assertImageSize(t, resp.Body, tc.Width, tc.Height)
			})
		})
	})

	f.Scenario("Duyên provides an invalid size", func(_ *testing.T, f *bdd.Feature) {
		type testCase struct {
			Size string `field:"<size>"`
		}

		testCases := map[string]testCase{
			"0x0":   {"0x0"},
			"0x1":   {"0x1"},
			"ax10":  {"ax10"},
			"10,10": {"10,10"},
			"-1x10": {"-1x10"},
			"xxx":   {"xxx"},
			"1x1x1": {"1x1x1"},
		}

		f.TestCases(testCases, func(t *testing.T, f *bdd.Feature, tc testCase) {
			t.Parallel()

			th := newTestHelper(t)
			q := make(url.Values)
			var resp *http.Response

			f.Given("Duyên selects the size <size>", func() {
				q.Set("size", tc.Size)
			})
			f.And("a link to an existen picture of a bigger size", func() {
				photoURL := newPhotoServer(t)
				q.Set("url", photoURL)
			})
			f.When("Duyên calls an endpoint", func() {
				resp = th.CallAPIGet(t, pathImageResize, q)
			})
			f.Then("she receives an error", func() {
				require.Equal(t, http.StatusBadRequest, resp.StatusCode)
			})
		})
	})

	f.Scenario("Duyên provides an invalid link", func(t *testing.T, f *bdd.Feature) {
		t.Parallel()

		th := newTestHelper(t)
		q := make(url.Values)
		var resp *http.Response

		ts := httptest.NewServer(http.NotFoundHandler())
		t.Cleanup(ts.Close)

		f.Given("Duyên selects the size 256x256", func() {
			q.Set("size", "256x256")
		})
		f.And("a link to an unexistent picture", func() {
			q.Set("url", ts.URL)
		})
		f.When("Duyên calls an endpoint", func() {
			resp = th.CallAPIGet(t, pathImageResize, q)
		})
		f.Then("she receives an error", func() {
			require.Equal(t, http.StatusNotFound, resp.StatusCode)
		})
	})

	f.Scenario("Duyên provides not jpeg image", func(t *testing.T, f *bdd.Feature) {
		t.Parallel()

		th := newTestHelper(t)
		q := make(url.Values)
		var resp *http.Response

		f.Given("Duyên selects the size 256x256", func() {
			q.Set("size", "256x256")
		})
		f.And("a link to an existent plain/text file", func() {
			plainTextURL := newPlainTextServer(t)
			q.Set("url", plainTextURL)
		})
		f.When("Duyên calls an endpoint", func() {
			resp = th.CallAPIGet(t, pathImageResize, q)
		})
		f.Then("she receives an error", func() {
			require.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
		})
	})
}

func assertImageSize(tb testing.TB, imageBody io.Reader, width int, height int) {
	tb.Helper()

	im, err := jpeg.Decode(imageBody)
	if assert.NoError(tb, err) {
		bounds := im.Bounds()
		assert.Equal(tb, bounds.Max.X, width)
		assert.Equal(tb, bounds.Max.Y, height)
	}
}

// Photo by Joyston Judah from Pexels
// Photo: https://www.pexels.com/photo/white-and-black-mountain-wallpaper-933054
// License: https://www.pexels.com/license/
//
//go:embed assets/pexels-photo-933054.jpg
var photoBody []byte

func newPhotoServer(tb testing.TB) (url string) {
	tb.Helper()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")

		_, err := w.Write(photoBody)
		assert.NoError(tb, err)
	}))

	tb.Cleanup(ts.Close)

	return ts.URL
}

func newPlainTextServer(tb testing.TB) (url string) {
	tb.Helper()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "plain/text")
		_, err := w.Write([]byte("hello world"))
		assert.NoError(tb, err)
	}))

	tb.Cleanup(ts.Close)

	return ts.URL
}
