package dto

type ResponseType[T any] struct {
	Body T
}

type ListQuery struct {
	Page     int    `query:"page" json:"page" doc:"Page number, starting from 1" default:"1" minimum:"1"`
	PerPage  int    `query:"per_page" json:"per_page" doc:"Number of items per page" default:"10" minimum:"1" maximum:"200"`
	SortBy   string `query:"sort_by" json:"sort_by" doc:"Sort by field" default:"created_at"`
	SortDir  string `query:"sort_dir" json:"sort_dir" doc:"Sort direction, either 'asc' or 'desc'" enum:"asc,desc" default:"desc"`
	Filters  string `query:"filters" json:"filters" doc:"Filters in JSON" default:"{}"`
	Search   string `query:"search" json:"search" doc:"Search query" default:""`
	Includes string `query:"includes" json:"includes" doc:"Includes in JSON" default:"{}"`
}

type AuthHeader struct {
	Authorization string `header:"Authorization" doc:"Bearer Token of the user" required:"true"`
}
