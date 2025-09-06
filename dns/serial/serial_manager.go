/**
 * Copyright (C) 2025 Brian Curnow
 *
 * This file is part of zonemgr.
 *
 * zonemgr is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * zonemgr is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with zonemgr.  If not, see <https://www.gnu.org/licenses/>.
 */

package serial

import (
	"fmt"
	"path/filepath"

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/utils"
	"github.com/hashicorp/go-hclog"
)

const (
	changeIndexFileExtension        = "serial"
	initialChangeIndex       uint32 = 1
)

type SerialManager interface {
	Next(zoneName string) (string, error)
}

var (
	generator Generator                  = &TimeBasedGenerator{}
	fs        utils.FileSystemOperations = &utils.FileSystem{}
)

type fileSerialManager struct {
	changeIndexDirectory string
	indexFile            utils.YamlFile[*models.SerialIndex]
}

func FileSerialManager(changeIndexDirectory string) SerialManager {
	return &fileSerialManager{changeIndexDirectory: changeIndexDirectory, indexFile: &utils.SerialIndexYamlFile{}}
}

func (m *fileSerialManager) Next(zoneName string) (string, error) {
	if err := fs.MkdirAll(m.changeIndexDirectory, 0750); err != nil {
		return "", err
	}

	path := filepath.Join(m.changeIndexDirectory, fmt.Sprintf("%s.%s", zoneName, changeIndexFileExtension))

	// Make sure the file exists
	var serialIndex *models.SerialIndex
	if fs.Exists(path) {
		hclog.L().Trace("Serial change index file exists, processing", "file", path)
		si, err := m.incrementAndUpdate(path)
		if err != nil {
			return "", err
		}
		serialIndex = si
	} else {
		hclog.L().Trace("Serial change index file does not exist, creating new one", "file", path)
		si, err := m.initFile(path)
		if err != nil {
			return "", err
		}
		serialIndex = si
	}

	serialNumber, err := generator.FromSerialIndex(serialIndex)
	if err != nil {
		return "", err
	}
	hclog.L().Trace("Returning next serial number", "serialNumber", serialNumber)
	return serialNumber, nil
}

func (m *fileSerialManager) initFile(path string) (*models.SerialIndex, error) {
	hclog.L().Debug("Creating new serial file", "file", path)
	//Lock the file so no other process modifies it while we're updating
	fileLock, err := fs.Flock(path)
	if err != nil {
		return nil, err
	}
	defer fileLock.Unlock()

	base, err := generator.GenerateBase()
	if err != nil {
		return nil, err
	}

	// This is a bit strange, however, I don't want an initial value that can be changed
	// Since you can't get a pointer to a constant, this is the work around
	changeIndex := initialChangeIndex
	// Create a new SerialIndex structure. We hardcode 1 as this is a new file
	serialIndex := &models.SerialIndex{Base: base, ChangeIndex: &changeIndex}
	if err := m.indexFile.Write(path, serialIndex); err != nil {
		return nil, err
	}

	return serialIndex, nil
}

func (m *fileSerialManager) incrementAndUpdate(path string) (*models.SerialIndex, error) {
	//Lock the file so no other process modifies it while we're updating
	fileLock, err := fs.Flock(path)
	if err != nil {
		return nil, err
	}
	defer fileLock.Unlock()

	serialIndex, err := m.indexFile.Read(path)
	if err != nil {
		return nil, err
	}

	//Generate a new base serial number and compare to the base in the file, if they aren't the same, it's a different day
	//and we should start back at initialChangeIndex
	newBase, err := generator.GenerateBase()
	if err != nil {
		return nil, err
	}

	hclog.L().Trace("Comparing base serial numbers", "current", *serialIndex.Base, "new", *newBase)
	if *serialIndex.Base != *newBase {
		serialIndex.Base = newBase
		// Again with the constant/pointer workaround
		changeIndex := initialChangeIndex
		serialIndex.ChangeIndex = &changeIndex
	} else {
		*serialIndex.ChangeIndex++
	}

	hclog.L().Trace("Writing updated serial change index file", "file", path, "baseSerialNumber", *serialIndex.Base, "changeIndex", *serialIndex.ChangeIndex)
	// Write the updated values back to the file
	err = m.indexFile.Write(path, serialIndex)
	if err != nil {
		return nil, err
	}

	return serialIndex, nil
}
