package parser

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/Misora000/easyhtml"
	"github.com/Misora000/mangadl/logging"
	"github.com/Misora000/mangadl/types"
)

// LoveHeavenParser handles loveha.net.
type LoveHeavenParser struct {
}

// NewLoveHeavenParser new a LoveHeavenParser object.
func NewLoveHeavenParser() types.Parser {
	return &LoveHeavenParser{}
}

// ParseChapterList implemets types.Parser.
func (p *LoveHeavenParser) ParseChapterList(
	ctx context.Context, URL string) (*types.ChapterList, error) {

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	u, err := url.Parse(URL)
	if err != nil {
		return nil, types.ErrInvalidURL
	}

	chapList := types.NewChapterList(ctx, URL)
	body, err := chapList.GetPage()
	if err != nil {
		return nil, err
	}

	z := easyhtml.NewTokenizer(body)

	// Title
	_, eof := z.JumpToTag("h1")
	if eof {
		return nil, types.ErrInsuficientChapterInfo
	}
	if chapList.Name, eof = z.GetNextText(); eof {
		return nil, types.ErrInsuficientChapterInfo
	}

	// Chapter is
	// <a class="chapter" href='xxx.html' title="xxx Chapter xx">...</a>
	for {
		attr, eof := z.JumpToClass("a", "chapter")
		if eof {
			break
		}

		if href, exists := attr["href"]; exists {
			// Generate full url.
			chapURL := u.Scheme + "://" + u.Host + "/" + href

			chap := types.NewChapter(ctx, chapURL)
			chapList.Chapters = append(chapList.Chapters, chap)

			// Parse chapter name & number.
			if title, exists := attr["title"]; exists {
				seg := strings.Split(title, " Chapter ")
				if len(seg) == 2 {
					chap.Name = seg[0]
					chap.ChapNo, _ = strconv.Atoi(seg[1])
				}
			}
		}
	}

	return chapList, nil
}

// ParseChapter implemets types.Parser.
func (p *LoveHeavenParser) ParseChapter(
	ctx context.Context, URL string) (*types.Chapter, error) {

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	chap := types.NewChapter(ctx, URL)

	body, err := chap.GetPage()
	if err != nil {
		return nil, err
	}

	z := easyhtml.NewTokenizer(body)

	if err := p.parseChapterInfo(z, chap); err != nil {
		return nil, err
	}

	for {
		attr, eof := z.JumpToClass("img", "chapter-img")
		if eof {
			break
		}

		src, exists := attr["data-src"]
		if !exists {
			continue
		}

		if raw, err := base64.StdEncoding.DecodeString(src); err != nil {
			logging.Error(err.Error())
		} else {
			chap.PicsURL = append(chap.PicsURL,
				strings.ReplaceAll(string(raw), "\n", ""))
		}
	}

	return chap, nil
}

func (p *LoveHeavenParser) parseChapterInfo(
	z *easyhtml.Tokenizer, chap *types.Chapter) error {

	// Dummy.
	for i := 0; i < 2; i++ {
		if _, eof := z.JumpToTagAttr("span", "itemprop", "name"); eof {
			return types.ErrInsuficientChapterInfo
		}
	}

	// Name.
	if _, eof := z.JumpToTagAttr("span", "itemprop", "name"); eof {
		return types.ErrInsuficientChapterInfo
	}
	text, eof := z.GetNextText()
	if eof {
		return types.ErrInsuficientChapterInfo
	}
	chap.Name = text

	// Chapter No.
	if _, eof := z.JumpToTagAttr("span", "itemprop", "name"); eof {
		return types.ErrInsuficientChapterInfo
	}
	if text, eof = z.GetNextText(); eof {
		return types.ErrInsuficientChapterInfo
	}
	// Format: CHAPTER 26
	cNo, err := strconv.Atoi(text[len("CHAPTER "):])
	if err != nil {
		return fmt.Errorf("invalid chapter: %v", text)
	}
	chap.ChapNo = cNo

	return nil
}
