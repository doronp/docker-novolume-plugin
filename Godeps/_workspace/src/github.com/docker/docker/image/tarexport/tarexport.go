package tarexport

import (
	"github.com/docker/docker/tag"
	"github.com/runcom/docker-novolume-plugin/Godeps/_workspace/src/github.com/docker/docker/image"
	"github.com/runcom/docker-novolume-plugin/Godeps/_workspace/src/github.com/docker/docker/layer"
)

const (
	manifestFileName           = "manifest.json"
	legacyLayerFileName        = "layer.tar"
	legacyConfigFileName       = "json"
	legacyVersionFileName      = "VERSION"
	legacyRepositoriesFileName = "repositories"
)

type manifestItem struct {
	Config   string
	RepoTags []string
	Layers   []string
}

type tarexporter struct {
	is image.Store
	ls layer.Store
	ts tag.Store
}

// NewTarExporter returns new ImageExporter for tar packages
func NewTarExporter(is image.Store, ls layer.Store, ts tag.Store) image.Exporter {
	return &tarexporter{
		is: is,
		ls: ls,
		ts: ts,
	}
}
