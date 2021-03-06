// Copyright 2017 The ACH Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package ach

import (
	"fmt"
	"strings"
	"time"
)

// ErrServiceClass

// BatchHeader identifies the originating entity and the type of transactions
// contained in the batch (i.e., the standard entry class, PPD for consumer, CCD
// or CTX for corporate). This record also contains the effective date, or desired
// settlement date, for all entries contained in this batch. The settlement date
// field is not entered as it is determined by the ACH operato
type BatchHeader struct {
	// RecordType defines the type of record in the block. 5
	recordType string
	// ServiceClassCode ACH Mixed Debits and Credits ‘200’
	// ACH Credits Only ‘220’
	// ACH Debits Only ‘225'
	ServiceClassCode int
	// CompanyName the company originating the entries in the batch
	CompanyName string

	// CompanyDiscretionaryData allows Originators and/or ODFIs to include codes (one or more),
	// of significance only to them, to enable specialized handling of all
	// subsequent entries in that batch. There will be no standardized
	// interpretation for the value of the field. This field must be returned
	// intact on any return entry.
	CompanyDiscretionaryData string
	// CompanyIdentification The 9 digit FEIN number (proceeded by a predetermined
	// alpha or numeric character) of the entity in the company name field
	CompanyIdentification string

	// StandardEntryClassCode PPD’ for consumer transactions, ‘CCD’ or ‘CTX’ for corporate
	StandardEntryClassCode string

	// CompanyEntryDescription A description of the entries contained in the batch
	//
	//The Originator establishes the value of this field to provide a
	// description of the purpose of the entry to be displayed back to
	// the receive For example, "GAS BILL," "REG. SALARY," "INS. PREM,"
	// "SOC. SEC.," "DTC," "TRADE PAY," "PURCHASE," etc.
	//
	// This field must contain the word "REVERSAL" (left justified) when the
	// batch contains reversing entries.
	//
	// This field must contain the word "RECLAIM" (left justified) when the
	// batch contains reclamation entries.
	//
	// This field must contain the word "NONSETTLED" (left justified) when the
	// batch contains entries which could not settle.
	CompanyEntryDescription string

	// CompanyDescriptiveDate except as otherwise noted below, the Originator establishes this field
	// as the date it would like to see displayed to the receiver for
	// descriptive purposes. This field is never used to control timing of any
	// computer or manual operation. It is solely for descriptive purposes.
	// The RDFI should not assume any specific format. Examples of possible
	// entries in this field are "011392,", "01 92," "JAN 13," "JAN 92," etc.
	CompanyDescriptiveDate string

	// EffectiveEntryDate the date on which the entries are to settle
	EffectiveEntryDate time.Time

	// SettlementDate Leave blank, this field is inserted by the ACH operator
	settlementDate string

	// OriginatorStatusCode '1'
	OriginatorStatusCode int

	//ODFIIdentification First 8 digits of the originating DFI transit routing number
	ODFIIdentification int

	// BatchNumber is assigned in ascending sequence to each batch by the ODFI
	// or its Sending Point in a given file of entries. Since the batch number
	// in the Batch Header Record and the Batch Control Record is the same,
	// the ascending sequence number should be assigned by batch and not by
	// record.
	BatchNumber int
	// validator is composed for data validation
	validator
	// converters is composed for ACH to golang Converters
	converters
}

// NewBatchHeader returns a new BatchHeader with default valus for none exported fields
func NewBatchHeader() *BatchHeader {
	return &BatchHeader{
		recordType:           "5",
		OriginatorStatusCode: 1,
		BatchNumber:          1,
	}
}

