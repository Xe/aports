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

	shouldWork(ctx, nil, wd, "rsync", "-avz","/home/xena/packages/", "greedo.xeserv.us:public_html/pkg/alpine/edge/")
}

var wg sync.WaitGroup

func BuildAll() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	coreDir := filepath.Join(wd, "./core")
	goDir := filepath.Join(wd, "./go")

	coreFin, err := os.Open(coreDir)
	qod.ANE(err)
	defer coreFin.Close()

	coreFis, err := coreFin.Readdir(999)
	qod.ANE(err)

	for _, fi := range coreFis {
		go buildPackage(ctx, fi, coreDir)
	}

	goFin, err := os.Open(goDir)
	qod.ANE(err)
	defer goFin.Close()

	goFis, err := goFin.Readdir(999)
	qod.ANE(err)

	for _, fi := range goFis {
		go buildPackage(ctx, fi, goDir)
	}

	time.Sleep(time.Second)

	wg.Wait()
}


func buildPackage(ctx context.Context, fi os.FileInfo, parentDir string) {
	packDir := filepath.Join(parentDir, filepath.Base(fi.Name()))

	wg.Add(1)
	defer wg.Done()

	shouldWork(ctx, nil, packDir, "abuild", "checksum")
	shouldWork(ctx, nil, packDir, "abuild", "-r")
}
