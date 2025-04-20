package service

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"time"

	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
	"github.com/Ranik23/avito-tech-spring/internal/repository"
)

type PVZService interface {
	CreatePVZ(ctx context.Context, city string) (*domain.Pvz, error)

	GetPVZSInfo(ctx context.Context, start time.Time, end time.Time, offset int, limit int) ([]domain.PvzInfo, error)
	GetPVZList(ctx context.Context) ([]domain.Pvz, error)

	AddProduct(ctx context.Context, pvzID string, productType string) (*domain.Product, error)
	DeleteLastProduct(ctx context.Context, pvzID string) error

	StartReception(ctx context.Context, pvzID string) (*domain.Reception, error)
	CloseReception(ctx context.Context, pvzID string) (*domain.Reception, error)
}

type pvzService struct {
	logger        *slog.Logger
	pvzRepo       repository.PvzRepository
	receptionRepo repository.ReceptionRepository
	productRepo   repository.ProductRepository
	txManager     repository.TxManager
	cities        []string
}

func NewPVZService(pvzRepo repository.PvzRepository, receptionRepo repository.ReceptionRepository, cities []string,
	productRepo repository.ProductRepository, manager repository.TxManager, logger *slog.Logger) PVZService {
	return &pvzService{
		logger:        logger,
		pvzRepo:       pvzRepo,
		receptionRepo: receptionRepo,
		productRepo:   productRepo,
		txManager:     manager,
		cities:        cities,
	}
}

func (p *pvzService) GetPVZList(ctx context.Context) ([]domain.Pvz, error) {
	pvzs, err := p.pvzRepo.GetListOfPVZS(ctx)
	if err != nil {
		p.logger.Error("Failed to get PVZ list", slog.String("error", err.Error()))
		return nil, err
	}
	return pvzs, nil
}

