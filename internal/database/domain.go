package database

import (
	"time"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ClientDB struct {
	gorm.Model
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	FullName     string
	Email        string    `gorm:"unique"`
	Phone        string
	PasswordHash string
	Photo        string
	Type         string
	Permissions  pq.StringArray `gorm:"type:text[]"`
	Balance      float64   `gorm:"default:0"`
	Orders       []Order   `gorm:"foreignKey:ClientID"`
	Reviews      []Review  `gorm:"foreignKey:ClientID"`
}

type CompanyDB struct {
	gorm.Model
	ID            uint      `gorm:"primaryKey;autoIncrement"`
	CompanyName   string
	FullName      string
	PositionAgent string
	IDCompany     string
	Email         string   `gorm:"unique"`
	Phone         string
	Address       string
	TypeService   string
	PasswordHash  string
	Photo         string
	Documents     pq.StringArray `gorm:"type:text[]"`
	Stars         float64  `gorm:"default:0"`
	ReviewCount   int      `gorm:"default:0"`
	Type          string
	Permissions   pq.StringArray `gorm:"type:text[]"`
	Balance       float64  `gorm:"default:0"`
	Cards         []Card   `gorm:"foreignKey:CompanyID"`
	Orders        []Order  `gorm:"foreignKey:CompanyID"`
	Reviews       []Review `gorm:"foreignKey:CompanyID"`
}

type Card struct {
	gorm.Model
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Title       string
	Description string
	Price       float64
	Location    string
	Category    string
	IsActive    bool      `gorm:"default:true"`
	CompanyID   uint
	Company     CompanyDB `gorm:"foreignKey:CompanyID"`
	Orders      []Order   `gorm:"foreignKey:CardID"`
}

type Order struct {
	gorm.Model
	ID                uint              `gorm:"primaryKey;autoIncrement"`
	ClientID          uint
	Client            ClientDB          `gorm:"foreignKey:ClientID"`
	CompanyID         uint
	Company           CompanyDB         `gorm:"foreignKey:CompanyID"`
	CardID            uint
	Card              Card              `gorm:"foreignKey:CardID"`
	Amount            float64
	Status            string            `gorm:"default:'created'"` // created, paid, in_progress, completed, finished, cancelled
	PaymentStatus     string            `gorm:"default:'pending'"` // pending, paid, refunded
	Description       string
	WorkerCompleteURL string            // Одноразовая ссылка для работника
	EscrowTransactions []EscrowTransaction `gorm:"foreignKey:OrderID"`
	Notifications     []Notification    `gorm:"foreignKey:OrderID"`
	CompletedAt       *time.Time
}

type EscrowTransaction struct {
	gorm.Model
	ID        uint    `gorm:"primaryKey;autoIncrement"`
	OrderID   uint
	Order     Order   `gorm:"foreignKey:OrderID"`
	Amount    float64
	Type      string  // hold, release, refund
	Status    string  // pending, completed, failed
	FromUser  string  // client, company, system
	ToUser    string  // client, company, escrow
}

type Review struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	ClientID  uint
	Client    ClientDB  `gorm:"foreignKey:ClientID"`
	CompanyID uint
	Company   CompanyDB `gorm:"foreignKey:CompanyID"`
	OrderID   uint
	Order     Order     `gorm:"foreignKey:OrderID"`
	Rating    int       `gorm:"check:rating >= 1 AND rating <= 5"`
	Comment   string
}

type Notification struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	UserID    uint   // ID клиента или компании
	UserType  string // client, company
	OrderID   *uint
	Order     *Order `gorm:"foreignKey:OrderID"`
	Title     string
	Message   string
	IsRead    bool   `gorm:"default:false"`
	Type      string // order_update, payment, review, system
}

type BalanceTransaction struct {
	gorm.Model
	ID       uint    `gorm:"primaryKey;autoIncrement"`
	UserID   uint    // ID клиента или компании
	UserType string  // client, company
	Amount   float64
	Type     string  // deposit, withdrawal, payment, refund
	Status   string  // pending, completed, failed
	OrderID  *uint   // Связь с заказом, если транзакция связана с заказом
	Description string
}

type WorkerLink struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	OrderID   uint
	Order     Order     `gorm:"foreignKey:OrderID"`
	Token     string    `gorm:"unique"`
	IsUsed    bool      `gorm:"default:false"`
	ExpiresAt time.Time
}
