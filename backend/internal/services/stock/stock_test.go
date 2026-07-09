package stock

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"supermarket/internal/http/dto"
	slogdiscard "supermarket/internal/lib/logger"
	testhelper "supermarket/internal/lib/test"
	stockmock "supermarket/internal/mocks/stock"
	"supermarket/internal/models"
	"supermarket/internal/repository"
	"supermarket/internal/services"
)

var errRepository = errors.New("error in repository")

// newTestService creates new stock service for tests.
func newTestService(stockR *stockmock.MockRepository) Service {
	return New(
		slogdiscard.NewDiscardLogger(),
		stockR,
	)
}

// TestGetStocks tests method GetStocks().
//
//nolint:funlen
func TestGetStocks(t *testing.T) {
	type args struct {
		ctx       context.Context
		search    string
		productID *uuid.UUID
	}

	stockModels := []models.Stock{
		newStockModel(),
	}

	testCases := []struct {
		name    string
		args    args
		mockFn  func(repo *stockmock.MockRepository, args args)
		wantLen int
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:       context.Background(),
				search:    "молоко",
				productID: testhelper.Ptr(uuid.New()),
			},
			mockFn: func(stockR *stockmock.MockRepository, args args) {
				setupGetByParams(
					args.ctx,
					stockR,
					args.search,
					args.productID,
					stockModels,
					nil,
				)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "nil arguments",
			args: args{
				ctx:       context.Background(),
				search:    "",
				productID: nil,
			},
			mockFn: func(stockR *stockmock.MockRepository, args args) {
				setupGetByParams(
					args.ctx,
					stockR,
					args.search,
					args.productID,
					stockModels,
					nil,
				)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "error get from repository",
			args: args{
				ctx:       context.Background(),
				search:    "молоко",
				productID: testhelper.Ptr(uuid.New()),
			},
			mockFn: func(stockR *stockmock.MockRepository, args args) {
				setupGetByParams(
					args.ctx,
					stockR,
					args.search,
					args.productID,
					nil,
					errRepository,
				)
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			stockR := stockmock.NewMockRepository(t)
			s := newTestService(stockR)

			tc.mockFn(stockR, tc.args)

			got, err := s.GetStocks(
				tc.args.ctx,
				tc.args.search,
				tc.args.productID,
			)

			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, got)

				return
			}

			require.NoError(t, err)
			require.Len(t, got, tc.wantLen)

			if len(got) > 0 {
				requireStockDTO(t, stockModels[0], got[0])
			}
		})
	}
}

// TestIncreaseStock tests method IncreaseStock().
//
//nolint:funlen
func TestIncreaseStock(t *testing.T) {
	type args struct {
		ctx       context.Context
		productID uuid.UUID
		count     decimal.Decimal
	}

	testCases := []struct {
		name    string
		args    args
		mockFn  func(stockR *stockmock.MockRepository, args args)
		wantErr error
	}{
		{
			name: "success",
			args: args{
				ctx:       context.Background(),
				productID: uuid.New(),
				count:     decimal.NewFromInt32(35),
			},
			mockFn: func(stockR *stockmock.MockRepository, args args) {
				setupFirstOrCreateByProductID(args.ctx, stockR, args.productID, nil)
				setupSetCountByProductID(args.ctx, stockR, args.productID, args.count, nil)
			},
			wantErr: nil,
		},
		{
			name: "error get or create in repository",
			args: args{
				ctx:       context.Background(),
				productID: uuid.New(),
				count:     decimal.NewFromInt32(35),
			},
			mockFn: func(stockR *stockmock.MockRepository, args args) {
				setupFirstOrCreateByProductID(args.ctx, stockR, args.productID, errRepository)
			},
			wantErr: errRepository,
		},
		{
			name: "error set in repository",
			args: args{
				ctx:       context.Background(),
				productID: uuid.New(),
				count:     decimal.NewFromInt32(35),
			},
			mockFn: func(stockR *stockmock.MockRepository, args args) {
				setupFirstOrCreateByProductID(args.ctx, stockR, args.productID, nil)
				setupSetCountByProductID(args.ctx, stockR, args.productID, args.count, errRepository)
			},
			wantErr: errRepository,
		},
		{
			name: "error not enough in repository",
			args: args{
				ctx:       context.Background(),
				productID: uuid.New(),
				count:     decimal.NewFromInt32(35),
			},
			mockFn: func(stockR *stockmock.MockRepository, args args) {
				setupFirstOrCreateByProductID(args.ctx, stockR, args.productID, nil)
				setupSetCountByProductID(args.ctx, stockR, args.productID, args.count, repository.ErrNotEnoughStock)
			},
			wantErr: services.ErrNotEnoughStock,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			stockR := stockmock.NewMockRepository(t)
			s := newTestService(stockR)

			tc.mockFn(stockR, tc.args)

			err := s.IncreaseStock(tc.args.ctx, tc.args.productID, tc.args.count)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)

				return
			}

			require.NoError(t, err)
		})
	}
}

// setupGetByParams configures the expected behavior method GetByParams() of a productSupply repository.
func setupGetByParams(
	ctx context.Context,
	stockR *stockmock.MockRepository,
	search string,
	productID *uuid.UUID,
	stockModels []models.Stock,
	err error,
) {
	stockR.EXPECT().
		GetByParams(
			ctx,
			search,
			productID,
		).
		Return(stockModels, err).
		Once()
}

// setupSetCountByProductID configures the expected behavior method SetCountByProductID() of a stock repository.
func setupSetCountByProductID(
	ctx context.Context,
	stockR *stockmock.MockRepository,
	productID uuid.UUID,
	count decimal.Decimal,
	err error,
) {
	stockR.EXPECT().
		SetCountByProductID(
			ctx,
			productID,
			count,
		).
		Return(err).
		Once()
}

// setupFirstOrCreateByProductID sets the expected behavior method FirstOrCreateByProductID() of a stock repository.
func setupFirstOrCreateByProductID(
	ctx context.Context,
	stockR *stockmock.MockRepository,
	productID uuid.UUID,
	err error,
) {
	stockR.EXPECT().
		FirstOrCreateByProductID(ctx, productID).
		Return(&models.Stock{}, err).
		Once()
}

func requireStockDTO(t *testing.T, want models.Stock, got dto.StockDTO) {
	t.Helper()

	require.Equal(t, want.ID, got.ID)
	require.Equal(t, want.Quantity, got.Quantity)
}

// newStockModel creates and returns stock model with arbitrary values for tests.
func newStockModel() models.Stock {
	return models.Stock{
		ID:        uuid.New(),
		ProductID: testhelper.Ptr(uuid.New()),
		Quantity:  decimal.NewFromInt(10),
	}
}
