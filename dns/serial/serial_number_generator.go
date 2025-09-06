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
	"strconv"
	"time"

	"github.com/bcurnow/zonemgr/models"
)

var (
	sprintf   = fmt.Sprintf
	parseUint = strconv.ParseUint
)

type Generator interface {
	GenerateBase() (*uint32, error)
	FromString(serialString string) (*uint32, error)
	FromSerialIndex(si *models.SerialIndex) (string, error)
}

type TimeBasedGenerator struct{}

var _ Generator = &TimeBasedGenerator{}

// Generates a time-based serial number in the format YYYYMMDD
func (g *TimeBasedGenerator) GenerateBase() (*uint32, error) {
	t := time.Now()
	serialString := sprintf("%04d%02d%02d", t.Year(), t.Month(), t.Day())
	return g.FromString(serialString)
}

// Validates that the serialString contains a uint32 value in string form
func (g *TimeBasedGenerator) FromString(serialString string) (*uint32, error) {
	parsedSerial, err := parseUint(serialString, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("unable to generate a serial number from string '%s': %w", serialString, err)
	}

	// Explicitly convert to a uint32 to make sure it fits
	// We could just return the original string passed in but, just in case the strconv result is is a different number from the string
	// we'll convert the parsed serial back into an int32
	serial := uint32(parsedSerial)
	return &serial, nil
}

func (g *TimeBasedGenerator) FromSerialIndex(si *models.SerialIndex) (string, error) {
	if si == nil || si.Base == nil || si.ChangeIndex == nil {
		return "", fmt.Errorf("unable to convert SerialIndex to a serial number: '%v'", si)
	}
	serialString := fmt.Sprintf("%d%02d", *si.Base, *si.ChangeIndex)

	return serialString, nil
}
