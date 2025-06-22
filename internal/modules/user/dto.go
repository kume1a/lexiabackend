package user

type updateUserDTO struct {
	Name string `json:"name" valid:"optional"`
}
