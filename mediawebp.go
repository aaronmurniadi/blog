package main

import (
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/chai2010/webp"
)

const mediaWebPConvertConcurrency = 4

// Require leading quote so https://other.example/media/x.jpg is not rewritten.
var mediaRasterURLRE = regexp.MustCompile(`(?i)(["'])(/media/[^"'>\s]+\.)(jpe?g|png|gif)`)

func rewriteMediaRasterToWebP(html string) string {
	return mediaRasterURLRE.ReplaceAllString(html, `${1}${2}webp`)
}

func mediaRasterToWebPExt(ext string) bool {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg", ".png", ".gif":
		return true
	default:
		return false
	}
}

func copyMediaTreeAsWebP(srcDir, dstDir string) error {
	var convertJobs []struct{ src, dst string }
	var plainCopies []struct{ src, dst string }

	err := filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dstDir, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		ext := filepath.Ext(path)
		if mediaRasterToWebPExt(ext) {
			target = strings.TrimSuffix(target, ext) + ".webp"
			convertJobs = append(convertJobs, struct{ src, dst string }{path, target})
			return nil
		}
		plainCopies = append(plainCopies, struct{ src, dst string }{path, target})
		return nil
	})
	if err != nil {
		return err
	}

	log.Printf("media: queued %d raster -> webp, %d binary copy", len(convertJobs), len(plainCopies))

	sem := make(chan struct{}, mediaWebPConvertConcurrency)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var convErrs []error
	for _, job := range convertJobs {
		job := job
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			rel, relErr := filepath.Rel(srcDir, job.src)
			if relErr != nil {
				rel = job.src
			}
			if err := convertImageFileToWebP(job.src, job.dst); err != nil {
				mu.Lock()
				convErrs = append(convErrs, err)
				mu.Unlock()
				log.Printf("media: webp FAIL %s: %v", filepath.ToSlash(rel), err)
				return
			}
			log.Printf("media: webp ok %s", filepath.ToSlash(rel))
		}()
	}
	wg.Wait()
	if len(convErrs) > 0 {
		return errors.Join(convErrs...)
	}

	for _, job := range plainCopies {
		rel, relErr := filepath.Rel(srcDir, job.src)
		if relErr != nil {
			rel = job.src
		}
		if err := copyFile(job.src, job.dst); err != nil {
			return fmt.Errorf("copy %s: %w", filepath.ToSlash(rel), err)
		}
		log.Printf("media: copy %s", filepath.ToSlash(rel))
	}
	return nil
}

func convertImageFileToWebP(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	img, _, err := image.Decode(in)
	if err != nil {
		return fmt.Errorf("decode %s: %w", src, err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if err := webp.Encode(out, img, &webp.Options{Quality: 90}); err != nil {
		return fmt.Errorf("encode %s: %w", dst, err)
	}
	return nil
}
