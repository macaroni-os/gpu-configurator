/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package specs

func NewEglExternalPlatformFiles(dir string) *EglExternalPlatformFiles {
	return &EglExternalPlatformFiles{
		Path:  dir,
		Files: make(map[string]*ICDJson, 0),
	}
}

func (e *EglExternalPlatformFiles) AddFile(file string, icdjson *ICDJson) {
	e.Files[file] = icdjson
}
