package domain

import (
	"gorm.io/gorm"
)

type ProductQueryParams struct {
	SellerId    *int
	CategoryID  *int
	MinPrice    *int
	MaxPrice    *int
	IsActive    *bool
	Disposable  *bool
	SearchQuery *string
	Limit       *int
	Offset      *int
}

type ProductUpdateData struct {
	Title       *string
	Description *string
	Price       *int64
	CategoryID  *int
	IsActive    *bool
	Disposable  *bool
	MaxSales    *int
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id int) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id int) error
	GetAll() ([]User, error)
	AssignRole(userID int, roleName string) error
	RemoveRole(userID int, roleName string) error
	GetRoles(userID int) ([]string, error)
	LastLoginUpdate(userID int) error
}

type RoleRepository interface {
	Create(role *Role) error
	GetByID(id int) (*Role, error)
	GetByName(name string) (*Role, error)
	Update(role *Role) error
	Delete(id int) error
	GetAll() ([]Role, error)
	AssignRoleToUser(userID int, roleName string) error
	GetUserRoles(userID int) ([]Role, error)
}

type ProductRepository interface {
	Create(product *Product) (int, error)
	GetById(id int) (*Product, error)
	Update(id int, data ProductUpdateData) error
	Delete(id int) error
	IsAvailable(id int) (bool, error)
	IncrementSoldCount(id int, delta int) error
	ReserveProduct(id int) error
	CheckAvailabilityAndIncrementSoldCount(id int, delta int) (bool, error)
	Find(params ProductQueryParams) ([]*Product, error)
	Count(params ProductQueryParams) (int, error)
	GetAll() ([]*Product, error)
	GetBySellerID(sellerID int) ([]*Product, error)
}

type CategoryRepository interface {
	Create(category *Category) error
	GetCategoryById(id int) (*Category, error)
	GetByID(id int) (*Category, error)
	Update(category *Category) (*Category, error)
	Delete(category *Category) error
	GetAllCategories() (*[]Category, error)
}

type ConversationRepository interface {
	Create(conversation *Conversation) error
	GetByID(id int) (*Conversation, error)
	GetByUserID(userID int) ([]*Conversation, error)
	GetByProductID(productID int) ([]*Conversation, error)
	Update(conversation *Conversation) error
	Delete(id int) error
	AddParticipant(conversationID, userID int) error
	RemoveParticipant(conversationID, userID int) error
	BeginTx() (*gorm.DB, error)
	CreateTx(tx *gorm.DB, conversation *Conversation) error
	AddParticipantTx(tx *gorm.DB, conversationID, userID int) error
	Commit(tx *gorm.DB) error
	Rollback(tx *gorm.DB) error
}

type ParticipantRepository interface {
	Create(participant *Participant) error
	GetByConversationID(conversationID int) ([]*Participant, error)
	GetByUserID(userID int) ([]*Participant, error)
	Delete(conversationID, userID int) error
	IsParticipant(conversationID, userID int) (bool, error)
}

type OrderRepository interface {
	Create(order *Order) error
	GetByID(id int) (*Order, error)
	GetByCustomerID(customerID int) ([]*Order, error)
	GetBySellerID(sellerID int) ([]*Order, error)
	GetByProductID(productID int) ([]*Order, error)
	GetByStatus(status string) ([]*Order, error)
	GetAll() ([]*Order, error)
	Update(order *Order) error
	Delete(id int) error
	UpdateStatus(id int, status string) error
}

type MessageRepository interface {
	Create(message *Message) error
	GetByID(id int) (*Message, error)
	GetByConversationID(conversationID int, limit, offset int) ([]*Message, error)
	GetByOrderID(orderID int) ([]*Message, error)
	Update(message *Message) error
	Delete(id int) error
	GetSystemMessages(conversationID int) ([]*Message, error)
}
