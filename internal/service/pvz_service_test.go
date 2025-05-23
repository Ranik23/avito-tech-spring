//go:build unit

package service

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
	"github.com/Ranik23/avito-tech-spring/internal/repository"
	repomock "github.com/Ranik23/avito-tech-spring/internal/repository/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreatePVZ_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := repomock.NewMockTxManager(ctrl)
	cities := []string{"Moscow"}

	ctx := context.Background()
	city := "Moscow"
	fakePVZ := &domain.Pvz{
		ID:   "fakeID",
		City: "Moscow",
	}

	mockPVZRepo.EXPECT().CreatePVZ(gomock.Any(), city).Return(fakePVZ, nil).Times(1)
	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	).Times(1)

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, cities, mockProductRepo, mockTxManager, slog.Default())

	pvz, err := pvzService.CreatePVZ(ctx, city)

	assert.NoError(t, err)
	assert.Equal(t, pvz.City, "Moscow")
	assert.Equal(t, pvz.ID, "fakeID")
}

func TestPVZService_StartReception_CreateReceptionFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := repomock.NewMockTxManager(ctrl)

	exampleCtx := context.Background()
	examplePvzID := "pvz456"

	mockReceptionRepo.EXPECT().FindOpen(gomock.Any(), examplePvzID).Return(nil, nil)
	mockReceptionRepo.EXPECT().CreateReception(gomock.Any(), examplePvzID).Return(nil, errors.New("db error"))

	mockTxManager.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, []string{"Moscow"}, mockProductRepo, mockTxManager, slog.Default())

	_, err := pvzService.StartReception(exampleCtx, examplePvzID)
	assert.Error(t, err)
}

func TestCreatePVZ_UnknownCity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := repomock.NewMockTxManager(ctrl)

	cities := []string{"Moscow"}

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, cities, mockProductRepo, mockTxManager, slog.Default())

	_, err := pvzService.CreatePVZ(context.Background(), "London")

	assert.ErrorIs(t, err, ErrInvalidCity)
}

func TestCreatePVZ_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := repomock.NewMockTxManager(ctrl)
	cities := []string{"Moscow"}

	ctx := context.Background()
	city := "Moscow"

	mockPVZRepo.EXPECT().CreatePVZ(gomock.Any(), city).Return(nil, repository.ErrAlreadyExists).Times(1)
	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	).Times(1)

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, cities, mockProductRepo, mockTxManager, slog.Default())

	_, err := pvzService.CreatePVZ(ctx, city)

	assert.ErrorIs(t, err, ErrAlreadyExists)
}

func TestPVZService_CloseReceptionFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := repomock.NewMockTxManager(ctrl)
	cities := []string{"Moscow"}

	ctx := context.Background()
	pvzID := "pvz456"

	mockReceptionRepo.EXPECT().FindOpen(gomock.Any(), pvzID).Return(nil, nil)

	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, cities, mockProductRepo, mockTxManager, slog.Default())

	_, err := pvzService.CloseReception(ctx, pvzID)
	assert.Equal(t, err, ErrAllReceptionsClosed)
}

func TestPVZService_CloseReceptionSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := repomock.NewMockTxManager(ctrl)
	cities := []string{"Moscow"}

	ctx := context.Background()
	examplepvzID := "pvz456"
	exampleReception := &domain.Reception{
		ID: "1",
		Status: "open",
		PvzID: examplepvzID,
	}

	mockReceptionRepo.EXPECT().FindOpen(gomock.Any(), examplepvzID).Return(exampleReception, nil)
	mockReceptionRepo.EXPECT().UpdateReceptionStatus(gomock.Any(), exampleReception.ID, "closed").Return(nil)

	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, cities, mockProductRepo, mockTxManager, slog.Default())

	reception, err := pvzService.CloseReception(ctx, examplepvzID)
	require.NoError(t, err)
	require.Equal(t, reception.ID, "1")
	require.Equal(t, reception.PvzID, examplepvzID)
	require.Equal(t, reception.Status, "open")
}

