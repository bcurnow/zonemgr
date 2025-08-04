/*
Copyright © 2025 Brian Curnow

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package parse

import (
	"fmt"
	"os"

	"github.com/bcurnow/zonemgr/schema"
	"gopkg.in/yaml.v3"
)

func ToZones(inputFile string) (map[string]*schema.Zone, error) {
	logger.Debug("Opening input file", "inputFile", inputFile)
	inputBytes, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to open input %s: %w", inputFile, err)
	}

	logger.Debug("Unmarshaling YAML", "inputFile", inputFile)
	zones, err := unmarshal(inputBytes)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal input bytes: %w", err)
	}

	// Normalize the zones
	logger.Debug("Normalizing the zones", "inputFile", inputFile, "zoneCount", len(zones))
	zones, err = normalize(zones)
	if err != nil {
		return nil, fmt.Errorf("Failed to normalize zones: %w", err)
	}
	return zones, nil
}

func unmarshal(inputBytes []byte) (map[string]*schema.Zone, error) {
	var zones map[string]*schema.Zone
	err := yaml.Unmarshal(inputBytes, &zones)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse input YAML: %w", err)
	}
	if len(zones) == 0 {
		return nil, fmt.Errorf("No zones found in input file")
	}
	return zones, nil
}
