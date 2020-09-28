package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/url"
	"os"
	"time"

	"github.com/shopspring/decimal"

	"github.com/Misora000/mangadl/engine"
	"github.com/Misora000/mangadl/types"
)

var (
	reqURL   = flag.String("url", "", "request URL")
	listMode = flag.Bool("l", false, "is chapter list")
	preview  = flag.Bool("preview", false, "preview only, not download")
	info     = flag.Bool("info", false, "show information")
	limit    = flag.Int("limit", -1, "how many chapters to download. -1 means no limit")
	max      = flag.Int("max", -1, "")
	min      = flag.Int("min", -1, "")

	randSleepDuration = time.Duration(rand.Intn(15)) * time.Second

	config Config
)

// Config defines the user config.
type Config struct {
	DownloadRoot string `json:"download_root"`
}

func loadConfig() error {
	f, err := os.Open("config.json")
	if err != nil {
		return err
	}

	text, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	fmt.Printf("config:\n%v\n-------------------------------\n", string(text))

	return json.Unmarshal(text, &config)
}

func main() {

	flag.Parse()

	if err := loadConfig(); err != nil {
		panic(err)
	}

	engine.Initialize()
	defer engine.Finalize()

	if *info {
		engine.PrintSupportedSites()
		return
	}

	if _, err := url.Parse(*reqURL); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := engine.GetParser(*reqURL)
	chapters := []*types.Chapter{}

	// Collect request chapters.
	if *listMode {
		chapList, err := p.ParseChapterList(ctx, *reqURL)
		if err != nil {
			panic(err)
		}
		fmt.Println("Name:", chapList.Name)
		fmt.Println("Chapter Num:", len(chapList.Chapters))
		chapters = chapList.Chapters
	} else {
		chapters = append(chapters, types.NewChapter(ctx, *reqURL))
	}

	dlCount := 0
	defer func() {
		fmt.Println("total downloaded chapters:", dlCount)
	}()

	// Download pics of all chapters.
	for i, v := range chapters {
		if *limit > 0 && dlCount >= *limit {
			fmt.Println("terminate by chapters limit:", *limit)
			break
		}

		chap, err := p.ParseChapter(ctx, v.PageURL)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if skipThisChapter(chap.ChapNo) {
			fmt.Println("skip chapter:", chap.ChapNo)
			continue
		}

		fmt.Println("----", chap.Name, fmt.Sprintf("#%v", chap.ChapNo), "----")

		if *preview {
			for _, p := range chap.PicsURL {
				fmt.Println(p)
			}
			continue
		}

		totalDL := chap.DownloadPics(config.DownloadRoot)
		fmt.Println("totoal downloaded:", totalDL, "pics")

		dlCount++

		// Random sleeping might be like a human.
		if i < len(chapters)-1 {
			rand.Seed(time.Now().UnixNano())
			time.Sleep(randSleepDuration)
		}
	}

	return
}

func skipThisChapter(chapNo string) bool {
	no, err := decimal.NewFromString(chapNo)
	if err != nil {
		fmt.Println("none decimal chapter no:", chapNo)
		return false
	}

	if *max > -1 && no.GreaterThan(decimal.NewFromInt(int64(*max))) {
		return true
	}

	if *min > -1 && no.LessThan(decimal.NewFromInt(int64(*min))) {
		return true
	}

	return false
}
