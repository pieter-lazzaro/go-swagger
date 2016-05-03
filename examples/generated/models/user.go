package models

import "github.com/go-openapi/strfmt"

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

/*User user

swagger:model User
*/
type User struct {

	/* email
	 */
	Email string `json:"email,omitempty"`

	/* first name
	 */
	FirstName string `json:"firstName,omitempty"`

	/* id
	 */
	ID int64 `json:"id,omitempty"`

	/* last name
	 */
	LastName string `json:"lastName,omitempty"`

	/* password
	 */
	Password string `json:"password,omitempty"`

	/* phone
	 */
	Phone string `json:"phone,omitempty"`

	/* User Status
	 */
	UserStatus int32 `json:"userStatus,omitempty"`

	/* username
	 */
	Username string `json:"username,omitempty"`
}

// Validate validates this user
func (m *User) Validate(formats strfmt.Registry) error {
	return nil
}
