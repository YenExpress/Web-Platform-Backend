package dto

type RegisterDeveloperBody struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
}

// type LoginBody struct {
// 	Email     string `json:"email" binding:"required,email"`
// 	Password  string `json:"password" binding:"required"`
// 	IPAddress string
// }
