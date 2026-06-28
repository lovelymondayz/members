package dto

// ─── Member DTOs ───────────────────────────────────────

type CreateMemberRequest struct {
	StoreID    string `json:"store_id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	MemberCode string `json:"member_code" binding:"required"`
	Tier       string `json:"tier"`
}

type UpdateMemberRequest struct {
	MemberCode string `json:"member_code"`
	Tier       string `json:"tier"`
}

// ─── Auth DTOs ─────────────────────────────────────────

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"`
}

// ─── Admin DTOs ────────────────────────────────────────

type CreateAdminRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// ─── Invoice DTOs ──────────────────────────────────────

type CreateInvoiceRequest struct {
	MemberID    string  `json:"member_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
	DueDate     string  `json:"due_date"`
}

type RecordPaymentRequest struct {
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	Method    string  `json:"method"`
	Reference string  `json:"reference"`
	Note      string  `json:"note"`
}
