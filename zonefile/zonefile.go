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
package zonefile

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/bcurnow/zonemgr/parse/schema"
)

type zoneTemplateData struct {
	Name string
	Zone *schema.Zone
}

func generateZone(name string, zone *schema.Zone, outputDir string, tmpl string) error {
	zoneFileTemplate, err := template.New("zonefile.tmpl").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("Failed to parse template: %w", err)
	}

	outputFile, err := os.Create(filepath.Join(outputDir, name))
	if err != nil {
		return fmt.Errorf("Failed to create output file for zone %s: %w", name, err)
	}
	defer outputFile.Close()

	fmt.Printf("Generating %s for zone %s\n", outputFile.Name(), name)
	err = zoneFileTemplate.Execute(outputFile, zoneTemplateData{Name: name, Zone: zone})
	if err != nil {
		return fmt.Errorf("Failed to execute template for zone %s: %w", name, err)
	}
	return nil
}
