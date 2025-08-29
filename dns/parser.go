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
package dns

import (
	"fmt"
	"os"

	"github.com/bcurnow/zonemgr/models"
	"github.com/hashicorp/go-hclog"
	"gopkg.in/yaml.v3"
)

type ZoneParser interface {
	Parse(inputFile string, globalConfig *models.Config) (map[string]*models.Zone, error)
}

type yamlZoneParser struct {
	ZoneParser
	normalizer Normalizer
}

func YamlZoneParser(normalizer Normalizer) ZoneParser {
	return &yamlZoneParser{normalizer: normalizer}
}

func (p *yamlZoneParser) Parse(inputFile string, globalConfig *models.Config) (map[string]*models.Zone, error) {
	hclog.L().Debug("Opening input file", "inputFile", inputFile)
	inputBytes, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open input %s: %w", inputFile, err)
	}

	hclog.L().Debug("Unmarshaling YAML", "inputFile", inputFile)
	zones, err := p.unmarshal(inputBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal from %s: %w", inputFile, err)
	}

	for name, zone := range zones {
		// It is possible for the zone itself to be nil, this happens if a file is parsed which only contains the name of the zone and no other info
		if zone == nil {
			return nil, fmt.Errorf("invalid input file %s, no zone information for zone %s", inputFile, name)
		}
		// Make sure we always have a complete config
		if nil == zone.Config {
			zone.Config = globalConfig
		}

	}

	// Normalize the zones
	if err = p.normalizer.Normalize(zones); err != nil {
		return nil, fmt.Errorf("failed to normalize zones: %w", err)
	}
	return zones, nil
}

func (p *yamlZoneParser) unmarshal(inputBytes []byte) (map[string]*models.Zone, error) {
	var zones map[string]*models.Zone
	err := yaml.Unmarshal(inputBytes, &zones)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input YAML: %w", err)
	}
	if len(zones) == 0 {
		return nil, fmt.Errorf("no zones found in input file")
	}
	return zones, nil
}
