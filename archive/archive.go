package archive

import (
	"fmt"
	"os"
	"path"

	"github.com/mholt/archiver"
	clog "github.com/openbiox/bioctl/log"
	"github.com/openbiox/bioctl/stringo"
)

var log = clog.Logger

// UnarchiveLog use github.com/mholt/archiver to uncompress and unarchive files
func UnarchiveLog(source, destination string) (err error) {
	sourceTmp := path.Join(path.Dir(source), "_"+path.Base(source))
	err = os.Symlink(source, sourceTmp)
	archive, err := archiver.ByExtension(sourceTmp)
	if err == nil {
		archiveType := fmt.Sprintf("%T", archive)
		for _, v := range []string{"*archiver.Xz", "*archiver.Gz", "*archiver.Zstd", "archiver.Snappy", "archiver.Lz4",
			"*archiver.Bz2", "*archiver.Brotli"} {
			if v == archiveType {
				destination = path.Join(destination, stringo.StrReplaceAll(path.Base(source), ".[a-zA-Z]*$", ""))
				log.Infof("Uncompressing %s => %s", source, destination)
				err = archiver.DecompressFile(sourceTmp, destination)
			}
		}
		for _, v := range []string{"*archiver.Rar", "*archiver.Tar", "*archiver.TarBrotli", "*archiver.TarBz2",
			"*archiver.TarGz", "*archiver.TarLz4", "*archiver.TarXz", "*archiver.TarSz", "*archiver.Zip", "*archiver.TarZstd"} {
			if v == archiveType {
				log.Infof("Uncompressing %s => %s", source, destination)
				err = archiver.Unarchive(source, destination)
			}
		}
	}
	os.Remove(sourceTmp)
	return err
}
