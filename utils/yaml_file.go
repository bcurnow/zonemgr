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
	"errors"
	"fmt"
	"os"

	"github.com/bcurnow/zonemgr/models"
	"github.com/hashicorp/go-hclog"
	"gopkg.in/yaml.v3"
)

type YamlFile[T any] interface {
	// Reads the supplied path as a YAML file and returns the unmarshalled contents
	Read(path string) (T, error)
	Write(path string, content T) error
}

type ZoneYamlFile struct {
}

type SerialIndexYamlFile struct {
}

var (
	_         YamlFile[map[string]*models.Zone] = &ZoneYamlFile{}
	_         YamlFile[*models.SerialIndex]     = &SerialIndexYamlFile{}
	unmarshal                                   = yaml.Unmarshal
	marshal                                     = yaml.Marshal
	open                                        = os.Open
)

func (yr *ZoneYamlFile) Read(path string) (map[string]*models.Zone, error) {
	return unmarshalYaml[map[string]*models.Zone](path)
}

// We don't need to write out a Zone back to a file (or do we?)
func (yr *ZoneYamlFile) Write(path string, content map[string]*models.Zone) error {
	return errors.New("not implemented")
}

func (sir *SerialIndexYamlFile) Read(path string) (*models.SerialIndex, error) {
	return unmarshalYaml[*models.SerialIndex](path)
}

func (sir *SerialIndexYamlFile) Write(path string, content *models.SerialIndex) error {
	return marshalYaml(path, content)
}

func unmarshalYaml[T any](path string) (T, error) {
	var nilT T
	hclog.L().Debug("Opening file", "path", path)
	inputBytes, err := readFile(path)
	if err != nil {
		return nilT, fmt.Errorf("failed to open '%s': %w", path, err)
	}

	hclog.L().Debug("Unmarshaling YAML", "path", path)
	var yaml T
	if err := unmarshal(inputBytes, &yaml); err != nil {
		return nilT, fmt.Errorf("failed to parse input YAML: %w", err)
	}

	return yaml, nil
}

func marshalYaml[T any](path string, content T) error {
	file, err := open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	yamlBytes, err := marshal(content)
	if err != nil {
		return err
	}

	if _, err := file.Write(yamlBytes); err != nil {
		return err
	}

	return nil
}