// Parse takes the input record string and parses the BatchHeader values
func (bh *BatchHeader) Parse(record string) {
	// 1-1 Always "5"
	bh.recordType = "5"
	// 2-4 If the entries are credits, always "220". If the entries are debits, always "225"
	bh.ServiceClassCode = bh.parseNumField(record[1:4])
	// 5-20 Your company's name. This name may appear on the receivers’ statements prepared by the RDFI.
	bh.CompanyName = strings.TrimSpace(record[4:20])
	// 21-40 Optional field you may use to describe the batch for internal accounting purposes
	bh.CompanyDiscretionaryData = strings.TrimSpace(record[20:40])
	// 41-50 A 10-digit number assigned to you by the ODFI once they approve you to
	// originate ACH files through them. This is the same as the "Immediate origin" field in File Header Record
	bh.CompanyIdentification = strings.TrimSpace(record[40:50])
	// 51-53 If the entries are PPD (credits/debits towards consumer account), use "PPD".
	// If the entries are CCD (credits/debits towards corporate account), use "CCD".
	// The difference between the 2 class codes are outside of the scope of this post, but generally most ACH transfers to consumer bank accounts should use "PPD"
	bh.StandardEntryClassCode = record[50:53]
	// 54-63 Your description of the transaction. This text will appear on the receivers’ bank statement.
	// For example: "Payroll   "
	bh.CompanyEntryDescription = strings.TrimSpace(record[53:63])
	// 64-69 The date you choose to identify the transactions in YYMMDD format.
	// This date may be printed on the receivers’ bank statement by the RDFI
	bh.CompanyDescriptiveDate = strings.TrimSpace(record[63:69])
	// 70-75 Date transactions are to be posted to the receivers’ account.
	// You almost always want the transaction to post as soon as possible, so put tomorrow's date in YYMMDD format
	bh.EffectiveEntryDate = bh.parseSimpleDate(record[69:75])
	// 76-79 Always blank (just fill with spaces)
	bh.settlementDate = "   "
	// 79-79 Always 1
	bh.OriginatorStatusCode = bh.parseNumField(record[78:79])
	// 80-87 Your ODFI's routing number without the last digit. The last digit is simply a
	// checksum digit, which is why it is not necessary
	bh.ODFIIdentification = bh.parseNumField(record[79:87])
	// 88-94 Sequential number of this Batch Header Recor
	// For example, put "1" if this is the first Batch Header Record in the file
	bh.BatchNumber = bh.parseNumField(record[87:94])
}

// String writes the BatchHeader struct to a 94 character string.
func (bh *BatchHeader) String() string {
	return fmt.Sprintf("%v%v%v%v%v%v%v%v%v%v%v%v%v",
		bh.recordType,
		bh.ServiceClassCode,
		bh.CompanyNameField(),
		bh.CompanyDiscretionaryDataField(),
		bh.CompanyIdentificationField(),
		bh.StandardEntryClassCode,
		bh.CompanyEntryDescriptionField(),
		bh.CompanyDescriptiveDateField(),
		bh.EffectiveEntryDateField(),
		bh.settlementDateField(),
		bh.OriginatorStatusCode,
		bh.ODFIIdentificationField(),
		bh.BatchNumberField(),
	)
}

// Validate performs NACHA format rule checks on the record and returns an error if not Validated
// The first error encountered is returned and stops that parsing.
func (bh *BatchHeader) Validate() error {
	if err := bh.fieldInclusion(); err != nil {
		return err
	}
	if bh.recordType != "5" {
		return &ValidateError{FieldName: "recordType", Value: bh.recordType, Err: ErrRecordType}
	}
	if err := bh.isServiceClass(bh.ServiceClassCode); err != nil {
		return &ValidateError{FieldName: "ServiceClassCode", Value: string(bh.ServiceClassCode), Err: err}
	}
	if err := bh.isSECCode(bh.StandardEntryClassCode); err != nil {
		return &ValidateError{FieldName: "ServiceClassCode", Value: string(bh.ServiceClassCode), Err: err}
	}
	if err := bh.isOriginatorStatusCode(bh.OriginatorStatusCode); err != nil {
		return &ValidateError{FieldName: "OriginatorStatusCode", Value: string(bh.OriginatorStatusCode), Err: err}
	}
	if err := bh.isAlphanumeric(bh.CompanyName); err != nil {
		return &ValidateError{FieldName: "CompanyName", Value: string(bh.CompanyName), Err: err}
	}
	if err := bh.isAlphanumeric(bh.CompanyDiscretionaryData); err != nil {
		return &ValidateError{FieldName: "CompanyDiscretionaryData", Value: string(bh.CompanyDiscretionaryData), Err: err}
	}
	if err := bh.isAlphanumeric(bh.CompanyIdentification); err != nil {
		return &ValidateError{FieldName: "CompanyIdentification", Value: string(bh.CompanyIdentification), Err: err}
	}
	if err := bh.isAlphanumeric(bh.CompanyEntryDescription); err != nil {
		return &ValidateError{FieldName: "CompanyEntryDescription", Value: string(bh.CompanyEntryDescription), Err: err}
	}

	return nil
}

