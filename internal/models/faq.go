package models

// FAQ Request DTOs
type FAQItemRequest struct {
	Question   string `json:"question" binding:"required,min=5,max=500"`
	QuestionEN string `json:"question_en" binding:"required,min=5,max=500"`
	Answer     string `json:"answer" binding:"required,min=10,max=2000"`
	AnswerEN   string `json:"answer_en" binding:"required,min=10,max=2000"`
}

type FAQUpdateRequest struct {
	Judul    *string          `json:"judul" binding:"omitempty,min=2,max=200"`
	JudulEN  *string          `json:"judul_en" binding:"omitempty,min=2,max=200"`
	IsActive *bool            `json:"is_active"`
	Items    []FAQItemRequest `json:"items" binding:"omitempty,min=1,dive"`
}

type FAQReorderRequest struct {
	Index     int    `json:"index" binding:"min=0"`
	Direction string `json:"direction" binding:"required,oneof=up down"`
}

// FAQ Response DTOs
type FAQItemResponse struct {
	Index      int    `json:"index"`
	Question   string `json:"question"`
	QuestionEN string `json:"question_en"`
	Answer     string `json:"answer"`
	AnswerEN   string `json:"answer_en"`
}

type FAQAdminResponse struct {
	ID        string            `json:"id"`
	Judul     string            `json:"judul"`
	JudulEN   string            `json:"judul_en"`
	Slug      string            `json:"slug"`
	IsActive  bool              `json:"is_active"`
	Items     []FAQItemResponse `json:"items"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}