func TestPVZService_DeleteLastProductFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := repomock.NewMockTxManager(ctrl)
	cities := []string{"Moscow"}

	exampleCtx := context.Background()
	examplePvzID := "pvz456"

	exampleReception := domain.Reception{}

	mockReceptionRepo.EXPECT().FindOpen(gomock.Any(), examplePvzID).Return(&exampleReception, nil)
	mockProductRepo.EXPECT().FindTheLastProduct(gomock.Any(), examplePvzID).Return(nil, nil)

	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, cities, mockProductRepo, mockTxManager, slog.Default())

	err := pvzService.DeleteLastProduct(exampleCtx, examplePvzID)
	assert.Equal(t, err, ErrReceptionEmpty)
}

func TestPVZService_DeleteLastProductSucces(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := repomock.NewMockTxManager(ctrl)
	cities := []string{"Moscow"}

	exampleCtx := context.Background()
	examplePvzID := "pvz456"
	exampleProduct := &domain.Product{
		ID: "123",
	}

	exampleReception := domain.Reception{}

	mockReceptionRepo.EXPECT().FindOpen(exampleCtx, examplePvzID).Return(&exampleReception, nil)
	mockProductRepo.EXPECT().FindTheLastProduct(exampleCtx, examplePvzID).Return(exampleProduct, nil)
	mockProductRepo.EXPECT().DeleteProduct(exampleCtx, exampleProduct.ID).Return(nil)

	mockTxManager.EXPECT().Do(exampleCtx, gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, cities, mockProductRepo, mockTxManager, slog.Default())

	err := pvzService.DeleteLastProduct(exampleCtx, examplePvzID)
	assert.NoError(t, err)
}

func TestPVZService_DeleteLastProductNoOpen(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := repomock.NewMockTxManager(ctrl)
	cities := []string{"Moscow"}

	exampleCtx := context.Background()
	examplePvzID := "pvz456"

	mockReceptionRepo.EXPECT().FindOpen(gomock.Any(), examplePvzID).Return(nil, nil)

	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, cities, mockProductRepo, mockTxManager, slog.Default())

	err := pvzService.DeleteLastProduct(exampleCtx, examplePvzID)
	assert.Equal(t, err, ErrAllReceptionsClosed)
}

func TestPVZService_StartReception_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := repomock.NewMockTxManager(ctrl)
	cities := []string{"Moscow"}

	exampleCtx := context.Background()
	examplePvzID := "pvz456"
	exampleReception := &domain.Reception{
		ID: "rec789",
	}

	mockReceptionRepo.EXPECT().FindOpen(gomock.Any(), examplePvzID).Return(nil, nil)
	mockReceptionRepo.EXPECT().CreateReception(gomock.Any(), examplePvzID).Return(exampleReception, nil)

	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, cities, mockProductRepo, mockTxManager, slog.Default())

	reception, err := pvzService.StartReception(exampleCtx, examplePvzID)
	assert.NoError(t, err)
	assert.Equal(t, exampleReception.ID, reception.ID)
}

func TestPVZService_StartReception_AlreadyOpen(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := repomock.NewMockTxManager(ctrl)
	cities := []string{"Moscow"}

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

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, cities, mockProductRepo, mockTxManager, slog.Default())

	_, err := pvzService.StartReception(exampleCtx, examplePvzID)
	assert.Equal(t, err, ErrAlreadyOpen)
}

func TestPVZService_AddProduct_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := repomock.NewMockTxManager(ctrl)
	cities := []string{"Moscow"}

	exampleCtx := context.Background()
	examplePvzID := "pvz123"
	exampleProductType := "box"
	exampleReception := &domain.Reception{ID: "reception123"}
	expectedProduct := &domain.Product{
		ID: "id",
	}

	mockReceptionRepo.EXPECT().FindOpen(gomock.Any(), examplePvzID).Return(exampleReception, nil)
	mockProductRepo.EXPECT().CreateProduct(gomock.Any(), exampleProductType, exampleReception.ID).Return(expectedProduct, nil)

	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, cities, mockProductRepo, mockTxManager, slog.Default())

	product, err := pvzService.AddProduct(exampleCtx, examplePvzID, exampleProductType)
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct.ID, product.ID)
}

