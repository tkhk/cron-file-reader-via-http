package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

var filePath string

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func fatalWithStackTrace(err error) {
	st, ok := err.(stackTracer)
	if !ok {
		log.Fatalf("err does not implement stackTracer: raw err: %+v", err)
	}

	log.Fatal(
		fmt.Sprintf("%+v", err),
		fmt.Sprintf("%+v", st.StackTrace()),
	)
}

func init() {
	fileName := flag.String("f", "", "filepath")
	flag.Parse()
	if len(*fileName) == 0 {
		fatalWithStackTrace(errors.New("f option is empty"))
	}

	currentPath, err := os.Getwd()
	if err != nil {
		fatalWithStackTrace(errors.Wrap(err, "could not os.Getwd"))
	}

	p := filepath.Join(currentPath, *fileName)
	if _, err := os.Stat(p); err != nil {
		fatalWithStackTrace(errors.Wrap(err, fmt.Sprintf("file(%s) does not exists", p)))
	}

	filePath = p
}

func read(filePath string) (string, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return string(b), nil
}

func cron(ctx context.Context, filePath string, fileBody *string) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Done!")
			return
		case <-ticker.C:
			v, err := read(filePath)
			if err != nil {
				return
			}
			*fileBody = v
		}
	}
}

func main() {
	var body string

	ctx, cancel := context.WithCancel(context.Background())
	go cron(ctx, filePath, &body)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	})

	// cancel 使いたかっただけ
	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		cancel()
	})
	http.ListenAndServe(":8080", nil)
}
