package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fonero-project/fnod/chaincfg"
	"github.com/fonero-project/fnod/fnoutil"
	"github.com/fonero-project/fnotime/fnotimed/backend"
	"github.com/fonero-project/fnotime/fnotimed/backend/filesystem"
)

var (
	defaultHomeDir = fnoutil.AppDataDir("fnotimed", false)

	file        = flag.String("file", "", "journal of modifications if used (will be written despite -fix)")
	fix         = flag.Bool("fix", false, "Try to correct correctable failures")
	fnodataHost = flag.String("host", "", "fnodata block explorer")
	printHashes = flag.Bool("printhashes", false, "Print all hashes")
	fsRoot      = flag.String("source", "", "Source directory")
	testnet     = flag.Bool("testnet", false, "Use testnet port")
	verbose     = flag.Bool("v", false, "Print more information during run")
)

func _main() error {
	flag.Parse()

	root := *fsRoot
	if root == "" {
		root = filepath.Join(defaultHomeDir, "data")
		if *testnet {
			root = filepath.Join(root, chaincfg.TestNetParams.Name)
		} else {
			root = filepath.Join(root, chaincfg.MainNetParams.Name)
		}
	}

	if *fnodataHost == "" {
		if *testnet {
			*fnodataHost = "https://testnet.fonero.org/api/tx/"
		} else {
			*fnodataHost = "https://explorer.fonero.org/api/tx/"
		}
	} else {
		if !strings.HasSuffix(*fnodataHost, "/") {
			*fnodataHost += "/"
		}
	}

	fmt.Printf("=== Root: %v\n", root)

	fs, err := filesystem.NewDump(root)
	if err != nil {
		return err
	}
	defer fs.Close()

	return fs.Fsck(&backend.FsckOptions{
		Verbose:     *verbose,
		PrintHashes: *printHashes,
		Fix:         *fix,
		URL:         *fnodataHost,
		File:        *file,
	})
}

func main() {
	err := _main()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