func TestPVZService_AddProduct_AllReceptionsClosed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := repomock.NewMockTxManager(ctrl)
	cities := []string{"Moscow"}

	exampleCtx := context.Background()
	examplePvzID := "pvz789"
	exampleProductType := "box"

	mockReceptionRepo.EXPECT().FindOpen(gomock.Any(), examplePvzID).Return(nil, nil)

	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, cities, mockProductRepo, mockTxManager, slog.Default())

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
	mockTxManager := repomock.NewMockTxManager(ctrl)
	cities := []string{"Moscow"}

	exampleCtx := context.Background()
	examplePvzID := "pvz789"
	exampleProductType := "box"
	exampleReception := &domain.Reception{ID: "reception123"}

	mockReceptionRepo.EXPECT().FindOpen(gomock.Any(), examplePvzID).Return(exampleReception, nil).Times(1)
	mockProductRepo.EXPECT().CreateProduct(gomock.Any(), exampleProductType, exampleReception.ID).Return(nil, errors.New("fail"))

	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, cities, mockProductRepo, mockTxManager, slog.Default())

	_, err := pvzService.AddProduct(exampleCtx, examplePvzID, exampleProductType)
	assert.Error(t, err)
}

func TestPVZService_GetPVZSInfo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := repomock.NewMockTxManager(ctrl)
	cities := []string{"Moscow"}

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

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, cities, mockProductRepo, mockTxManager, slog.Default())

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
	mockTxManager := repomock.NewMockTxManager(ctrl)
	cities := []string{"Moscow"}

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

	selectedPvz := allPvzs[1]
	exampleReception := &domain.Reception{ID: "reception2"}
	exampleProduct := domain.Product{ID: "product2"}

	mockPVZRepo.EXPECT().GetPVZS(gomock.Any(), offset, limit).Return([]domain.Pvz{selectedPvz}, nil)
	mockReceptionRepo.EXPECT().GetReceptionsFiltered(gomock.Any(), selectedPvz.ID, start, end).Return([]*domain.Reception{exampleReception}, nil)
	mockProductRepo.EXPECT().GetProducts(gomock.Any(), exampleReception.ID).Return([]domain.Product{exampleProduct}, nil)

	mockTxManager.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, cities, mockProductRepo, mockTxManager, slog.Default())

	pvzInfos, err := pvzService.GetPVZSInfo(ctx, start, end, offset, limit)
	assert.NoError(t, err)
	assert.Len(t, pvzInfos, 1)
	assert.Equal(t, "pvz2", pvzInfos[0].Pvz.ID)
}

func TestPVZService_GetPVZSInfo_ProductFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := repomock.NewMockTxManager(ctrl)

	ctx := context.Background()
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()
	offset := 0
	limit := 1

	examplePvz := domain.Pvz{ID: "pvz1"}
	exampleReception := &domain.Reception{ID: "reception1"}

	mockPVZRepo.EXPECT().GetPVZS(gomock.Any(), offset, limit).Return([]domain.Pvz{examplePvz}, nil)
	mockReceptionRepo.EXPECT().GetReceptionsFiltered(gomock.Any(), examplePvz.ID, start, end).Return([]*domain.Reception{exampleReception}, nil)
	mockProductRepo.EXPECT().GetProducts(gomock.Any(), exampleReception.ID).Return(nil, errors.New("db error"))

	mockTxManager.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, []string{"Moscow"}, mockProductRepo, mockTxManager, slog.Default())

	_, err := pvzService.GetPVZSInfo(ctx, start, end, offset, limit)
	assert.Error(t, err)
}

func TestPVZService_GetPVZList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := repomock.NewMockPvzRepository(ctrl)
	mockReceptionRepo := repomock.NewMockReceptionRepository(ctrl)
	mockProductRepo := repomock.NewMockProductRepository(ctrl)
	mockTxManager := repomock.NewMockTxManager(ctrl)
	cities := []string{"Moscow"}

	examplePVZS := []domain.Pvz{
		{ID: "1"},
		{ID: "2"},
	}

	mockPVZRepo.EXPECT().GetListOfPVZS(gomock.Any()).Return(examplePVZS, nil).Times(1)

	pvzService := NewPVZService(mockPVZRepo, mockReceptionRepo, cities, mockProductRepo, mockTxManager, slog.Default())

	pvzs, err := pvzService.GetPVZList(context.Background())
	require.NoError(t, err)

	require.Equal(t, len(examplePVZS), len(pvzs))
	require.Equal(t, examplePVZS[0].ID, pvzs[0].ID)
	require.Equal(t, examplePVZS[1].ID, pvzs[1].ID)

}
