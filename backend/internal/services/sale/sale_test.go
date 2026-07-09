package sale

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"supermarket/internal/http/dto"
	slogdiscard "supermarket/internal/lib/logger"
	testhelper "supermarket/internal/lib/test"
	productsalemock "supermarket/internal/mocks/product_sale"
	salemock "supermarket/internal/mocks/sale"
	stockmock "supermarket/internal/mocks/stock"
	transactionsnmock "supermarket/internal/mocks/transactions"
	"supermarket/internal/models"
	productsale "supermarket/internal/repository/product_sale"
	sale "supermarket/internal/repository/sale"
	"supermarket/internal/repository/stock"
	"supermarket/internal/repository/transactions"
)

var errRepository = errors.New("error in repository")

// newTestService creates new sale service for tests.
func newTestService(
	saleR *salemock.MockRepository,
	uow *transactionsnmock.MockUnitOfWork) Service {
	return New(
		slogdiscard.NewDiscardLogger(),
		saleR,
		uow,
	)
}

// TestGetSales tests method GetSales().
//
//nolint:funlen
func TestGetSales(t *testing.T) {
	type args struct {
		ctx      context.Context
		dateFrom *time.Time
		dateTo   *time.Time
	}

	saleModels := []models.Sale{
		newSaleModel(),
	}

	ctx := context.Background()

	testCases := []struct {
		name    string
		args    args
		mockFn  func(saleR *salemock.MockRepository, args args)
		wantLen int
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:      ctx,
				dateFrom: testhelper.Ptr(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
				dateTo:   testhelper.Ptr(time.Date(2027, 12, 31, 0, 0, 0, 0, time.UTC)),
			},
			mockFn: func(saleR *salemock.MockRepository, args args) {
				setupGetByParams(
					args.ctx,
					saleR,
					args.dateFrom,
					args.dateTo,
					saleModels,
					nil,
				)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "nil arguments",
			args: args{
				ctx:      ctx,
				dateFrom: nil,
				dateTo:   nil,
			},
			mockFn: func(saleR *salemock.MockRepository, args args) {
				setupGetByParams(
					args.ctx,
					saleR,
					args.dateFrom,
					args.dateTo,
					saleModels,
					nil,
				)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "error get in repository",
			args: args{
				ctx:      ctx,
				dateFrom: testhelper.Ptr(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
				dateTo:   testhelper.Ptr(time.Date(2027, 12, 31, 0, 0, 0, 0, time.UTC)),
			},
			mockFn: func(saleR *salemock.MockRepository, args args) {
				setupGetByParams(
					args.ctx,
					saleR,
					args.dateFrom,
					args.dateTo,
					saleModels,
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

			saleR := salemock.NewMockRepository(t)
			s := newTestService(saleR, nil)

			tc.mockFn(saleR, tc.args)
			got, err := s.GetSales(
				tc.args.ctx,
				tc.args.dateFrom,
				tc.args.dateTo,
			)

			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, got)

				return
			}

			require.NoError(t, err)
			require.Len(t, got, tc.wantLen)

			if len(got) > 0 {
				requireSaleDTO(t, saleModels[0], got[0])
			}
		})
	}
}

// TestCreateSale tests method CreateSale().
//
//nolint:funlen
func TestCreateSale(t *testing.T) {
	type args struct {
		ctx  context.Context
		sale *models.Sale
	}

	saleModel := newSaleModel()

	ctx := context.Background()

	testCases := []struct {
		name   string
		args   args
		mockFn func(
			saleR *salemock.MockRepository,
			args args,
		)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:  ctx,
				sale: &saleModel,
			},
			mockFn: func(
				saleR *salemock.MockRepository,
				args args,
			) {
				setupCreate(args.ctx, saleR, args.sale, nil)
			},
			wantErr: false,
		},
		{
			name: "error create sale in repository",
			args: args{
				ctx:  ctx,
				sale: &saleModel,
			},
			mockFn: func(
				saleR *salemock.MockRepository,
				args args,
			) {
				setupCreate(args.ctx, saleR, args.sale, errRepository)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			saleR := salemock.NewMockRepository(t)
			uow := transactionsnmock.NewMockUnitOfWork(t)
			s := newTestService(saleR, uow)

			tc.mockFn(saleR, tc.args)

			err := s.CreateSale(tc.args.ctx, tc.args.sale)

			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.NotEqual(t, uuid.Nil, tc.args.sale.ID)
		})
	}
}

