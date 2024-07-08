package dto

type AccountCreateRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=254"`
	Email string `json:"email" binding:"required,min=2,max=254,email"`
}

type AccountUpdateRequest struct {
	Name string `json:"name" binding:"required,max=254"`
}

type AccountOAuth2Request struct {
	Code string `json:"code" binding:"required,min=10,max=1024"`
}
