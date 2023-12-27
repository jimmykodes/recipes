package dist

import (
	"io"
	"os"
	"path/filepath"

	"github.com/jimmykodes/gommand"
	"github.com/jimmykodes/gommand/flags"
)

const (
	distDirFlag = "dist"
)

func Flags() *flags.FlagSet {
	fs := flags.NewFlagSet()

	fs.String(distDirFlag, "./dist", "dir location to write dist files to")

	return fs
}

func New(ctx *gommand.Context) (*Dist, error) {
	return &Dist{
		location: ctx.Flags().String(distDirFlag),
	}, nil
}

type Dist struct {
	location string
}

func (d Dist) Clear() error {
	if err := os.RemoveAll(d.location); err != nil {
		return err
	}
	if err := os.Mkdir(d.location, 0755); err != nil {
		return err
	}
	return nil
}

func (d Dist) Dir() string {
	return d.location
}

func (d Dist) Join(path string) string {
	return filepath.Join(d.location, path)
}

func (d Dist) Create(path string) (io.WriteCloser, error) {
	return os.Create(d.Join(path))
}