// TestDeleteSale tests method DeleteSale().
//
//nolint:funlen
func TestDeleteSale(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}

	ctx := context.Background()

	testCases := []struct {
		name   string
		args   args
		mockFn func(
			uow *transactionsnmock.MockUnitOfWork,
			saleR *salemock.MockRepository,
			stockR *stockmock.MockRepository,
			productSaleR *productsalemock.MockRepository,
			args args,
		)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: ctx,
				id:  uuid.New(),
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				saleR *salemock.MockRepository,
				stockR *stockmock.MockRepository,
				productSaleR *productsalemock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, saleR, stockR, productSaleR, nil)
				setupDeleteBySale(args.ctx, productSaleR, args.id, nil)
				setupUpdateStockOnDeleteSale(args.ctx, stockR, args.id, nil)
				setupDelete(args.ctx, saleR, args.id, nil)
			},
			wantErr: false,
		},
		{
			name: "error delete product in sale",
			args: args{
				ctx: ctx,
				id:  uuid.New(),
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				saleR *salemock.MockRepository,
				stockR *stockmock.MockRepository,
				productSaleR *productsalemock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, saleR, stockR, productSaleR, nil)
				setupDeleteBySale(args.ctx, productSaleR, args.id, errRepository)
			},
			wantErr: true,
		},
		{
			name: "error update stock",
			args: args{
				ctx: ctx,
				id:  uuid.New(),
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				saleR *salemock.MockRepository,
				stockR *stockmock.MockRepository,
				productSaleR *productsalemock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, saleR, stockR, productSaleR, nil)
				setupDeleteBySale(args.ctx, productSaleR, args.id, nil)
				setupUpdateStockOnDeleteSale(args.ctx, stockR, args.id, errRepository)
			},
			wantErr: true,
		},
		{
			name: "error delete sale",
			args: args{
				ctx: ctx,
				id:  uuid.New(),
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				saleR *salemock.MockRepository,
				stockR *stockmock.MockRepository,
				productSaleR *productsalemock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, saleR, stockR, productSaleR, nil)
				setupDeleteBySale(args.ctx, productSaleR, args.id, nil)
				setupUpdateStockOnDeleteSale(args.ctx, stockR, args.id, nil)
				setupDelete(args.ctx, saleR, args.id, errRepository)
			},
			wantErr: true,
		},
		{
			name: "error uow",
			args: args{
				ctx: ctx,
				id:  uuid.New(),
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				saleR *salemock.MockRepository,
				stockR *stockmock.MockRepository,
				productSaleR *productsalemock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, saleR, stockR, productSaleR, errRepository)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			saleR := salemock.NewMockRepository(t)
			stockR := stockmock.NewMockRepository(t)
			uow := transactionsnmock.NewMockUnitOfWork(t)
			productSaleR := productsalemock.NewMockRepository(t)
			s := newTestService(saleR, uow)

			tc.mockFn(uow, saleR, stockR, productSaleR, tc.args)

			err := s.DeleteSale(tc.args.ctx, tc.args.id)

			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}

