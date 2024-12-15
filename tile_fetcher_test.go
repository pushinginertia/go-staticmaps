package sm_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/pushinginertia/go-staticmaps"
)

func TestTileFetcher_storeCache(t *testing.T) {
	path, err := ioutil.TempDir("", "test_storeCache")
	if err != nil {
		panic(err)
	}

	defer func() {
		fmt.Printf("Removing: %v\n", path)
		os.RemoveAll(path)
	}()

	fileName := filepath.Join(path, "tile")

	provider := sm.NewTileProviderOpenStreetMaps()
	cache := sm.NewTileCache(path, 0700)

	fetcher := sm.NewTileFetcher(provider, cache)

	file, err := os.OpenFile(
		fileName,
		os.O_RDWR|os.O_CREATE|os.O_EXCL,
		0600,
	)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = fetcher.StoreCache(fileName, []byte{0, 1, 2, 3, 4, 5, 6, 7})
	if err != nil {
		panic(err)
	}
	//print err
	//assert err != nil
}
