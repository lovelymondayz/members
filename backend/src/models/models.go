package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─── Roles ──────────────────────────────────────────────

type Role struct {
	RoleID    uint      `gorm:"primaryKey" json:"role_id"`
	Name      string    `gorm:"unique;not null;size:50" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Role) TableName() string { return "roles" }

const (
	RoleSuperAdmin = "super_admin"
	RoleAdmin      = "admin"
	RoleMember     = "member"
)

// ─── Users ──────────────────────────────────────────────

type User struct {
	UserID     uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"user_id"`
	RoleID     uint           `gorm:"not null;index" json:"role_id"`
	Email      string         `gorm:"unique;not null;size:255" json:"email"`
	Password   string         `gorm:"size:255" json:"-"`
	GoogleID   *string        `gorm:"unique;size:255" json:"google_id,omitempty"`
	Name       string         `gorm:"not null;size:255" json:"name"`
	AvatarURL  string         `gorm:"size:255;column:avatar_url" json:"avatar_url,omitempty"`
	IsActive   bool           `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// Role is loaded separately to avoid circular FK dependency

func (User) TableName() string { return "users" }

// ─── Stores ─────────────────────────────────────────────

type Store struct {
	StoreID      uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"store_id"`
	AdminID      uuid.UUID      `gorm:"type:uuid;unique;not null;index" json:"admin_id"`
	Name         string         `gorm:"not null;size:255" json:"name"`
	LogoURL      string         `gorm:"size:255;column:logo_url" json:"logo_url,omitempty"`
	Address      string         `gorm:"type:text" json:"address,omitempty"`
	Phone        string         `gorm:"size:50" json:"phone,omitempty"`
	CardColorHex string         `gorm:"default:'#1E40AF';size:7" json:"card_color_hex"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// Admin is loaded separately

func (Store) TableName() string { return "stores" }

// ─── Members ────────────────────────────────────────────

type Member struct {
	MemberID       uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"member_id"`
	StoreID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"store_id"`
	UserID         *uuid.UUID     `gorm:"type:uuid;unique;index" json:"user_id,omitempty"`
	MemberCode     string         `gorm:"unique;not null;size:50" json:"member_code"`
	Tier           string         `gorm:"default:'standard';size:50" json:"tier"`
	JoinedAt       *time.Time     `json:"joined_at,omitempty"`
	UserName       string         `gorm:"-" json:"user_name,omitempty"`
	StoreName      string         `gorm:"-" json:"store_name,omitempty"`
	StoreCardColor string         `gorm:"-" json:"store_card_color,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Member) TableName() string { return "members" }

// ─── Invoices ───────────────────────────────────────────

type InvoiceStatus string

const (
	InvoiceDraft    InvoiceStatus = "draft"
	InvoiceSent     InvoiceStatus = "sent"
	InvoicePaid     InvoiceStatus = "paid"
	InvoiceOverdue  InvoiceStatus = "overdue"
	InvoiceCancelled InvoiceStatus = "cancelled"
)

type Invoice struct {
	InvoiceID     uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"invoice_id"`
	StoreID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"store_id"`
	MemberID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"member_id"`
	InvoiceNumber string         `gorm:"unique;not null;size:50" json:"invoice_number"`
	Amount        float64        `gorm:"type:decimal(12,2);not null" json:"amount"`
	Description   string         `gorm:"type:text" json:"description,omitempty"`
	Status        InvoiceStatus  `gorm:"default:'draft';size:20;index" json:"status"`
	DueDate       *time.Time     `json:"due_date,omitempty"`
	MemberName    string         `gorm:"-" json:"member_name,omitempty"`
	MemberCode    string         `gorm:"-" json:"member_code,omitempty"`
	Payments      []Payment      `gorm:"-" json:"payments,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (Invoice) TableName() string { return "invoices" }

// ─── Payments ───────────────────────────────────────────

type Payment struct {
	PaymentID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"payment_id"`
	InvoiceID uuid.UUID `gorm:"not null;index" json:"invoice_id"`
	Amount    float64   `gorm:"type:decimal(12,2);not null" json:"amount"`
	Method    string    `gorm:"default:'manual';size:50" json:"method"`
	Reference string    `gorm:"size:255" json:"reference,omitempty"`
	Note      string    `gorm:"type:text" json:"note,omitempty"`
	PaidAt    time.Time `json:"paid_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (Payment) TableName() string { return "payments" }