func (p *pvzService) AddProduct(ctx context.Context, pvzID string, productType string) (*domain.Product, error) {
	var product *domain.Product

	err := p.txManager.Do(ctx, func(txCtx context.Context) error {
		reception, err := p.receptionRepo.FindOpen(txCtx, pvzID)
		if err != nil {
			p.logger.Error("Failed to find the open reception",
				slog.String("pvzID", pvzID), slog.String("productType", productType),
				slog.String("error", err.Error()))
			return err
		}

		if reception == nil {
			p.logger.Warn("No open reception found", slog.String("pvzID", pvzID),
				slog.String("productType", productType))
			return ErrAllReceptionsClosed
		}

		product, err = p.productRepo.CreateProduct(txCtx, productType, reception.ID)
		if err != nil {
			p.logger.Error("Failed to create the product", slog.String("pvzID", pvzID),
				slog.String("productType", productType), slog.String("error", err.Error()))
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *pvzService) CloseReception(ctx context.Context, pvzID string) (*domain.Reception, error) {
	var receptionToReturn *domain.Reception

	err := p.txManager.Do(ctx, func(txCtx context.Context) error {
		reception, err := p.receptionRepo.FindOpen(txCtx, pvzID)
		if err != nil {
			p.logger.Error("Failed to find the open reception", slog.String("pvzID", pvzID),
				slog.String("error", err.Error()))
			return err
		}

		if reception == nil {
			p.logger.Warn("No open reception found", slog.String("pvzID", pvzID))
			return ErrAllReceptionsClosed
		}

		err = p.receptionRepo.UpdateReceptionStatus(txCtx, reception.ID, "closed")
		if err != nil {
			p.logger.Error("Failed to update the reception status", slog.String("pvzID", pvzID),
				slog.String("error", err.Error()))
			return err
		}

		receptionToReturn = reception
		return nil
	})

	if err != nil {
		return nil, err
	}

	return receptionToReturn, nil
}

func (p *pvzService) CreatePVZ(ctx context.Context, city string) (*domain.Pvz, error) {
	if !slices.Contains(p.cities, city) {
		p.logger.Warn("Invalid city for PVZ creation", slog.String("city", city))
		return nil, ErrInvalidCity
	}

	var pvz *domain.Pvz

	err := p.txManager.Do(ctx, func(txCtx context.Context) error {
		var err error
		pvz, err = p.pvzRepo.CreatePVZ(txCtx, city)
		if err != nil {
			if errors.Is(err, repository.ErrAlreadyExists) {
				p.logger.Warn("PVZ already exists", slog.String("city", city))
				return ErrAlreadyExists
			}
			p.logger.Error("Failed to create PVZ", slog.String("city", city),
				slog.String("error", err.Error()))
			return err
		}
		return nil
	})
	if err != nil {
		p.logger.Error("Failed to CreatePVZ", slog.String("city", city),
			slog.String("error", err.Error()))
		return nil, err
	}

	return pvz, nil
}

func (p *pvzService) DeleteLastProduct(ctx context.Context, pvzID string) error {
	err := p.txManager.Do(ctx, func(txCtx context.Context) error {
		reception, err := p.receptionRepo.FindOpen(txCtx, pvzID)
		if err != nil {
			p.logger.Error("Failed to find the open reception", slog.String("pvzID", pvzID),
				slog.String("error", err.Error()))
			return err
		}

		if reception == nil {
			p.logger.Warn("No open reception found", slog.String("pvzID", pvzID))
			return ErrAllReceptionsClosed
		}

		lastProduct, err := p.productRepo.FindTheLastProduct(txCtx, pvzID)
		if err != nil {
			p.logger.Error("Failed to find the last product", slog.String("pvzID", pvzID),
				slog.String("error", err.Error()))
			return err
		}

		if lastProduct == nil {
			return ErrReceptionEmpty
		}

		err = p.productRepo.DeleteProduct(txCtx, lastProduct.ID)
		if err != nil {
			p.logger.Error("Failed to delete the last product", slog.String("pvzID", pvzID),
				slog.String("error", err.Error()))
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (p *pvzService) GetPVZSInfo(ctx context.Context, start time.Time, end time.Time, offset int, limit int) ([]domain.PvzInfo, error) {
	var pvZsInfo []domain.PvzInfo

	err := p.txManager.Do(ctx, func(txCtx context.Context) error {
		pvzs, err := p.pvzRepo.GetPVZS(txCtx, offset, limit)
		if err != nil {
			p.logger.Error("Failed to get PVZS", slog.String("error", err.Error()))
			return err
		}
		for _, pvz := range pvzs {
			receptions, err := p.receptionRepo.GetReceptionsFiltered(txCtx, pvz.ID, start, end)
			if err != nil {
				p.logger.Error("Failed to get filtered receptions", slog.String("pvzID", pvz.ID),
					slog.String("error", err.Error()))
				return err
			}

			var receptionsInfo []domain.ReceptionInfo
			for _, reception := range receptions {
				products, err := p.productRepo.GetProducts(txCtx, reception.ID)
				if err != nil {
					p.logger.Error("Failed to get products", slog.String("receptionID", reception.ID),
						slog.String("error", err.Error()))
					return err
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
		return nil, err
	}

	return pvZsInfo, nil
}

func (p *pvzService) StartReception(ctx context.Context, pvzID string) (*domain.Reception, error) {
	var receptionToReturn *domain.Reception

	err := p.txManager.Do(ctx, func(txCtx context.Context) error {
		reception, err := p.receptionRepo.FindOpen(txCtx, pvzID)
		if err != nil {
			p.logger.Error("Failed to find open reception", slog.String("pvzID", pvzID),
				slog.String("error", err.Error()))
			return err
		}

		if reception != nil {
			p.logger.Warn("Reception already open", slog.String("pvzID", pvzID))
			return ErrAlreadyOpen
		}

		receptionToReturn, err = p.receptionRepo.CreateReception(txCtx, pvzID)
		if err != nil {
			p.logger.Error("Failed to create reception", slog.String("pvzID", pvzID),
				slog.String("error", err.Error()))
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return receptionToReturn, nil
}
