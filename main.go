package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	facebookCountURL = "https://graph.facebook.com/?id="
	hatenaCountURL   = "http://api.b.st-hatena.com/entry.count?url="
)

type share struct {
	CommentCount int `json:"comment_count"`
	ShareCount   int `json:"share_count"`
}

type ogObject struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	UpdatedTime string `json:"updated_time"`
}

type facebook struct {
	Share    share    `json:"share"`
	OgObject ogObject `json:"og_object"`
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Fprintf(os.Stdout, "What manual page do you want?")
		os.Exit(1)
	}

	url, err := transeformURL(args[1])

	if err != nil {
		os.Exit(man(args[1:]))
	}

	os.Exit(webinfo(url))
}

// webinfo prints page's title and description
func webinfo(url string) int {
	printMeta(url)
	printSNSCount(url)
	return 0
}

// printMeta prints site meta data
func printMeta(url string) error {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "goquery Error:%s\n", err.Error())
		return err
	}
	fmt.Fprint(os.Stdout, "Page Infomation: [ページメタ情報]\n")
	doc.Find("head").Each(func(i int, s *goquery.Selection) {
		// title
		fmt.Fprint(os.Stdout, fmt.Sprintf("Title :\t\t\t\t%s\n", s.Find("title").Text()))
		s.Find("meta").Each(func(i int, s *goquery.Selection) {
			// description
			if name, _ := s.Attr("name"); name == "description" {
				description, _ := s.Attr("content")
				fmt.Fprint(os.Stdout, fmt.Sprintf("Description :\t\t\t%s\n", description))
			}
			// og:title
			if prop, _ := s.Attr("property"); prop == "og:title" {
				title, _ := s.Attr("content")
				fmt.Fprint(os.Stdout, fmt.Sprintf("og:title :\t\t\t%s\n", title))
			}

			// og:description
			if prop, _ := s.Attr("property"); prop == "og:description" {
				description, _ := s.Attr("content")
				fmt.Fprint(os.Stdout, fmt.Sprintf("og:description :\t\t%s\n", description))
			}

			// og:image
			if prop, _ := s.Attr("property"); prop == "og:image" {
				image, _ := s.Attr("content")
				fmt.Fprint(os.Stdout, fmt.Sprintf("og:image :\t\t\t%s\n", image))
			}
		})
	})

	return err
}

// printSNSCount prints SNS count
func printSNSCount(url string) {
	printHatena(url)
	printFacebook(url)
}

// printFacebook prints comment and share count
func printFacebook(url string) error {
	c, s, err := getFacebookCount(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "getFacebookCount Error:%s\n", err.Error())
		return err
	}

	fmt.Fprint(os.Stdout, fmt.Sprintf("FacebookCount(comment):\t\t%d\n", c))
	fmt.Fprint(os.Stdout, fmt.Sprintf("FacebookCount(share):\t\t%d\n", s))

	return nil
}

// printHatena prints hatenabookmark count.
func printHatena(url string) error {
	c, err := getHatenaCount(url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "getHatenaCount Error:%s\n", err.Error())
		return err
	}
	fmt.Fprint(os.Stdout, fmt.Sprintf("HatenaBookMark:\t\t\t%d\n", c))

	return nil
}

// getHatenaCount gets hatenabookmark count.
func getHatenaCount(url string) (int, error) {
	u := fmt.Sprintf("%s%s", hatenaCountURL, url)
	var client = &http.Client{Timeout: 10 * time.Second}
	r, err := client.Get(u)
	if err != nil {
		return 0, err
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return 0, err
	}

	c, err := strconv.Atoi(string(b))

	return c, err
}

// getFacebookCount returns facebook comment and share counts.
func getFacebookCount(url string) (int, int, error) {
	u := fmt.Sprintf("%s%s", facebookCountURL, url)
	var client = &http.Client{Timeout: 10 * time.Second}
	r, err := client.Get(u)
	if err != nil {
		return 0, 0, err
	}
	defer r.Body.Close()

	f := facebook{}
	json.NewDecoder(r.Body).Decode(&f)
	return f.Share.CommentCount, f.Share.ShareCount, err
}

// man runs man command.
func man(options []string) int {
	cmd := exec.Command("man", options...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "man Error:%s\n", err.Error())
		return 1
	}
	return 0
}

// transeformURL returns url string and error
func transeformURL(str string) (string, error) {
	u, err := url.Parse(str)
	if u.Host != "" && err == nil {
		return str, nil
	}
	return "", errors.New("str is not url format")
}
