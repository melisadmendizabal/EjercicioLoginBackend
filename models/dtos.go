package models

// --- Request DTOs ---

// LoginRequest define la estructura esperada para el request de login.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterRequest define la estructura para el registro (idéntica a LoginRequest en este ejemplo).
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// --- Response DTOs ---

// LoginSuccessData define la estructura para la data de login exitoso.
type LoginSuccessData struct {
	UserID   int64  `json:"userId"`
	Username string `json:"username"`
}

// RegisterSuccessData define la estructura para la data de registro exitoso.
type RegisterSuccessData struct {
	UserID   int64  `json:"userId"`
	Username string `json:"username"`
}

// UserResponse define la data pública que se muestra para un usuario.
type UserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// --- Estructuras Genéricas de Respuesta ---

// ErrorDetail provee detalles estructurados del error.
type ErrorDetail struct {
	Message string `json:"message"`
}

// APIResponse es un wrapper para estandarizar respuestas.
type APIResponse struct {
	Success bool         `json:"success"`
	Data    any          `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
}

// NewSuccessResponse crea una respuesta exitosa.
func NewSuccessResponse(data any) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
		Error:   nil,
	}
}

// NewErrorResponse crea una respuesta de error.
func NewErrorResponse(errorMessage string) APIResponse {
	return APIResponse{
		Success: false,
		Data:    nil,
		Error: &ErrorDetail{
			Message: errorMessage,
		},
	}
}
