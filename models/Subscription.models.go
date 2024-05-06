package models

import (
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v78"
	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model
	ID string // Use ID from stripe

	UserID string
	User   User

	Status stripe.SubscriptionStatus

	MemoryGB  int
	StorageGB int
	ArchiveGB int

	MemoryPriceID uuid.UUID
	MemoryPrice   Price

	StoragePriceID uuid.UUID
	StoragePrice   Price

	ArchivePriceID uuid.UUID
	ArchivePrice   Price
}
