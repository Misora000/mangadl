package types

import "context"

// Parser defines the common methods of each site-parser.
type Parser interface {
	ParseChapterList(ctx context.Context, url string) (*ChapterList, error)
	ParseChapter(ctx context.Context, url string) (*Chapter, error)
}
