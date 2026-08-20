package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	videoURL := flag.String("url", "", "direct HTTPS video URL")
	outputPath := flag.String("out", "video.mp4", "output file path")
	flag.Parse()

	if err := run(*videoURL, *outputPath); err != nil {
		fmt.Fprintln(os.Stderr, "download failed:", err)
		os.Exit(1)
	}
}

func run(videoURL, outputPath string) error {
	if err := validateURL(videoURL); err != nil {
		return err
	}

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(videoURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	written, err := io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("downloaded %d bytes to %s\n", written, outputPath)
	return nil
}

func validateURL(raw string) error {
	if raw == "" {
		return errors.New("missing --url")
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return err
	}
	if parsed.Scheme != "https" {
		return errors.New("only HTTPS URLs are supported")
	}
	if parsed.Host == "" {
		return errors.New("missing URL host")
	}
	return nil
}
