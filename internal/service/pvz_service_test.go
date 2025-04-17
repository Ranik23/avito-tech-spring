package service

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
	"github.com/Ranik23/avito-tech-spring/internal/repository"
	managermock "github.com/Ranik23/avito-tech-spring/internal/repository/manager/mock"
	repomock "github.com/Ranik23/avito-tech-spring/internal/repository/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreatePVZ_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := managermock.NewMockTxManager(ctrl)

	ctx := context.Background()
	city := "Moscow"
	fakePVZID := "pvz123"

	mockPVZRepo.EXPECT().CreatePVZ(gomock.Any(), city).Return(fakePVZID, nil).Times(1)
	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	).Times(1)

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, mockProductRepo, mockTxManager, slog.Default())

	pvzID, err := pvzService.CreatePVZ(ctx, city)

	assert.NoError(t, err)
	assert.Equal(t, fakePVZID, pvzID)
}

func TestCreatePVZ_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := managermock.NewMockTxManager(ctrl)

	ctx := context.Background()
	city := "Moscow"

	mockPVZRepo.EXPECT().CreatePVZ(gomock.Any(), city).Return("", repository.ErrAlreadyExists).Times(1)
	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	).Times(1)

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, mockProductRepo, mockTxManager, slog.Default())

	pvzID, err := pvzService.CreatePVZ(ctx, city)

	assert.ErrorIs(t, err, ErrAlreadyExists)
	assert.Empty(t, pvzID)
}


func TestPVZService_CloseReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := managermock.NewMockTxManager(ctrl)

	ctx := context.Background()
	pvzID := "pvz456"

	mockReceptionRepo.EXPECT().FindOpen(gomock.Any(), pvzID).Return(nil, nil)

	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, mockProductRepo, mockTxManager, slog.Default())

	id, err := pvzService.CloseReception(ctx, pvzID)
	assert.Equal(t, err, ErrNotFound) 
	assert.Equal(t, id, "")
}

func TestPVZService_DeleteLastProductFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := managermock.NewMockTxManager(ctrl)

	exampleCtx := context.Background()
	examplePvzID := "pvz456"

	exampleReception := domain.Reception{}

	mockReceptionRepo.EXPECT().FindOpen(gomock.Any(), examplePvzID).Return(&exampleReception, nil)
	mockProductRepo.EXPECT().FindTheLastProduct(gomock.Any(), examplePvzID).Return(nil, repository.ErrNotFound)

	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, mockProductRepo, mockTxManager, slog.Default())

	err := pvzService.DeleteLastProduct(exampleCtx, examplePvzID)
	assert.Equal(t, err, ErrEmpty) 
}


func TestPVZService_DeleteLastProductSucces(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := managermock.NewMockTxManager(ctrl)

	exampleCtx := context.Background()
	examplePvzID := "pvz456"
	exampleProduct := &domain.Product{
		ID: "123",
	}

	exampleReception := domain.Reception{}

	mockReceptionRepo.EXPECT().FindOpen(gomock.Any(), examplePvzID).Return(&exampleReception, nil)
	mockProductRepo.EXPECT().FindTheLastProduct(gomock.Any(), examplePvzID).Return(exampleProduct, nil)
	mockProductRepo.EXPECT().DeleteProduct(gomock.Any(), exampleProduct.ID).Return(nil)

	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, mockProductRepo, mockTxManager, slog.Default())

	err := pvzService.DeleteLastProduct(exampleCtx, examplePvzID)
	assert.NoError(t, err)
}

func TestPVZService_DeleteLastProductNoOpen(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := managermock.NewMockTxManager(ctrl)

	exampleCtx := context.Background()
	examplePvzID := "pvz456"

	mockReceptionRepo.EXPECT().FindOpen(gomock.Any(), examplePvzID).Return(nil, nil)

	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, mockProductRepo, mockTxManager, slog.Default())

	err := pvzService.DeleteLastProduct(exampleCtx, examplePvzID)
	assert.Equal(t, err, ErrAllReceptionsClosed)
}

func TestPVZService_StartReception_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := managermock.NewMockTxManager(ctrl)

	exampleCtx := context.Background()
	examplePvzID := "pvz456"
	exampleReception := &domain.Reception{
		ID: "rec789",
	}

	mockReceptionRepo.EXPECT().FindOpen(gomock.Any(), examplePvzID).Return(nil, nil)
	mockReceptionRepo.EXPECT().CreateReception(gomock.Any(), examplePvzID).Return(exampleReception.ID, nil)

	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, mockProductRepo, mockTxManager, slog.Default())

	receptionID, err := pvzService.StartReception(exampleCtx, examplePvzID)
	assert.NoError(t, err)
	assert.Equal(t, exampleReception.ID, receptionID)
}


func TestPVZService_StartReception_AlreadyOpen(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := managermock.NewMockTxManager(ctrl)

	exampleCtx := context.Background()
	examplePvzID := "pvz456"
	exampleReception := &domain.Reception{
		ID: "openRec123",
	}

	mockReceptionRepo.EXPECT().FindOpen(gomock.Any(), examplePvzID).Return(exampleReception, nil)

	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, mockProductRepo, mockTxManager, slog.Default())

	receptionID, err := pvzService.StartReception(exampleCtx, examplePvzID)
	assert.Error(t, err)
	assert.Equal(t, "", receptionID)
	assert.Equal(t, err, ErrAlreadyOpen)
}


