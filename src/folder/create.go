package folder

import (
	"fmt"
	"os"

	"github.com/aryadiahmad4689/infra-lib/src/entity" // Sesuaikan dengan path yang benar
)

// Folder struct untuk mengelola direktori infrastruktur
type Folder struct {
	BasePath string
}

// NewFolder adalah constructor untuk Folder struct
func NewFolder(basePath string) *Folder {
	return &Folder{BasePath: basePath}
}

// Create membuat seluruh struktur direktori di dalam base path yang ditentukan
func (f *Folder) Create() error {
	directories := []entity.Directory{
		{Path: f.BasePath + "/app/app1"},
		{Path: f.BasePath + "/app/app2"},
		{Path: f.BasePath + "/terraform/dev"},
		{Path: f.BasePath + "/terraform/production"},
		{Path: f.BasePath + "/terraform/staging"},
		{Path: f.BasePath + "/tools/app-tools"},
		{Path: f.BasePath + "/tools/argocd-installer-prod"},
		{Path: f.BasePath + "/tools/argocd-installer-staging"},
		{Path: f.BasePath + "/tools/list-tools/argocd"},
		{Path: f.BasePath + "/tools/list-tools/prometheus"},
		{Path: f.BasePath + "/tools/list-tools/elk"},
		{Path: f.BasePath + "/tools/list-tools/grafana"},
		{Path: f.BasePath + "/tools/list-tools/redis"},
		{Path: f.BasePath + "/tools/list-tools/mongo"},
		{Path: f.BasePath + "/tools/list-tools/kafka"},
		{Path: f.BasePath + "/tools/list-tools/outline"},
		{Path: f.BasePath + "/tools/list-tools/vault"},
	}

	for _, dir := range directories {
		err := os.MkdirAll(dir.Path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating directory %s: %v", dir.Path, err)
		}
		fmt.Printf("Created directory: %s\n", dir.Path)
	}

	return nil
}
