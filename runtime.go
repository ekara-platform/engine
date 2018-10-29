package engine

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/ekara-platform/model"
)

func getCurrentDirectoryURL(ctx *context) (*url.URL, error) {

	wd, err := os.Getwd()
	if err != nil {
		ctx.logger.Println("Error getting the working directory")
		return nil, err
	}
	ctx.logger.Printf("wd: %s\n", wd)

	absWd, err := filepath.Abs(wd)
	if err != nil {
		ctx.logger.Println("Error getting the absolute working directory")
		return nil, err
	}
	ctx.logger.Printf("absWd: %s\n", absWd)

	if strings.HasPrefix(absWd, "/") {
		absWd = "file://" + filepath.ToSlash(absWd)
	} else {
		absWd = "file:///" + filepath.ToSlash(absWd)
	}

	if err != nil {
		ctx.logger.Println("Error3 parsing")
		return nil, err
	}

	wdUrl, err := url.Parse(absWd)
	if err != nil {
		ctx.logger.Println("Error4 parsing")
		return nil, err
	}

	ctx.logger.Printf("wdUrl: %s\n", wdUrl)
	wdUrl, err = model.NormalizeUrl(wdUrl)
	if err != nil {
		ctx.logger.Println("Error5 parsing")
		return nil, err
	}
	return wdUrl, nil

}
