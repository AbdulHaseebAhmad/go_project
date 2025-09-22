package types

type Student struct {
	Id    int    `json:"id"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
	Age   int    `json:"age" validate:"required"`
}

type Credentials struct {
	Username    string `json:"username" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=6"`
	Phonenumber string `json:"phonenumber" validate:"required"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Token    string `json:"token"`
}
