package countrycodes

import (
	"fmt"
)

type Assignment int

const (
	// OfficiallyAssigned
	// http://en.wikipedia.org/wiki/ISO_3166-1_alpha-2#Officially_assigned_code_elements
	// Assigned to a country, territory, or area of geographical interest.
	OfficiallyAssigned Assignment = iota

	// UserAssigned
	// http://en.wikipedia.org/wiki/ISO_3166-1_alpha-2#User-assigned_code_elements
	// Free for assignment at the disposal of users.
	UserAssigned

	// ExceptionallyReserved
	// http://en.wikipedia.org/wiki/ISO_3166-1_alpha-2#Exceptional_reservations
	// Reserved on request for restricted use
	ExceptionallyReserved

	// TransitionallyReserved
	// http://en.wikipedia.org/wiki/ISO_3166-1_alpha-2#Transitional_reservations
	// Deleted from ISO 3166-1 but reserved transitionally
	TransitionallyReserved

	// IndeterminatelyReserved
	// http://en.wikipedia.org/wiki/ISO_3166-1_alpha-2#Indeterminate_reservations
	// Used in coding systems associated with ISO 3166-1
	IndeterminatelyReserved

	// NotUsed
	// http://en.wikipedia.org/wiki/ISO_3166-1_alpha-2#Codes_currently_agreed_not_to_use
	// Not used in ISO 3166-1 in deference to international property organization names
	//
	NotUsed

	// Invalid
	// A placeholder value
	//
	Invalid
)

func (a *Assignment) Valid() bool {
	return *a < NotUsed
}

func NewAssignment(a string) (Assignment, error) {
	switch a {
	case "OfficiallyAssigned":
		return OfficiallyAssigned, nil
	case "UserAssigned":
		return UserAssigned, nil
	case "ExceptionallyReserved":
		return ExceptionallyReserved, nil
	case "TransitionallyReserved":
		return TransitionallyReserved, nil
	case "IndeterminatelyReserved":
		return IndeterminatelyReserved, nil
	case "NotUsed":
		return NotUsed, nil
	}
	return Invalid, fmt.Errorf("invalid assignment")
}
