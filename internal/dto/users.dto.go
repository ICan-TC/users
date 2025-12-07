package dto

type CreateUserReq struct {
	Body struct {
		Username string `json:"username" doc:"Username of the user" minLength:"3" MaxLength:"255" required:"true"`
		Email    string `json:"email" doc:"Email of the user" Email:"true" required:"true" format:"email"`
		Password string `json:"password" doc:"Password of the user" MinLength:"8" MaxLength:"255" required:"true"`
	}
}

type CreateUserRes struct{ Body CreateUserResBody }
type CreateUserResBody struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt uint   `json:"created_at"`
	UpdatedAt uint   `json:"updated_at"`
}

type UpdateUserReq struct {
	Body struct {
		Username string `json:"username" doc:"Username of the user" MinLength:"3" MaxLength:"255"`
		Email    string `json:"email" doc:"Email of the user" Email:"true"`
		Password string `json:"password" doc:"Password of the user" MinLength:"8" MaxLength:"255"`
	}
}
type UpdateUserRes struct{ Body UpdateUserResBody }
type UpdateUserResBody struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt uint   `json:"created_at"`
	UpdatedAt uint   `json:"updated_at"`
}

type GetUserByFieldReq struct {
	Body struct {
		Field string `query:"field" doc:"Field to find the user by" enums:"username,emaid,username,email"`
		Value string `query:"value" doc:"Value to find the user by"`
	}
}
type GetUserByFieldRes struct{ Body GetUserResBody }

type GetUserByIDReq struct {
	ID string `path:"id" doc:"ID of the user" required:"true"`
}
type GetUserByIDRes struct{ Body GetUserResBody }

type GetUserResBody struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt uint   `json:"created_at"`
	UpdatedAt uint   `json:"updated_at"`
}

type DeleteUserReq struct {
	ID string `path:"id" doc:"ID of the user" required:"true"`
}

type DeleteUserRes struct {
	Body struct {
		ID string `json:"id"`
	}
}

type ListUsersReq struct {
	ListQuery
}

type ListUsersRes struct {
	Body struct {
		Users     []GetUserResBody `json:"users"`
		Total     int              `json:"total"`
		ListQuery ListQuery        `json:"query"`
	}
}
