package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	os.Exit(run(os.Args))
}

// webinfo prints page's title and description
func webinfo(url string) int {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Printf("goquery Error:%s\n", err.Error())
		return 1
	}
	fmt.Fprint(os.Stdout, "Page Infomation: [ページメタ情報]\n")
	doc.Find("head").Each(func(i int, s *goquery.Selection) {
		fmt.Fprint(os.Stdout, fmt.Sprintf("Title :\t\t\t\t%s\n", s.Find("title").Text()))
		s.Find("meta").Each(func(i int, s *goquery.Selection) {
			if name, _ := s.Attr("name"); name == "description" {
				description, _ := s.Attr("content")
				fmt.Fprint(os.Stdout, fmt.Sprintf("Description :\t\t\t%s\n", description))
			}
		})
	})
	return 0
}

// man runs man command.
func man(options []string) int {
	cmd := exec.Command("man", options...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("man Error:%s\n", err.Error())
		return 1
	}
	return 0
}

// transeformURL returns url string
func transeformURL(str string) (string, error) {
	u, err := url.Parse(str)
	if u.Host != "" {
		return str, err
	}
	return "url parse error", errors.New("str is not url format")
}

func run(args []string) int {
	if len(args) < 2 {
		fmt.Printf("What manual page do you want?")
		return 1
	}

	url, err := transeformURL(args[1])

	if err != nil {
		return man(args[1:])
	}
	return webinfo(url)
}
