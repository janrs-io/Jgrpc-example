// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: v1/productservice.proto

package productPBV1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on CreateRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *CreateRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in CreateRequestMultiError, or
// nil if none found.
func (m *CreateRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if l := utf8.RuneCountInString(m.GetName()); l < 1 || l > 255 {
		err := CreateRequestValidationError{
			field:  "Name",
			reason: "value length must be between 1 and 255 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetPrice() < 1 {
		err := CreateRequestValidationError{
			field:  "Price",
			reason: "value must be greater than or equal to 1",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Desc

	if l := utf8.RuneCountInString(m.GetTitle()); l < 1 || l > 100 {
		err := CreateRequestValidationError{
			field:  "Title",
			reason: "value length must be between 1 and 100 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Stock

	if _, ok := _CreateRequest_IsDisable_InLookup[m.GetIsDisable()]; !ok {
		err := CreateRequestValidationError{
			field:  "IsDisable",
			reason: "value must be in list [1 2]",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return CreateRequestMultiError(errors)
	}

	return nil
}

// CreateRequestMultiError is an error wrapping multiple validation errors
// returned by CreateRequest.ValidateAll() if the designated constraints
// aren't met.
type CreateRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateRequestMultiError) AllErrors() []error { return m }

// CreateRequestValidationError is the validation error returned by
// CreateRequest.Validate if the designated constraints aren't met.
type CreateRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateRequestValidationError) ErrorName() string { return "CreateRequestValidationError" }

// Error satisfies the builtin error interface
func (e CreateRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateRequestValidationError{}

var _CreateRequest_IsDisable_InLookup = map[int64]struct{}{
	1: {},
	2: {},
}

// Validate checks the field values on UpdateRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *UpdateRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in UpdateRequestMultiError, or
// nil if none found.
func (m *UpdateRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetId() < 1 {
		err := UpdateRequestValidationError{
			field:  "Id",
			reason: "value must be greater than or equal to 1",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.Name != nil {

		if utf8.RuneCountInString(m.GetName()) > 255 {
			err := UpdateRequestValidationError{
				field:  "Name",
				reason: "value length must be at most 255 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.Price != nil {
		// no validation rules for Price
	}

	if m.Desc != nil {

		if utf8.RuneCountInString(m.GetDesc()) > 255 {
			err := UpdateRequestValidationError{
				field:  "Desc",
				reason: "value length must be at most 255 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.Title != nil {

		if utf8.RuneCountInString(m.GetTitle()) > 100 {
			err := UpdateRequestValidationError{
				field:  "Title",
				reason: "value length must be at most 100 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.Stock != nil {
		// no validation rules for Stock
	}

	if m.IsDisable != nil {

		if _, ok := _UpdateRequest_IsDisable_InLookup[m.GetIsDisable()]; !ok {
			err := UpdateRequestValidationError{
				field:  "IsDisable",
				reason: "value must be in list [1 2]",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if len(errors) > 0 {
		return UpdateRequestMultiError(errors)
	}

	return nil
}

// UpdateRequestMultiError is an error wrapping multiple validation errors
// returned by UpdateRequest.ValidateAll() if the designated constraints
// aren't met.
type UpdateRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateRequestMultiError) AllErrors() []error { return m }

// UpdateRequestValidationError is the validation error returned by
// UpdateRequest.Validate if the designated constraints aren't met.
type UpdateRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateRequestValidationError) ErrorName() string { return "UpdateRequestValidationError" }

// Error satisfies the builtin error interface
func (e UpdateRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateRequestValidationError{}

var _UpdateRequest_IsDisable_InLookup = map[int64]struct{}{
	1: {},
	2: {},
}

// Validate checks the field values on DeleteRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *DeleteRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in DeleteRequestMultiError, or
// nil if none found.
func (m *DeleteRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetId() < 1 {
		err := DeleteRequestValidationError{
			field:  "Id",
			reason: "value must be greater than or equal to 1",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return DeleteRequestMultiError(errors)
	}

	return nil
}

// DeleteRequestMultiError is an error wrapping multiple validation errors
// returned by DeleteRequest.ValidateAll() if the designated constraints
// aren't met.
type DeleteRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteRequestMultiError) AllErrors() []error { return m }

// DeleteRequestValidationError is the validation error returned by
// DeleteRequest.Validate if the designated constraints aren't met.
type DeleteRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteRequestValidationError) ErrorName() string { return "DeleteRequestValidationError" }

// Error satisfies the builtin error interface
func (e DeleteRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteRequestValidationError{}

// Validate checks the field values on DetailRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *DetailRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DetailRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in DetailRequestMultiError, or
// nil if none found.
func (m *DetailRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *DetailRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetId() < 1 {
		err := DetailRequestValidationError{
			field:  "Id",
			reason: "value must be greater than or equal to 1",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return DetailRequestMultiError(errors)
	}

	return nil
}

// DetailRequestMultiError is an error wrapping multiple validation errors
// returned by DetailRequest.ValidateAll() if the designated constraints
// aren't met.
type DetailRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DetailRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DetailRequestMultiError) AllErrors() []error { return m }

// DetailRequestValidationError is the validation error returned by
// DetailRequest.Validate if the designated constraints aren't met.
type DetailRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DetailRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DetailRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DetailRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DetailRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DetailRequestValidationError) ErrorName() string { return "DetailRequestValidationError" }

// Error satisfies the builtin error interface
func (e DetailRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDetailRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DetailRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DetailRequestValidationError{}

// Validate checks the field values on ListRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ListRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ListRequestMultiError, or
// nil if none found.
func (m *ListRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ListRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Page

	// no validation rules for PageSize

	// no validation rules for Name

	if len(errors) > 0 {
		return ListRequestMultiError(errors)
	}

	return nil
}

// ListRequestMultiError is an error wrapping multiple validation errors
// returned by ListRequest.ValidateAll() if the designated constraints aren't met.
type ListRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListRequestMultiError) AllErrors() []error { return m }

// ListRequestValidationError is the validation error returned by
// ListRequest.Validate if the designated constraints aren't met.
type ListRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListRequestValidationError) ErrorName() string { return "ListRequestValidationError" }

// Error satisfies the builtin error interface
func (e ListRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListRequestValidationError{}

// Validate checks the field values on ListResponse with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ListResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ListResponseMultiError, or
// nil if none found.
func (m *ListResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *ListResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Total

	for idx, item := range m.GetList() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ListResponseValidationError{
						field:  fmt.Sprintf("List[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ListResponseValidationError{
						field:  fmt.Sprintf("List[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListResponseValidationError{
					field:  fmt.Sprintf("List[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ListResponseMultiError(errors)
	}

	return nil
}

// ListResponseMultiError is an error wrapping multiple validation errors
// returned by ListResponse.ValidateAll() if the designated constraints aren't met.
type ListResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListResponseMultiError) AllErrors() []error { return m }

// ListResponseValidationError is the validation error returned by
// ListResponse.Validate if the designated constraints aren't met.
type ListResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListResponseValidationError) ErrorName() string { return "ListResponseValidationError" }

// Error satisfies the builtin error interface
func (e ListResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListResponseValidationError{}

// Validate checks the field values on DecreaseStockRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DecreaseStockRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DecreaseStockRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DecreaseStockRequestMultiError, or nil if none found.
func (m *DecreaseStockRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *DecreaseStockRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetId() < 1 {
		err := DecreaseStockRequestValidationError{
			field:  "Id",
			reason: "value must be greater than or equal to 1",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetQuantity() < 1 {
		err := DecreaseStockRequestValidationError{
			field:  "Quantity",
			reason: "value must be greater than or equal to 1",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return DecreaseStockRequestMultiError(errors)
	}

	return nil
}

// DecreaseStockRequestMultiError is an error wrapping multiple validation
// errors returned by DecreaseStockRequest.ValidateAll() if the designated
// constraints aren't met.
type DecreaseStockRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DecreaseStockRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DecreaseStockRequestMultiError) AllErrors() []error { return m }

// DecreaseStockRequestValidationError is the validation error returned by
// DecreaseStockRequest.Validate if the designated constraints aren't met.
type DecreaseStockRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DecreaseStockRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DecreaseStockRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DecreaseStockRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DecreaseStockRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DecreaseStockRequestValidationError) ErrorName() string {
	return "DecreaseStockRequestValidationError"
}

// Error satisfies the builtin error interface
func (e DecreaseStockRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDecreaseStockRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DecreaseStockRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DecreaseStockRequestValidationError{}

// Validate checks the field values on ProductDetail with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ProductDetail) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ProductDetail with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ProductDetailMultiError, or
// nil if none found.
func (m *ProductDetail) ValidateAll() error {
	return m.validate(true)
}

func (m *ProductDetail) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for Name

	// no validation rules for Price

	// no validation rules for Desc

	// no validation rules for Title

	// no validation rules for Stock

	// no validation rules for IsDisable

	// no validation rules for CreateTime

	// no validation rules for UpdateTime

	if len(errors) > 0 {
		return ProductDetailMultiError(errors)
	}

	return nil
}

// ProductDetailMultiError is an error wrapping multiple validation errors
// returned by ProductDetail.ValidateAll() if the designated constraints
// aren't met.
type ProductDetailMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ProductDetailMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ProductDetailMultiError) AllErrors() []error { return m }

// ProductDetailValidationError is the validation error returned by
// ProductDetail.Validate if the designated constraints aren't met.
type ProductDetailValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ProductDetailValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ProductDetailValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ProductDetailValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ProductDetailValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ProductDetailValidationError) ErrorName() string { return "ProductDetailValidationError" }

// Error satisfies the builtin error interface
func (e ProductDetailValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sProductDetail.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ProductDetailValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ProductDetailValidationError{}

// Validate checks the field values on Response with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Response) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Response with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ResponseMultiError, or nil
// if none found.
func (m *Response) ValidateAll() error {
	return m.validate(true)
}

func (m *Response) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Code

	// no validation rules for Msg

	if all {
		switch v := interface{}(m.GetProtoAnyData()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ResponseValidationError{
					field:  "ProtoAnyData",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ResponseValidationError{
					field:  "ProtoAnyData",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetProtoAnyData()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ResponseValidationError{
				field:  "ProtoAnyData",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return ResponseMultiError(errors)
	}

	return nil
}

// ResponseMultiError is an error wrapping multiple validation errors returned
// by Response.ValidateAll() if the designated constraints aren't met.
type ResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ResponseMultiError) AllErrors() []error { return m }

// ResponseValidationError is the validation error returned by
// Response.Validate if the designated constraints aren't met.
type ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ResponseValidationError) ErrorName() string { return "ResponseValidationError" }

// Error satisfies the builtin error interface
func (e ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ResponseValidationError{}
