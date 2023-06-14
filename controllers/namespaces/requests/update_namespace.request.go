package requests

// UpdateNamespaceByIDRequest is a struct that represents the request body for updating a namespace
type UpdateNamespaceByIDRequest struct {
	Description string `json:"description" binding:"omitempty,min=3,max=1000"`
	UserID      string `json:"userId" binding:"required"`
}
