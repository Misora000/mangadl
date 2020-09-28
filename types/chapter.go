package types

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Misora000/mangadl/engine/downloader"
	"github.com/Misora000/mangadl/logging"
)

var (
	timeoutGetPage = 1 * time.Minute
	timeoutGetPic  = 30 * time.Second
)

// Chapter is a chapter with its pictures.
type Chapter struct {
	downloader.Downloader
	ctx      context.Context
	Name     string
	ChapNo   string
	PageURL  string
	PicsURL  []string
	ReferURL string
}

// // ChapterJob common methods interface.
// type ChapterJob interface {
// 	GetPage() (io.Reader, error)
// 	DownloadPics(dst string) int
// }

// NewChapter new a chapter object.
func NewChapter(ctx context.Context, url string) *Chapter {
	return &Chapter{
		ctx:     ctx,
		PageURL: url,
	}
}

// GetPage gets raw html of the given url.
func (c *Chapter) GetPage() (io.Reader, error) {
	if c.ctx.Err() != nil {
		return nil, c.ctx.Err()
	}

	ctx, cancel := context.WithTimeout(c.ctx, timeoutGetPage)
	defer cancel()

	ret, reqURL, _, err := c.Get(ctx, c.PageURL)
	// Because the request would been redireted to another url, we need to store
	// the final url here for https referer header.
	c.ReferURL = reqURL.String()
	return ret, err
}

// DownloadPics downloads all pictures of this chapter.
func (c *Chapter) DownloadPics(dst string) (succ int) {
	parent := downloader.WithReferer(c.ctx, c.ReferURL)

	for i, url := range c.PicsURL {
		if c.ctx.Err() != nil {
			return
		}

		if len(url) == 0 {
			continue
		}

		// Create folder.
		home := fmt.Sprintf("%v/%v", dst, c.Name)
		if err := os.Mkdir(home, 0644); err != nil && !os.IsExist(err) {
			logging.Error(err.Error())
			break
		}
		chap := fmt.Sprintf("%v/%v", home, c.ChapNo)
		if err := os.Mkdir(chap, 0644); err != nil && !os.IsExist(err) {
			logging.Error(err.Error())
			break
		}
		path := fmt.Sprintf("%v/%03v", chap, i+1)

		// TODO: skip existing files.
		// I don't know the exact file extension now. How to check exist?

		ctx, cancel := context.WithTimeout(parent, timeoutGetPic)
		if err := c.SaveAs(ctx, url, path); err != nil {
			logging.Log(err.Error())
		} else {
			succ++
		}

		cancel()
	}

	return
}

// ChapterList is the set of chapters.
type ChapterList struct {
	downloader.Downloader
	ctx      context.Context
	Name     string
	PageURL  string
	Chapters []*Chapter
}

// NewChapterList new a ChapterList object.
func NewChapterList(ctx context.Context, url string) *ChapterList {
	return &ChapterList{
		ctx:     ctx,
		PageURL: url,
	}
}

// GetPage gets raw html of the given url.
func (l *ChapterList) GetPage() (io.Reader, error) {
	if l.ctx.Err() != nil {
		return nil, l.ctx.Err()
	}

	ctx, cancel := context.WithTimeout(l.ctx, timeoutGetPage)
	defer cancel()

	ret, _, _, err := l.Get(ctx, l.PageURL)
	return ret, err
}