// fieldInclusion validate mandatory fields are not default values. If fields are
// invalid the ACH transfer will be returned.
func (bh *BatchHeader) fieldInclusion() error {
	if bh.recordType == "" {
		return &ValidateError{FieldName: "recordType", Value: bh.recordType, Err: ErrRecordType}
	}
	if bh.ServiceClassCode == 0 {
		return &ValidateError{FieldName: "ServiceClassCode", Value: string(bh.ServiceClassCode), Err: ErrServiceClass}
	}
	if bh.CompanyName == "" {
		return &ValidateError{FieldName: "CompanyName", Value: bh.CompanyName, Err: ErrValidFieldInclusion}
	}
	if bh.CompanyIdentification == "" {
		return &ValidateError{FieldName: "CompanyIdentification", Value: bh.CompanyIdentification, Err: ErrValidFieldInclusion}
	}
	if bh.StandardEntryClassCode == "" {
		return &ValidateError{FieldName: "StandardEntryClassCode", Value: bh.StandardEntryClassCode, Err: ErrSECCode}
	}
	if bh.CompanyEntryDescription == "" {
		return &ValidateError{FieldName: "CompanyEntryDescription", Value: bh.CompanyEntryDescription, Err: ErrValidFieldInclusion}
	}
	if bh.OriginatorStatusCode == 0 {
		return &ValidateError{FieldName: "OriginatorStatusCode", Value: string(bh.OriginatorStatusCode), Err: ErrValidFieldInclusion}
	}
	if bh.ODFIIdentification == 0 {
		return &ValidateError{FieldName: "ODFIIdentification", Value: string(bh.ODFIIdentification), Err: ErrValidFieldInclusion}
	}
	if bh.BatchNumber == 0 {
		return &ValidateError{FieldName: "BatchNumber", Value: string(bh.BatchNumber), Err: ErrValidFieldInclusion}
	}

	return nil
}

// CompanyNameField get the CompanyName left padded
func (bh *BatchHeader) CompanyNameField() string {
	return bh.alphaField(bh.CompanyName, 16)
}

// CompanyDiscretionaryDataField get the CompanyDiscretionaryData left padded
func (bh *BatchHeader) CompanyDiscretionaryDataField() string {
	return bh.alphaField(bh.CompanyDiscretionaryData, 20)
}

// CompanyIdentificationField get the CompanyIdentification left padded
func (bh *BatchHeader) CompanyIdentificationField() string {
	return bh.alphaField(bh.CompanyIdentification, 10)
}

// CompanyEntryDescriptionField get the CompanyEntryDescription left padded
func (bh *BatchHeader) CompanyEntryDescriptionField() string {
	return bh.alphaField(bh.CompanyEntryDescription, 10)
}

// CompanyDescriptiveDateField get the CompanyDescriptiveDate left padded
func (bh *BatchHeader) CompanyDescriptiveDateField() string {
	return bh.alphaField(bh.CompanyDescriptiveDate, 6)
}

// EffectiveEntryDateField get the EffectiveEntryDate in YYMMDD format
func (bh *BatchHeader) EffectiveEntryDateField() string {
	return bh.formatSimpleDate(bh.EffectiveEntryDate)
}

// ODFIIdentificationField get the odfi number zero padded
func (bh *BatchHeader) ODFIIdentificationField() string {
	return bh.numericField(bh.ODFIIdentification, 8)
}

// BatchNumberField get the batch number zero padded
func (bh *BatchHeader) BatchNumberField() string {
	return bh.numericField(bh.BatchNumber, 7)
}

func (bh *BatchHeader) settlementDateField() string {
	return bh.alphaField(bh.settlementDate, 3)
}
