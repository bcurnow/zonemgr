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
	"fmt"
	"strconv"
	"time"
)

type SerialNumberGenerator interface {
	GenerateBaseSerial() (*uint32, error)
	GenerateSerial(index uint32) (*uint32, error)
	ToSerial(serialString string) (*uint32, error)
}

type TimeBasedSerialGenerator struct{}

// Generates a time-based serial number in the format YYYYMMDD
func (g *TimeBasedSerialGenerator) GenerateBaseSerial() (*uint32, error) {
	t := time.Now()
	serialString := fmt.Sprintf("%04d%02d%02d", t.Year(), t.Month(), t.Day())
	return g.ToSerial(serialString)
}

// Generates a time-based serial number plus a numeric index
func (g *TimeBasedSerialGenerator) GenerateSerial(index uint32) (*uint32, error) {
	baseSerial, err := g.GenerateBaseSerial()
	if err != nil {
		return nil, err
	}
	serialString := fmt.Sprintf("%d%02d", *baseSerial, index)
	return g.ToSerial(serialString)
}

// Validates that the serialString contains a uint32 value in string form
func (g *TimeBasedSerialGenerator) ToSerial(serialString string) (*uint32, error) {
	parsedSerial, err := strconv.ParseUint(serialString, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("unable to generate a serial number from string '%s': %w", serialString, err)
	}

	// Explicitly convert to a uint32 to make sure it fits
	// We could just return the original string passed in but, just in case the strconv result is is a different number from the string
	// we'll convert the parsed serial back into an int32
	serial := uint32(parsedSerial)
	return &serial, nil
}
