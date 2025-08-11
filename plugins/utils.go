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

package plugins

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bcurnow/zonemgr/schema"
)


func RenderResourceWithoutValue(rr *schema.ResourceRecord) string {
	var record strings.Builder

	record.WriteString(fmt.Sprintf(ResourceRecordNameFormatString, rr.Name))
	record.WriteString(fmt.Sprintf(ResourceRecordTypeFormatString, rr.Type))
	if rr.Class != "" {
		record.WriteString(" ")
		record.WriteString(rr.Class)
		record.WriteString(" ")
	}

	if rr.TTL != nil {
		record.WriteString(" ")
		record.WriteString(strconv.Itoa(int(*rr.TTL)))
		record.WriteString(" ")
	}

	return record.String()
}

func RenderSingleValueResource(rr *schema.ResourceRecord) string {
	var record strings.Builder
	record.WriteString(RenderResourceWithoutValue(rr))
	record.WriteString(rr.Value)
	if rr.Comment != "" {
		record.WriteString(" ; ")
		record.WriteString(rr.Comment)
	}

	return record.String()
}

func RenderMultivalueResource(rr *schema.ResourceRecord) string {
	var record strings.Builder
	record.WriteString(RenderResourceWithoutValue(rr))
	record.WriteString("(\n")
	indentFormatString := "%" + strconv.Itoa(record.Len()-2) + "s"
	for _, value := range rr.Values {
		record.WriteString(fmt.Sprintf(indentFormatString, ""))                         // This will ensure that all the values are indented
		record.WriteString(fmt.Sprintf(ResourceRecordMultivalueIndentFormatString, "")) // This will add an indent inside the parens
		record.WriteString(fmt.Sprintf(ResourceRecordNameFormatString, value.Value))
		if value.Comment != "" {
			record.WriteString(" ; ")
			record.WriteString(value.Comment)
		}
		record.WriteString("\n")
	}
	record.WriteString(fmt.Sprintf(indentFormatString, ""))
	record.WriteString(")")

	return record.String()
}
