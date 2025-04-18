package service

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"time"

	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
	"github.com/Ranik23/avito-tech-spring/internal/repository"
	"github.com/Ranik23/avito-tech-spring/internal/repository/manager"
	"github.com/Ranik23/avito-tech-spring/pkg/errs"
)

type PVZService interface {
	CreatePVZ(ctx context.Context, city string) (*domain.Pvz, error)

	GetPVZSInfo(ctx context.Context, start time.Time, end time.Time, offset int, limit int) ([]domain.PvzInfo, error)
	GetPVZList(ctx context.Context) ([]domain.Pvz, error)

	AddProduct(ctx context.Context, pvzID string, product_type string) (*domain.Product, error)
	DeleteLastProduct(ctx context.Context, pvzID string) error

	StartReception(ctx context.Context, pvzID string) (*domain.Reception, error)
	CloseReception(ctx context.Context, pvzID string) (*domain.Reception, error)
}

type pvzService struct {
	logger        *slog.Logger
	pvzRepo       repository.PvzRepository
	receptionRepo repository.ReceptionRepository
	productRepo   repository.ProductRepository
	txManager     manager.TxManager
	cities        []string
}

// GetPVZList implements PVZService.
func (p *pvzService) GetPVZList(ctx context.Context) ([]domain.Pvz, error) {
	return p.pvzRepo.GetListPVZ(ctx)
}

func NewPVZService(pvzRepo repository.PvzRepository, receptionRepo repository.ReceptionRepository, cities []string,
	productRepo repository.ProductRepository, manager manager.TxManager, logger *slog.Logger) PVZService {
	return &pvzService{
		logger:        logger,
		pvzRepo:       pvzRepo,
		receptionRepo: receptionRepo,
		productRepo:   productRepo,
		txManager:     manager,
		cities:        cities,
	}
}

func (p *pvzService) AddProduct(ctx context.Context, pvzID string, productType string) (*domain.Product, error) {

	var product *domain.Product

	err := p.txManager.Do(ctx, func(txCtx context.Context) error {
		reception, err := p.receptionRepo.FindOpen(txCtx, pvzID)
		if err != nil {
			return errs.Wrap("failed to find the open reception", err)
		}

		if reception == nil {
			return ErrAllReceptionsClosed
		}

		product, err = p.productRepo.CreateProduct(txCtx, productType, reception.ID)
		if err != nil {
			return errs.Wrap("failed to create the product", err)
		}

		return nil
	})

	if err != nil {
		p.logger.Error("Failed to add the product to the open reception",
			slog.String("error", err.Error()),
			slog.String("pvzID", pvzID),
			slog.String("productType", productType))

		return nil, err
	}

	return product, nil
}

func (p *pvzService) CloseReception(ctx context.Context, pvzID string) (*domain.Reception, error) {

	var receptionToReturn *domain.Reception

	err := p.txManager.Do(ctx, func(txCtx context.Context) error {
		reception, err := p.receptionRepo.FindOpen(txCtx, pvzID)
		if err != nil {
			return errs.Wrap("failed to find the open", err)
		}

		if reception == nil {
			return ErrAllReceptionsClosed
		}

		err = p.receptionRepo.UpdateReceptionStatus(txCtx, reception.ID, "Closed")
		if err != nil {
			return errs.Wrap("failed to updata the status", err)
		}

		receptionToReturn = reception
		return nil
	})

	if err != nil {
		p.logger.Error("Failed to close the reception",
			slog.String("error", err.Error()),
			slog.String("pvzID", pvzID))

		return nil, err
	}

	return receptionToReturn, nil
}

func (p *pvzService) CreatePVZ(ctx context.Context, city string) (*domain.Pvz, error) {
	if !slices.Contains(p.cities, city) {
		return nil, ErrInvalidCity
	}

	var pvz *domain.Pvz

	err := p.txManager.Do(ctx, func(txCtx context.Context) error {
		var err error

		pvz, err = p.pvzRepo.CreatePVZ(txCtx, city)
		if err != nil {
			if errors.Is(err, repository.ErrAlreadyExists) {
				return ErrAlreadyExists
			}
			return errs.Wrap("failed to create PVZ", err)
		}
		return nil
	})
	if err != nil {
		p.logger.Error("Failed to do the transaction CreatePvz",
			slog.String("error", err.Error()),
			slog.String("city", city))

		return nil, err
	}

	return pvz, nil
}

func (p *pvzService) DeleteLastProduct(ctx context.Context, pvzID string) error {
	err := p.txManager.Do(ctx, func(txCtx context.Context) error {
		reception, err := p.receptionRepo.FindOpen(txCtx, pvzID)
		if err != nil {
			return errs.Wrap("failed to find the opne reception", err)
		}

		if reception == nil {
			return ErrAllReceptionsClosed
		}

		lastProduct, err := p.productRepo.FindTheLastProduct(txCtx, pvzID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrReceptionEmpty
			}
			return errs.Wrap("failed to find the last product", err)
		}

		err = p.productRepo.DeleteProduct(txCtx, lastProduct.ID)
		if err != nil {
			return errs.Wrap("failed to delete the product", err)
		}

		return nil
	})

	if err != nil {
		p.logger.Error("Failed to delete the last product",
			slog.String("error", err.Error()),
			slog.String("pvzID", pvzID))

		return err
	}

	return nil
}

func (p *pvzService) GetPVZSInfo(ctx context.Context, start time.Time, end time.Time, offset int, limit int) ([]domain.PvzInfo, error) {

	var pvZsInfo []domain.PvzInfo

	err := p.txManager.Do(ctx, func(txCtx context.Context) error {
		pvzs, err := p.pvzRepo.GetPVZS(txCtx, offset, limit)
		if err != nil {
			return errs.Wrap("failed to get pvzs", err)
		}
		for _, pvz := range pvzs {

			receptions, err := p.receptionRepo.GetReceptionsFiltered(txCtx, pvz.ID, start, end)
			if err != nil {
				return errs.Wrap("fail to get receptions filtered", err)
			}

			var receptionsInfo []domain.ReceptionInfo

			for _, reception := range receptions {
				products, err := p.productRepo.GetProducts(txCtx, reception.ID)
				if err != nil {
					return errs.Wrap("fail to get the products", err)
				}

				receptionsInfo = append(receptionsInfo, domain.ReceptionInfo{
					Reception: *reception,
					Products:  products,
				})
			}

			pvZsInfo = append(pvZsInfo, domain.PvzInfo{
				Pvz:        pvz,
				Receptions: receptionsInfo,
			})
		}

		return nil
	})
	if err != nil {
		p.logger.Error("Failed to get PVZS info",
			slog.String("error", err.Error()))
		return nil, err
	}

	return pvZsInfo, nil
}

// StartReception implements PVZService.
func (p *pvzService) StartReception(ctx context.Context, pvzID string) (*domain.Reception, error) {

	var receptionToReturn *domain.Reception

	err := p.txManager.Do(ctx, func(txCtx context.Context) error {
		reception, err := p.receptionRepo.FindOpen(txCtx, pvzID)
		if err != nil {
			return errs.Wrap("find open reception", err)
		}

		if reception != nil {
			return ErrAlreadyOpen
		}

		receptionToReturn, err = p.receptionRepo.CreateReception(txCtx, pvzID)
		if err != nil {
			return errs.Wrap("create reception", err)
		}

		return nil
	})

	if err != nil {
		p.logger.Error("Failed to start reception",
			slog.String("error", err.Error()),
			slog.String("pvzID", pvzID))

		return nil, err
	}

	return receptionToReturn, nil
}
