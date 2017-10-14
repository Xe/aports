// +build mage

package main

import (
	"context"
	"path/filepath"
	"os"
	"sync"
	"time"

	"github.com/jtolds/qod"
)

// Push pushes locally built packages to the remote package host
func Push() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shouldWork(ctx, nil, wd, "rsync", "-avz","/home/xena/packages/core/x86_64/", "greedo.xeserv.us:public_html/pkg/alpine/edge/core/x86_64")
}

func BuildAll() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	coreDir := filepath.Join(wd, "./core")

	var wg sync.WaitGroup

	fin, err := os.Open(coreDir)
	qod.ANE(err)

	fis, err := fin.Readdir(999)
	qod.ANE(err)

	for _, fi := range fis {
		go func(fi os.FileInfo) {
			packDir := filepath.Join(coreDir, filepath.Base(fi.Name()))

			wg.Add(1)
			defer wg.Done()

			shouldWork(ctx, nil, packDir, "abuild", "checksum")
			shouldWork(ctx, nil, packDir, "abuild", "-r")
		} (fi)
	}

	time.Sleep(time.Second)

	wg.Wait()
}
