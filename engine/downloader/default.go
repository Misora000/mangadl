package downloader

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strings"

	"github.com/Misora000/mangadl/logging"
)

type ctxKey string

// CtxKey enum.
var (
	CtxKeyReferer ctxKey = "referer"
)

// Downloader provides basic function for download page and file.
type Downloader struct {
	cookie []string
}

// WithReferer setups Referer header if the schem is https.
func WithReferer(parent context.Context, URL string) context.Context {
	if strings.HasPrefix(URL, "https://") {
		u := strings.Split(URL[len("https://"):], "/")
		r := "https://" + u[0] + "/"
		return context.WithValue(parent, CtxKeyReferer, r)
	}
	return parent
}

// Get implement HTTP GET.
// Return:
//   1. body          io.Reader
//   2. request URL   *url.URL
//   3. content-type  string
//   4. error         error
func (d *Downloader) Get(
	ctx context.Context, URL string) (io.Reader, *url.URL, string, error) {

	logging.Debug("[GET] %v", URL)

	if ctx.Err() != nil {
		return nil, nil, "", ctx.Err()
	}

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, nil, "", err
	}
	req.WithContext(ctx)

	// Set referer for https.
	if referer, ok := ctx.Value(CtxKeyReferer).(string); ok && len(referer) > 0 {
		req.Header.Add("referer", referer)
	}

	// for i, v := range d.header {
	// 	req.Header.Add(i, v)
	// }

	for _, v := range d.cookie {
		req.Header.Add("cookie", v)
	}

	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		return nil, nil, "", err
	}

	for i, v := range rsp.Header {
		if i == "Set-Cookie" {
			d.cookie = append(d.cookie, v...)
		}
	}

	buf, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, nil, "", err
	}
	rsp.Body.Close()

	return bytes.NewReader(buf), rsp.Request.URL, rsp.Header.Get("Content-Type"), nil
}

// SaveAs stores file from url with the file extension by parsing content-type.
func (d *Downloader) SaveAs(ctx context.Context, URL string, dst string) error {

	r, _, contentType, err := d.Get(ctx, URL)
	if err != nil {
		return err
	}

	// Decide the file extension.
	extension := ""
	exts, err := mime.ExtensionsByType(contentType)
	if err != nil {
		return err
	}
	if len(exts) > 0 {
		extension = exts[0]
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	// write file
	return ioutil.WriteFile(dst+extension, data, 0644)
}
