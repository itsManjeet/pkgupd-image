package utils

import (
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func Extractfile(file, outdir string) error {
	return exec.Command("bsdtar", "-xf", file, "-C", outdir).Run()
}

func Copyfile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
func Preparefile(file, outpath string) error {
	switch filepath.Ext(file) {
	case "tar", "xz", "gz", "zstd", "zip", "deb", "rlx":
		return Extractfile(file, outpath)
	default:
		return Copyfile(file, outpath+"/"+file)
	}
}

func PrepareSources(urls []string, srcdir, wrkdir string) error {
	for _, url := range urls {
		filename := path.Base(url)
		if err := Preparefile(srcdir+"/"+filename, wrkdir); err != nil {
			return err
		}
	}

	return nil
}
