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
	"github.com/go-playground/validator/v10"
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
	_               YamlFile[map[string]*models.Zone] = &ZoneYamlFile{}
	_               YamlFile[*models.SerialIndex]     = &SerialIndexYamlFile{}
	unmarshal                                         = yaml.Unmarshal
	marshal                                           = yaml.Marshal
	openFile                                          = os.OpenFile
	marshalFileMode                                   = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	validate                                          = validator.New()
)

func (yr *ZoneYamlFile) Read(path string) (map[string]*models.Zone, error) {
	zones, err := unmarshalYaml[map[string]*models.Zone](path)
	if err != nil {
		return nil, err
	}

	// Validate the zones
	for zoneName, zone := range zones {
		if zone == nil {
			continue // Skip nil zones, they will be handled by the parser
		}
		if err := validate.Struct(zone); err != nil {
			return nil, fmt.Errorf("validation failed for zone '%s': %w", zoneName, err)
		}
	}

	return zones, nil
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
	logger().Debug("opening file", "path", path)
	inputBytes, err := readFile(path)
	if err != nil {
		return nilT, fmt.Errorf("failed to open '%s': %w", path, err)
	}

	logger().Debug("unmarshaling YAML", "path", path)
	var yaml T
	if err := unmarshal(inputBytes, &yaml); err != nil {
		return nilT, fmt.Errorf("failed to parse input YAML: %w", err)
	}

	return yaml, nil
}

func marshalYaml[T any](path string, content T) error {
	file, err := openFile(path, marshalFileMode, 0600)
	if err != nil {
		return fmt.Errorf("failed to open '%s': %w", path, err)
	}
	defer file.Close()

	yamlBytes, err := marshal(content)
	if err != nil {
		return fmt.Errorf("failed to marshal '%v': %w", content, err)
	}

	if _, err := file.Write(yamlBytes); err != nil {
		return fmt.Errorf("failed to write to '%s': %w", path, err)
	}

	return nil
}
