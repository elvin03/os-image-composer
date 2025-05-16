package debutils

import (
	"fmt"

	"github.com/intel-innersource/os.linux.tiberos.os-curation-tool/internal/provider"
)

// ParsePrimary parses the repodata/primary.xml.gz file from a given base URL.
func ParsePrimary(baseURL, gzHref string) ([]provider.PackageInfo, error) {

	// Download the debian repo .gz file with all components meta data
	PkgMetaFile := "/tmp/Packages.gz"
	zipFiles, err := Download(gzHref, PkgMetaFile)
	if err != nil {
		return nil, fmt.Errorf("failed to download repo file: %v", err)
	}

	// Decompress the .gz file and store the decompressed file in the same location
	if len(zipFiles) == 0 {
		return []provider.PackageInfo{}, fmt.Errorf("no files downloaded from repo URL: %s", gzHref)
	}
	files, err := Decompress(zipFiles[0])
	if err != nil {
		return []provider.PackageInfo{}, err
	}
	fmt.Printf("decompressed files: %v\n", files)

	//todo: parse the decompressed file

	return []provider.PackageInfo{}, nil
}
