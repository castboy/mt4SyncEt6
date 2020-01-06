// Copyright 2016 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package decimal

import (
	"mt4SyncEt6/decimal/errors"

	"mt4SyncEt6/decimal/terror"
)

var (
	// ErrTruncated is returned when data has been truncated during conversion.
	ErrTruncated = terror.ClassTypes.New(codeTruncated, "Data Truncated")
	// ErrTruncatedWrongVal is returned when data has been truncated during conversion.
	// ErrOverflow is returned when data is out of range for a field type.
	ErrOverflow = terror.ClassTypes.New(codeOverflow, msgOverflow)
	// ErrDivByZero is return when do division by 0.
	ErrDivByZero = terror.ClassTypes.New(codeDivByZero, "Division by 0")
	// ErrTooBigDisplayWidth is return when display width out of range for column.
	ErrTooBigDisplayWidth = terror.ClassTypes.New(codeTooBigDisplayWidth, "Too Big Display width")
	// ErrTooBigFieldLength is return when column length too big for column.
	ErrTooBigFieldLength = terror.ClassTypes.New(codeTooBigFieldLength, "Too Big Field length")
	// ErrBadNumber is return when parsing an invalid binary decimal number.
	ErrBadNumber = terror.ClassTypes.New(codeBadNumber, "Bad Number")
	// ErrCastAsSignedOverflow is returned when positive out-of-range integer, and convert to it's negative complement.
	ErrCastAsSignedOverflow = terror.ClassTypes.New(codeUnknown, msgCastAsSignedOverflow)
	// ErrCastNegIntAsUnsigned is returned when a negative integer be casted to an unsigned int.
	ErrCastNegIntAsUnsigned = terror.ClassTypes.New(codeUnknown, msgCastNegIntAsUnsigned)
	// ErrInvalidTimeFormat is returned when the time format is not correct.
	ErrInvalidTimeFormat = terror.ClassTypes.New(ErrTruncatedWrongValue, "invalid time format: '%v'")
	// ErrInvalidWeekModeFormat is returned when the week mode is wrong.
	ErrInvalidWeekModeFormat = terror.ClassTypes.New(ErrTruncatedWrongValue, "invalid week mode format: '%v'")
	// ErrInvalidYearFormat is returned when the input is not a valid year format.
	ErrInvalidYearFormat = errors.New("invalid year format")
	// ErrInvalidYear is returned when the input value is not a valid year.
	ErrInvalidYear = errors.New("invalid year")
)

const (
	codeBadNumber terror.ErrCode = 1

	codeTruncated          = terror.ErrCode(WarnDataTruncated)
	codeOverflow           = terror.ErrCode(ErrDataOutOfRange)
	codeDivByZero          = terror.ErrCode(ErrDivisionByZero)
	codeTooBigDisplayWidth = terror.ErrCode(ErrTooBigDisplaywidth)
	codeTooBigFieldLength  = terror.ErrCode(ErrTooBigFieldlength)
	codeUnknown            = terror.ErrCode(ErrUnknown)
)

var (
	msgOverflow             = MySQLErrName[ErrDataOutOfRange]
	msgCastAsSignedOverflow = "Cast to signed converted positive out-of-range integer to it's negative complement"
	msgCastNegIntAsUnsigned = "Cast to unsigned converted negative integer to it's positive complement"
)

func init() {
	typesMySQLErrCodes := map[terror.ErrCode]uint16{
		codeTruncated:          WarnDataTruncated,
		codeOverflow:           ErrDataOutOfRange,
		codeDivByZero:          ErrDivisionByZero,
		codeTooBigDisplayWidth: ErrTooBigDisplaywidth,
		codeTooBigFieldLength:  ErrTooBigFieldlength,
		codeUnknown:            ErrUnknown,
	}
	terror.ErrClassToMySQLCodes[terror.ClassTypes] = typesMySQLErrCodes
}