func TestPVZService_AddProduct_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := managermock.NewMockTxManager(ctrl)

	exampleCtx := context.Background()
	examplePvzID := "pvz123"
	exampleProductType := "box"
	exampleReception := &domain.Reception{ID: "reception123"}
	expectedProductID := "product456"

	mockReceptionRepo.EXPECT().FindOpen(gomock.Any(), examplePvzID).Return(exampleReception, nil)
	mockProductRepo.EXPECT().CreateProduct(gomock.Any(), exampleProductType, exampleReception.ID).Return(expectedProductID, nil)

	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, mockProductRepo, mockTxManager, slog.Default())

	productID, err := pvzService.AddProduct(exampleCtx, examplePvzID, exampleProductType)
	assert.NoError(t, err)
	assert.Equal(t, expectedProductID, productID)
}



func TestPVZService_AddProduct_AllReceptionsClosed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := managermock.NewMockTxManager(ctrl)

	exampleCtx := context.Background()
	examplePvzID := "pvz789"
	exampleProductType := "box"

	mockReceptionRepo.EXPECT().FindOpen(gomock.Any(), examplePvzID).Return(nil, nil)

	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, mockProductRepo, mockTxManager, slog.Default())

	productID, err := pvzService.AddProduct(exampleCtx, examplePvzID, exampleProductType)
	assert.Error(t, err)
	assert.Equal(t, err, ErrAllReceptionsClosed)
	assert.Empty(t, productID)
}


func TestPVZService_AddProduct_AllReceptionsFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := managermock.NewMockTxManager(ctrl)

	exampleCtx := context.Background()
	examplePvzID := "pvz789"
	exampleProductType := "box"
	exampleReception := &domain.Reception{ID: "reception123"}

	mockReceptionRepo.EXPECT().FindOpen(gomock.Any(), examplePvzID).Return(exampleReception, nil).Times(1)
	mockProductRepo.EXPECT().CreateProduct(gomock.Any(), exampleProductType, exampleReception.ID).Return("", errors.New("fail"))

	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, mockProductRepo, mockTxManager, slog.Default())

	_, err := pvzService.AddProduct(exampleCtx, examplePvzID, exampleProductType)
	assert.Error(t, err)
}

func TestPVZService_GetPVZSInfo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := managermock.NewMockTxManager(ctrl)

	ctx := context.Background()
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()
	offset := 0
	limit := 10

	examplePvz := domain.Pvz{ID: "pvz1"}
	exampleReception := &domain.Reception{ID: "reception1"}
	exampleProduct := domain.Product{ID: "product1"}

	mockPVZRepo.EXPECT().GetPVZS(gomock.Any(), offset, limit).Return([]domain.Pvz{examplePvz}, nil)

	mockReceptionRepo.EXPECT().GetReceptionsFiltered(gomock.Any(), examplePvz.ID, start, end).Return([]*domain.Reception{exampleReception}, nil)

	mockProductRepo.EXPECT().GetProducts(gomock.Any(), exampleReception.ID).Return([]domain.Product{exampleProduct}, nil)

	mockTxManager.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, mockProductRepo, mockTxManager, slog.Default())

	pvzInfos, err := pvzService.GetPVZSInfo(ctx, start, end, offset, limit)
	assert.NoError(t, err)
	assert.Len(t, pvzInfos, 1)
	assert.Equal(t, examplePvz.ID, pvzInfos[0].Pvz.ID)
	assert.Len(t, pvzInfos[0].Receptions, 1)
	assert.Equal(t, exampleReception.ID, pvzInfos[0].Receptions[0].Reception.ID)
	assert.Len(t, pvzInfos[0].Receptions[0].Products, 1)
	assert.Equal(t, exampleProduct.ID, pvzInfos[0].Receptions[0].Products[0].ID)
}

func TestPVZService_GetPVZSInfo_PartialByOffsetLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := managermock.NewMockTxManager(ctrl)

	ctx := context.Background()
	start := time.Now().Add(-48 * time.Hour)
	end := time.Now()
	offset := 1
	limit := 1

	allPvzs := []domain.Pvz{
		{ID: "pvz1"},
		{ID: "pvz2"},
		{ID: "pvz3"},
	}

	selectedPvz := allPvzs[1] // pvz2
	exampleReception := &domain.Reception{ID: "reception2"}
	exampleProduct := domain.Product{ID: "product2"}

	mockPVZRepo.EXPECT().GetPVZS(gomock.Any(), offset, limit).Return([]domain.Pvz{selectedPvz}, nil)
	mockReceptionRepo.EXPECT().GetReceptionsFiltered(gomock.Any(), selectedPvz.ID, start, end).Return([]*domain.Reception{exampleReception}, nil)
	mockProductRepo.EXPECT().GetProducts(gomock.Any(), exampleReception.ID).Return([]domain.Product{exampleProduct}, nil)

	mockTxManager.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, mockProductRepo, mockTxManager, slog.Default())

	pvzInfos, err := pvzService.GetPVZSInfo(ctx, start, end, offset, limit)
	assert.NoError(t, err)
	assert.Len(t, pvzInfos, 1)
	assert.Equal(t, "pvz2", pvzInfos[0].Pvz.ID)
}

