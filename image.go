package main

import (
	"embed"
	"fmt"

	"fyne.io/fyne/v2"
)

//go:embed img/*
var imageFS embed.FS

type ImageFile struct {
	boardName string
	filePath  string
}

func (i ImageFile) LoadResource() fyne.Resource {
	content, err := imageFS.ReadFile(i.filePath)
	if err != nil {
		panic(fmt.Errorf("Cannot open file %s: %w", i.filePath, err))
	}
	return &Image{i.boardName, content}
}

type Image struct {
	name    string
	content []byte
}

func (i *Image) Name() string {
	return i.name
}

func (i *Image) Content() []byte {
	return i.content
}
