package dto

import "time"

type CreateUserReq struct {
	AuthHeader
	Body struct {
		Username string `json:"username" doc:"Username of the user" minLength:"3" maxLength:"255" required:"true"`
		Email    string `json:"email" doc:"Email of the user" Email:"true" required:"true" format:"email"`
		Password string `json:"password" doc:"Password of the user" minLength:"8" maxLength:"255" required:"true"`

		FirstName   *string `json:"first_name" doc:"First name of the user" required:"false"`
		FamilyName  *string `json:"family_name" doc:"Family name of the user" required:"false"`
		PhoneNumber *string `json:"phone_number" doc:"Phone number of the user" required:"false"`
		DateOfBirth *string `json:"date_of_birth" doc:"Date of birth of the user" required:"false"`
	}
}

type CreateUserRes struct{ Body CreateUserResBody }
type CreateUserResBody struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`

	FirstName   *string `json:"first_name"`
	FamilyName  *string `json:"family_name"`
	PhoneNumber *string `json:"phone_number"`
	DateOfBirth *string `json:"date_of_birth"`

	CreatedAt int `json:"created_at"`
	UpdatedAt int `json:"updated_at"`
}

type UpdateUserReq struct {
	AuthHeader
	Body struct {
		Id          string     `json:"id" doc:"ID of the user" required:"true"`
		Username    *string    `json:"username,omitempty" doc:"Username of the user" minLength:"3" maxLength:"255"`
		Email       *string    `json:"email" doc:"Email of the user" format:"email" required:"false"`
		Password    *string    `json:"password" doc:"Password of the user" minLength:"8" maxLength:"255" required:"false"`
		FirstName   *string    `json:"first_name" doc:"First name of the user" required:"false"`
		FamilyName  *string    `json:"family_name" doc:"Family name of the user" required:"false"`
		PhoneNumber *string    `json:"phone_number" doc:"Phone number of the user" required:"false"`
		DateOfBirth *time.Time `json:"date_of_birth" doc:"Date of birth of the user" required:"false"`
	}
}
type UpdateUserRes struct{ Body UpdateUserResBody }
type UpdateUserResBody struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`

	FirstName   *string    `json:"first_name"`
	FamilyName  *string    `json:"family_name"`
	PhoneNumber *string    `json:"phone_number"`
	DateOfBirth *time.Time `json:"date_of_birth"`

	CreatedAt int `json:"created_at"`
	UpdatedAt int `json:"updated_at"`
}

type GetUserByFieldReq struct {
	AuthHeader
	Field string `path:"field" doc:"Field to find the user by" enum:"username,email,id" required:"true"`
	Value string `path:"value" doc:"Value to find the user by" required:"true"`
}
type GetUserByFieldRes struct{ Body GetUserResBody }

type GetUserByIDReq struct {
	AuthHeader
	ID string `path:"id" doc:"ID of the user" required:"true"`
}
type GetUserByIDRes struct{ Body GetUserResBody }

type GetUserResBody struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`

	FirstName   *string    `json:"first_name"`
	FamilyName  *string    `json:"family_name"`
	PhoneNumber *string    `json:"phone_number"`
	DateOfBirth *time.Time `json:"date_of_birth"`

	CreatedAt int `json:"created_at"`
	UpdatedAt int `json:"updated_at"`
}

type DeleteUserReq struct {
	AuthHeader
	ID string `path:"id" doc:"ID of the user" required:"true"`
}

type DeleteUserResBody struct {
	ID string `json:"id"`
}
type DeleteUserRes struct {
	Body DeleteUserResBody
}

type ListUsersReq struct {
	AuthHeader
	ListQuery
}

type ListUsersResBody struct {
	Users     []GetUserResBody `json:"users"`
	Total     int              `json:"total"`
	ListQuery ListQuery        `json:"query"`
}

type ListUsersRes struct {
	Body ListUsersResBody
}
