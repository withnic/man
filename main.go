package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
)

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	var out []byte
	var err error
	if len(args) < 2 {
		fmt.Printf("What manual page do you want?")
		return 1
	}

	u, err := url.Parse(args[1])

	if u.Host != "" {
		out, err = exec.Command("whois", u.Host).Output()
		if err != nil {
			log.Fatal(err)
			return 1
		}
		fmt.Fprint(os.Stdout, string(out))
		return 0
	}

	out, err = exec.Command("man", args[1:]...).Output()
	if err != nil {
		log.Fatal(err)
		return 1
	}

	fmt.Fprint(os.Stdout, string(out))
	return 0
}
