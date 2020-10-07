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

	"github.com/Misora000/mangadl/engine"
	"github.com/Misora000/mangadl/logging"
	"github.com/Misora000/mangadl/types"

	"github.com/shopspring/decimal"
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
	LogLevel     int    `json:"log_level"`
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

	logging.Initialize(config.LogLevel)
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
		logging.Log("Name: %v", chapList.Name)
		logging.Log("Chapters: %v", len(chapList.Chapters))
		chapters = chapList.Chapters
	} else {
		chapters = append(chapters, types.NewChapter(ctx, *reqURL))
	}

	dlCount := 0
	defer func() {
		logging.Log("total downloaded chapters: %v", dlCount)
	}()

	// Download pics of all chapters.
	for i, v := range chapters {
		if *limit > 0 && dlCount >= *limit {
			logging.Log("terminate by chapters limit: %v", *limit)
			break
		}

		chap, err := p.ParseChapter(ctx, v.PageURL)
		if err != nil {
			logging.Error(err.Error())
			continue
		}

		if skipThisChapter(chap.ChapNo) {
			logging.Log("skip chapter: %v", chap.ChapNo)
			continue
		}

		logging.Log("---- %v %v ----", chap.Name, fmt.Sprintf("#%v", chap.ChapNo))

		if *preview {
			for _, p := range chap.PicsURL {
				logging.Log(p)
			}
			continue
		}

		totalDL := chap.DownloadPics(config.DownloadRoot)
		logging.Log("totoal downloaded: %v pics", totalDL)

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
		logging.Error("none decimal chapter no: %v", chapNo)
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
