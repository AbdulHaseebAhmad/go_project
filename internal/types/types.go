package types

type Student struct {
	Id    int    `json:"id"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
	Age   int    `json:"age" validate:"required"`
}

type Credentials struct {
	Username    string `json:username`
	Email       string `json:email`
	Password    string `json:password`
	Phonenumber string `json:phonenumber`
}
