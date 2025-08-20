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

package utils

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofrs/flock"
	"github.com/hashicorp/go-hclog"
	"gopkg.in/yaml.v3"
)

const serialChangeIndexFileExtension = ".serial"

var (
	initialChangeIndex    uint32                = 1
	serialNumberGenerator SerialNumberGenerator = &TimeBasedSerialGenerator{}
)

type SerialIndex struct {
	BaseSerialNumber *uint32 `yaml:"base_serial_number"`
	ChangeIndex      *uint32 `yaml:"change_index"`
}

func (si *SerialIndex) toSerial() string {
	return strconv.Itoa(int(*si.BaseSerialNumber)) + strconv.Itoa(int(*si.ChangeIndex))
}

type SerialIndexManager interface {
	GetNext(zoneName string) (string, error)
}

type FileSerialIndexManager struct {
}

func (m *FileSerialIndexManager) GetNext(zoneName string) (string, error) {
	if err := m.createSerialChangeIndexDirectory(); err != nil {
		return "", err
	}

	serialFile := m.fileName(zoneName)

	// Make sure the file exists
	if m.exists(serialFile) {
		hclog.L().Trace("Serial change index file exists, processing", "file", serialFile)
		serial, err := m.incrementAndUpdate(serialFile)
		if err != nil {
			return "", err
		}
		return serial, nil
	}

	hclog.L().Trace("Serial change index file does not exist, creating new one", "file", serialFile)
	serialIndex, err := m.initFile(serialFile)
	if err != nil {
		return "", err
	}
	hclog.L().Trace("Returning next serial number", "serialNumber", serialIndex.toSerial())
	return serialIndex.toSerial(), nil
}

func (m *FileSerialIndexManager) parseFile(serialFile string) (*SerialIndex, error) {
	hclog.L().Trace("Opening serial change index file", "file", serialFile)
	inputBytes, err := os.ReadFile(serialFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open serial change index file %s: %w", serialFile, err)
	}

	hclog.L().Trace("Unmarshaling YAML", "file", serialFile)
	serialIndex, err := m.unmarshal(inputBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal input bytes: %w", err)
	}

	return serialIndex, nil
}

func (m *FileSerialIndexManager) initFile(serialFile string) (*SerialIndex, error) {
	hclog.L().Debug("Creating new serial file", "file", serialFile)
	//Lock the file so no other process modifies it while we're updating
	fileLock, err := m.getLock(serialFile)
	if err != nil {
		return nil, err
	}
	defer fileLock.Unlock()

	file, err := os.Create(serialFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	baseSerial, err := serialNumberGenerator.GenerateBaseSerial()
	if err != nil {
		return nil, err
	}

	// Create a new SerialIndex structure. We hardcode 1 as this is a new file
	serialIndex := &SerialIndex{BaseSerialNumber: baseSerial, ChangeIndex: &initialChangeIndex}
	if err := m.marshal(file, serialIndex); err != nil {
		return nil, err
	}

	return serialIndex, nil
}

func (m *FileSerialIndexManager) incrementAndUpdate(serialFile string) (string, error) {
	//Lock the file so no other process modifies it while we're updating
	fileLock, err := m.getLock(serialFile)
	if err != nil {
		return "", err
	}
	defer fileLock.Unlock()

	serialIndex, err := m.parseFile(serialFile)
	if err != nil {
		return "", err
	}

	//Generate a new base serial number and compare to the base in the file, if they aren't the same, it's a different day
	//and we should start back at initialChangeIndex
	newBaseSerial, err := serialNumberGenerator.GenerateBaseSerial()
	if err != nil {
		return "", err
	}

	hclog.L().Trace("Comparing base serial numbers", "current", *serialIndex.BaseSerialNumber, "new", *newBaseSerial)
	if *serialIndex.BaseSerialNumber != *newBaseSerial {
		serialIndex.BaseSerialNumber = newBaseSerial
		serialIndex.ChangeIndex = &initialChangeIndex
	} else {
		*serialIndex.ChangeIndex++
	}

	hclog.L().Trace("Writing updated serial change index file", "file", serialFile, "baseSerialNumber", *serialIndex.BaseSerialNumber, "changeIndex", *serialIndex.ChangeIndex)
	// Write the updated values back to the file
	file, err := os.OpenFile(serialFile, os.O_WRONLY, os.FileMode(0650))
	if err != nil {
		return "", err
	}
	defer file.Close()

	if err := m.marshal(file, serialIndex); err != nil {
		return "", err
	}
	hclog.L().Trace("Returning next serial number", "serialNumber", serialIndex.toSerial())
	return serialIndex.toSerial(), nil

}

func (m *FileSerialIndexManager) unmarshal(inputBytes []byte) (*SerialIndex, error) {
	var serialIndex SerialIndex
	err := yaml.Unmarshal(inputBytes, &serialIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input YAML: %w", err)
	}
	return &serialIndex, nil
}

func (m *FileSerialIndexManager) marshal(file *os.File, serialIndex *SerialIndex) error {
	yamlBytes, err := yaml.Marshal(serialIndex)
	if err != nil {
		return err
	}

	if _, err := file.Write(yamlBytes); err != nil {
		return err
	}
	return nil
}

func (m *FileSerialIndexManager) fileName(zoneName string) string {
	return filepath.Join(SerialChangeIndexDirectory.Value, zoneName+serialChangeIndexFileExtension)
}

func (m *FileSerialIndexManager) exists(serialFile string) bool {
	// now check if the file exists
	_, err := os.Stat(serialFile)
	return !errors.Is(err, os.ErrNotExist)
}

func (m *FileSerialIndexManager) createSerialChangeIndexDirectory() error {
	if _, err := os.Stat(SerialChangeIndexDirectory.Value); os.IsNotExist(err) {
		if err := os.MkdirAll(SerialChangeIndexDirectory.Value, os.FileMode(0750)); err != nil {
			return err
		}
	}
	return nil
}

// Will try and get the lock for 10 seconds, returns an error if it can't get the lock
func (m *FileSerialIndexManager) getLock(serialFile string) (*flock.Flock, error) {
	// Create a file lock, this doesn't lock the file...yet
	fileLock := flock.New(serialFile)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	hclog.L().Trace("Attempting to lock serial change index file, will try for 10 seconds", "file", serialFile)
	locked, err := fileLock.TryLockContext(ctx, 500*time.Millisecond)
	if err != nil {
		return nil, err
	}

	if locked {
		hclog.L().Trace("Locked serial change index file", "file", serialFile)
		return fileLock, nil
	}

	// This should never happen
	return nil, fmt.Errorf("unexpected error, expected file to be locked")
}
