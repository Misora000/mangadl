# Manga Downloader

## Usage
```sh
> go run main.go -url WEBSITE_URL [attributes]

Attributes:
    -l          : indicate url is the chapter list page
    -limit NUM  : at most download NUM chapters
    -max NUM    : the maximum chapter no. to download
    -min NUM    : the minimum chapter no. to download
    -info       : show supported sites
    -preview    : show pics url only, not download
```

## Example
```sh
# Download the whole books
> go run main.go -url https://loveheaven.net/manga-one-piece.html -l

# Donwload chapter 981~986
> go run main.go -url https://loveheaven.net/manga-one-piece.html -l -max 986 -min 981

# Download the top 30 chapters
> go run main.go -url https://loveheaven.net/manga-one-piece.html -l -limit 30

# Donwload the specific chapter
> go run main.go -url https://loveheaven.net/read-one-piece-chapter-983.html
```

## Supported sites
- https://loveheaven.net/