// TestUpdateSale tests method UpdateSale().
//
//nolint:funlen
func TestUpdateSale(t *testing.T) {
	type args struct {
		ctx  context.Context
		sale *models.Sale
	}

	saleModel := newSaleModel()

	ctx := context.Background()

	testCases := []struct {
		name   string
		args   args
		mockFn func(
			saleR *salemock.MockRepository,
			args args,
		)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:  ctx,
				sale: &saleModel,
			},
			mockFn: func(
				saleR *salemock.MockRepository,
				args args,
			) {
				setupUpdate(args.ctx, saleR, args.sale, nil)
			},
			wantErr: false,
		},
		{
			name: "error update in repository",
			args: args{
				ctx:  ctx,
				sale: &saleModel,
			},
			mockFn: func(
				saleR *salemock.MockRepository,
				args args,
			) {
				setupUpdate(args.ctx, saleR, args.sale, errRepository)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			saleR := salemock.NewMockRepository(t)
			uow := transactionsnmock.NewMockUnitOfWork(t)
			s := newTestService(saleR, uow)

			tc.mockFn(saleR, tc.args)

			err := s.UpdateSale(tc.args.ctx, tc.args.sale)

			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}

// setupUOW configures the expected behavior method Do() of a transaction.
func setupDoUOW(
	ctx context.Context,
	uow *transactionsnmock.MockUnitOfWork,
	saleR sale.Repository,
	stockR stock.Repository,
	productSaleR productsale.Repository,
	err error,
) {
	uow.EXPECT().
		Do(ctx, mock.Anything).
		RunAndReturn(func(
			ctx context.Context,
			fn func(context.Context, transactions.Repositories) error,
		) error {
			if err != nil {
				return err
			}

			repos := transactions.Repositories{
				Sale:        saleR,
				ProductSale: productSaleR,
				Stock:       stockR,
			}

			return fn(ctx, repos)
		})
}

// setupGetByParams configures the expected behavior method GetByParams() of a sale repository.
func setupGetByParams(
	ctx context.Context,
	saleR *salemock.MockRepository,
	dateFrom *time.Time,
	dateTo *time.Time,
	saleModels []models.Sale,
	err error,
) {
	saleR.EXPECT().
		GetByParams(ctx, dateFrom, dateTo).
		Return(saleModels, err).
		Once()
}

// setupCreate configures the expected behavior method Create() of a sale repository.
func setupCreate(
	ctx context.Context,
	saleR *salemock.MockRepository,
	sale *models.Sale,
	err error,
) {
	saleR.EXPECT().
		Create(ctx, sale).
		Return(err).
		Once()
}

// setupUpdate configures the expected behavior method Update() of a sale repository.
func setupUpdate(
	ctx context.Context,
	saleR *salemock.MockRepository,
	sale *models.Sale,
	err error,
) {
	saleR.EXPECT().
		Update(ctx, sale).
		Return(err).
		Once()
}

// setupDelete configures the expected behavior method Delete() of a sale repository.
func setupDelete(
	ctx context.Context,
	saleR *salemock.MockRepository,

	id uuid.UUID,
	err error,
) {
	saleR.EXPECT().
		Delete(ctx, id).
		Return(err).
		Once()
}

// setupDeleteBySale configures the expected behavior method DeleteBySale() of a stock repository.
func setupDeleteBySale(
	ctx context.Context,
	proudtSaleR *productsalemock.MockRepository,
	saleID uuid.UUID,
	err error,
) {
	proudtSaleR.EXPECT().
		DeleteBySale(
			ctx,
			saleID,
		).
		Return(err).
		Once()
}

// setupUpdateStockBySale sets the expected behavior method UpdateStockBySale() of a stock repository.
func setupUpdateStockOnDeleteSale(
	ctx context.Context,
	stockR *stockmock.MockRepository,
	saleID uuid.UUID,
	err error,
) {
	stockR.EXPECT().
		UpdateStockOnDeleteSale(ctx, saleID).
		Return(err).
		Once()
}

func requireSaleDTO(t *testing.T, want models.Sale, got dto.SaleDTO) {
	t.Helper()

	require.Equal(t, want.ID, got.ID)
	require.Equal(t, want.DateTime.Format("02.01.2006 15:04"), got.DateTime)
	require.Equal(t, want.Discount, got.Discount)
	require.Equal(t, want.FullCost, got.FullCost)
	require.Equal(t, want.TotalCost, got.TotalCost)
}

// newSaleModel creates and returns sale model with arbitrary values for tests.
func newSaleModel() models.Sale {
	return models.Sale{
		ID:        uuid.New(),
		DateTime:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Discount:  decimal.NewFromInt(10),
		FullCost:  decimal.NewFromInt(1000),
		TotalCost: decimal.NewFromInt(900),
	}
}
