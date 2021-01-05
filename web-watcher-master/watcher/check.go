package watcher

import (
	"bytes"
	"compress/zlib"
	"log"

	"github.com/pmezard/go-difflib/difflib"
	"golang.org/x/net/html"

	"github.com/shellbear/web-watcher/models"
)

func (w *Watcher) checkChanges(task *models.Task, body []byte) (bool, error) {
	if task.Body == nil {
		return true, nil
	}

	r, err := zlib.NewReader(bytes.NewBuffer(task.Body))
	if err != nil {
		return false, err
	}

	defer r.Close()

	previousHTML, err := html.Parse(r)
	if err != nil {
		return false, err
	}

	newHTML, err := html.Parse(bytes.NewBuffer(body))
	if err != nil {
		return false, err
	}

	matcher := difflib.NewMatcher(extractTags(previousHTML), extractTags(newHTML))
	ratio := matcher.Ratio()

	if ratio < w.ChangeRatio {
		log.Printf("Changes detected for: %s. Changes ratio: %f < %f\n", task.URL, ratio, w.ChangeRatio)
		return true, nil
	}

	log.Println("No changed detected for:", task.URL)
	return false, nil
}
