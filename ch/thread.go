package ch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Thread struct {
	Posts []Post
}

type Post struct {
	No          int    `json:"no"`
	Now         string `json:"now"`
	Name        string `json:"name"`
	Com         string `json:"com,omitempty"`
	Filename    string `json:"filename,omitempty"`
	Ext         string `json:"ext,omitempty"`
	W           int    `json:"w,omitempty"`
	H           int    `json:"h,omitempty"`
	TnW         int    `json:"tn_w,omitempty"`
	TnH         int    `json:"tn_h,omitempty"`
	Tim         int64  `json:"tim,omitempty"`
	Time        int    `json:"time"`
	Md5         string `json:"md5,omitempty"`
	Fsize       int    `json:"fsize,omitempty"`
	Resto       int    `json:"resto"`
	Bumplimit   int    `json:"bumplimit,omitempty"`
	Imagelimit  int    `json:"imagelimit,omitempty"`
	SemanticURL string `json:"semantic_url,omitempty"`
	Replies     int    `json:"replies,omitempty"`
	Images      int    `json:"images,omitempty"`
	UniqueIps   int    `json:"unique_ips,omitempty"`
	TailSize    int    `json:"tail_size,omitempty"`
}

func GetMedia(board, id string) ([]string, error) {
	resp, err := http.Get(fmt.Sprintf("https://a.4cdn.org/%s/thread/%s.json", board, id))
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var thread Thread
	if err = json.Unmarshal(bytes, &thread); err != nil {
		return nil, err
	}

	media := make([]string, 0)
	for _, post := range thread.Posts {
		url := fmt.Sprintf("https://i.4cdn.org/%s/%d%s", board, post.Tim, post.Ext)
		media = append(media, url)
	}

	return media, nil
}
