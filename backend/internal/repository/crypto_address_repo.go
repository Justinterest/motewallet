package repository

import (
	"context"

	"gorm.io/gorm"
	"motewallet/internal/model"
)

type CryptoAddressRepository interface {
	Create(ctx context.Context, addr *model.CryptoAddress) error
	FindByID(ctx context.Context, id uint64) (*model.CryptoAddress, error)
	ListByMerchant(ctx context.Context, merchantID uint64) ([]*model.CryptoAddress, error)
	Delete(ctx context.Context, id uint64) error
}

type cryptoAddressRepository struct {
	db *gorm.DB
}

func NewCryptoAddressRepository(db *gorm.DB) CryptoAddressRepository {
	return &cryptoAddressRepository{db: db}
}

func (r *cryptoAddressRepository) Create(ctx context.Context, addr *model.CryptoAddress) error {
	return r.db.WithContext(ctx).Create(addr).Error
}

func (r *cryptoAddressRepository) FindByID(ctx context.Context, id uint64) (*model.CryptoAddress, error) {
	var addr model.CryptoAddress
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&addr).Error
	if err != nil {
		return nil, err
	}
	return &addr, nil
}

func (r *cryptoAddressRepository) ListByMerchant(ctx context.Context, merchantID uint64) ([]*model.CryptoAddress, error) {
	var addrs []*model.CryptoAddress
	err := r.db.WithContext(ctx).Where("merchant_id = ?", merchantID).Order("id DESC").Find(&addrs).Error
	if err != nil {
		return nil, err
	}
	return addrs, nil
}

func (r *cryptoAddressRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.CryptoAddress{}, id).Error
}
