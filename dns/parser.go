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

	"github.com/bcurnow/zonemgr/models"
	"github.com/bcurnow/zonemgr/utils"
)

type ZoneParser interface {
	Parse(inputFile string, globalConfig *models.Config) (map[string]*models.Zone, error)
}

type yamlZoneParser struct {
	ZoneParser
	normalizer Normalizer
	reader     *utils.ZoneYamlFile
}

func YamlZoneParser(normalizer Normalizer) ZoneParser {
	return &yamlZoneParser{normalizer: normalizer, reader: &utils.ZoneYamlFile{}}
}

func (p *yamlZoneParser) Parse(inputFile string, globalConfig *models.Config) (map[string]*models.Zone, error) {
	zones, err := p.reader.Read(inputFile)
	if err != nil {
		return nil, err
	}

	if len(zones) == 0 {
		return nil, fmt.Errorf("no zones found in input file")
	}

	for name, zone := range zones {
		// It is possible for the zone itself to be nil, this happens if a file is parsed which only contains the name of the zone and no other info
		if zone == nil {
			return nil, fmt.Errorf("invalid input file %s, no zone information for zone %s", inputFile, name)
		}
	}

	// Normalize the zones
	if err = p.normalizer.Normalize(zones, globalConfig); err != nil {
		return nil, fmt.Errorf("failed to normalize zones: %w", err)
	}
	return zones, nil
}
