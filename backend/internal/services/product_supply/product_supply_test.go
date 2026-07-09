package productsupply

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
	productSupplymock "supermarket/internal/mocks/product_supply"
	stockmock "supermarket/internal/mocks/stock"
	transactionsnmock "supermarket/internal/mocks/transactions"
	"supermarket/internal/models"
	productsupply "supermarket/internal/repository/product_supply"
	"supermarket/internal/repository/stock"
	"supermarket/internal/repository/transactions"
)

var errRepository = errors.New("error in repository")

// newTestService creates new productSupply service for tests.
func newTestService(
	productSupplyR *productSupplymock.MockRepository,
	uow *transactionsnmock.MockUnitOfWork) Service {
	return New(
		slogdiscard.NewDiscardLogger(),
		productSupplyR,
		uow,
	)
}

// TestGetProductSupplies tests method GetProductSupplies().
//
//nolint:funlen
func TestGetProductSupplies(t *testing.T) {
	type args struct {
		ctx        context.Context
		productID  *uuid.UUID
		supplierID *uuid.UUID
		dateFrom   *time.Time
		dateTo     *time.Time
	}

	productSupplyModels := []models.ProductSupply{
		newProductSupplyModel(),
	}

	ctx := context.Background()

	testCases := []struct {
		name    string
		args    args
		mockFn  func(productSupplyR *productSupplymock.MockRepository, args args)
		wantLen int
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:        ctx,
				productID:  testhelper.Ptr(uuid.New()),
				supplierID: testhelper.Ptr(uuid.New()),
				dateFrom:   testhelper.Ptr(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
				dateTo:     testhelper.Ptr(time.Date(2027, 12, 31, 0, 0, 0, 0, time.UTC)),
			},
			mockFn: func(productSupplyR *productSupplymock.MockRepository, args args) {
				setupGetByParams(
					args.ctx,
					productSupplyR,
					args.productID,
					args.supplierID,
					args.dateFrom,
					args.dateTo,
					productSupplyModels,
					nil,
				)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "nil arguments",
			args: args{
				ctx:        ctx,
				productID:  nil,
				supplierID: nil,
				dateFrom:   nil,
				dateTo:     nil,
			},
			mockFn: func(productSupplyR *productSupplymock.MockRepository, args args) {
				setupGetByParams(
					args.ctx,
					productSupplyR,
					args.productID,
					args.supplierID,
					args.dateFrom,
					args.dateTo,
					productSupplyModels,
					nil,
				)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "error get in repository",
			args: args{
				ctx:        ctx,
				productID:  testhelper.Ptr(uuid.New()),
				supplierID: testhelper.Ptr(uuid.New()),
				dateFrom:   testhelper.Ptr(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
				dateTo:     testhelper.Ptr(time.Date(2027, 12, 31, 0, 0, 0, 0, time.UTC)),
			},
			mockFn: func(productSupplyR *productSupplymock.MockRepository, args args) {
				setupGetByParams(
					args.ctx,
					productSupplyR,
					args.productID,
					args.supplierID,
					args.dateFrom,
					args.dateTo,
					productSupplyModels,
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

			productSupplyR := productSupplymock.NewMockRepository(t)
			s := newTestService(productSupplyR, nil)

			tc.mockFn(productSupplyR, tc.args)
			got, err := s.GetProductSupplies(
				tc.args.ctx,
				tc.args.productID,
				tc.args.supplierID,
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
				requireProductSupplyDTO(t, productSupplyModels[0], got[0])
			}
		})
	}
}

// TestCreateProductSupply tests method CreateProductSupply().
//
//nolint:funlen
func TestCreateProductSupply(t *testing.T) {
	type args struct {
		ctx           context.Context
		productSupply *models.ProductSupply
	}

	productSupplyModel := newProductSupplyModel()

	ctx := context.Background()

	testCases := []struct {
		name   string
		args   args
		mockFn func(
			uow *transactionsnmock.MockUnitOfWork,
			productSupplyR *productSupplymock.MockRepository,
			stockR *stockmock.MockRepository,
			args args,
		)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:           ctx,
				productSupply: &productSupplyModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				productSupplyR *productSupplymock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSupplyR, stockR, nil)
				setupCreate(args.ctx, productSupplyR, args.productSupply, nil)
				setupFirstOrCreateByProductID(args.ctx, stockR, args.productSupply.ProductID, nil)
				setupSetCountByProductID(args.ctx, stockR, args.productSupply, nil)
			},
			wantErr: false,
		},
		{
			name: "error create productSupply",
			args: args{
				ctx:           ctx,
				productSupply: &productSupplyModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				productSupplyR *productSupplymock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSupplyR, stockR, nil)
				setupCreate(args.ctx, productSupplyR, args.productSupply, errRepository)
			},
			wantErr: true,
		},
		{
			name: "error get stock for set stock",
			args: args{
				ctx:           ctx,
				productSupply: &productSupplyModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				productSupplyR *productSupplymock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSupplyR, stockR, nil)
				setupCreate(args.ctx, productSupplyR, args.productSupply, nil)
				setupFirstOrCreateByProductID(args.ctx, stockR, args.productSupply.ProductID, errRepository)
			},
			wantErr: true,
		},
		{
			name: "error set stock",
			args: args{
				ctx:           ctx,
				productSupply: &productSupplyModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				productSupplyR *productSupplymock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSupplyR, stockR, nil)
				setupCreate(args.ctx, productSupplyR, args.productSupply, nil)
				setupFirstOrCreateByProductID(args.ctx, stockR, args.productSupply.ProductID, nil)
				setupSetCountByProductID(args.ctx, stockR, args.productSupply, errRepository)
			},
			wantErr: true,
		},
		{
			name: "error uow",
			args: args{
				ctx:           ctx,
				productSupply: &productSupplyModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				productSupplyR *productSupplymock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSupplyR, stockR, errRepository)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			productSupplyR := productSupplymock.NewMockRepository(t)
			stockR := stockmock.NewMockRepository(t)
			uow := transactionsnmock.NewMockUnitOfWork(t)
			s := newTestService(productSupplyR, uow)

			tc.mockFn(uow, productSupplyR, stockR, tc.args)

			err := s.CreateProductSupply(tc.args.ctx, tc.args.productSupply)

			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.NotEqual(t, uuid.Nil, tc.args.productSupply.ID)
		})
	}
}

// TestDeleteProductSupply tests method DeleteProductSupply().
//
//nolint:funlen
func TestDeleteProductSupply(t *testing.T) {
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
			productSupplyR *productSupplymock.MockRepository,
			stockR *stockmock.MockRepository,
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
				productSupplyR *productSupplymock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSupplyR, stockR, nil)
				setupUpdateStockByProductSupply(args.ctx, stockR, args.id, decimal.Zero, nil)
				setupDelete(args.ctx, productSupplyR, args.id, nil)
			},
			wantErr: false,
		},
		{
			name: "error set stock",
			args: args{
				ctx: ctx,
				id:  uuid.New(),
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				productSupplyR *productSupplymock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSupplyR, stockR, nil)
				setupUpdateStockByProductSupply(args.ctx, stockR, args.id, decimal.Zero, errRepository)
			},
			wantErr: true,
		},
		{
			name: "error delete productSupply",
			args: args{
				ctx: ctx,
				id:  uuid.New(),
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				productSupplyR *productSupplymock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSupplyR, stockR, nil)
				setupUpdateStockByProductSupply(args.ctx, stockR, args.id, decimal.Zero, nil)
				setupDelete(args.ctx, productSupplyR, args.id, errRepository)
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
				productSupplyR *productSupplymock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSupplyR, stockR, errRepository)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			productSupplyR := productSupplymock.NewMockRepository(t)
			stockR := stockmock.NewMockRepository(t)
			uow := transactionsnmock.NewMockUnitOfWork(t)
			s := newTestService(productSupplyR, uow)

			tc.mockFn(uow, productSupplyR, stockR, tc.args)

			err := s.DeleteProductSupply(tc.args.ctx, tc.args.id)

			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}

// TestUpdateProductSupply tests method UpdateProductSupply().
//
//nolint:funlen
func TestUpdateProductSupply(t *testing.T) {
	type args struct {
		ctx           context.Context
		productSupply *models.ProductSupply
	}

	productSupplyModel := newProductSupplyModel()

	ctx := context.Background()

	testCases := []struct {
		name   string
		args   args
		mockFn func(
			uow *transactionsnmock.MockUnitOfWork,
			productSupplyR *productSupplymock.MockRepository,
			stockR *stockmock.MockRepository,
			args args,
		)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:           ctx,
				productSupply: &productSupplyModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				productSupplyR *productSupplymock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSupplyR, stockR, nil)
				setupUpdateStockByProductSupply(
					args.ctx,
					stockR,
					args.productSupply.ID,
					productSupplyModel.Quantity,
					nil,
				)
				setupUpdate(args.ctx, productSupplyR, args.productSupply, nil)
			},
			wantErr: false,
		},
		{
			name: "error set stock",
			args: args{
				ctx:           ctx,
				productSupply: &productSupplyModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				productSupplyR *productSupplymock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSupplyR, stockR, nil)
				setupUpdateStockByProductSupply(
					args.ctx,
					stockR,
					args.productSupply.ID,
					productSupplyModel.Quantity,
					errRepository,
				)
			},
			wantErr: true,
		},
		{
			name: "error update productSupply",
			args: args{
				ctx:           ctx,
				productSupply: &productSupplyModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				productSupplyR *productSupplymock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSupplyR, stockR, nil)
				setupUpdateStockByProductSupply(
					args.ctx,
					stockR,
					args.productSupply.ID,
					productSupplyModel.Quantity,
					nil,
				)
				setupUpdate(args.ctx, productSupplyR, args.productSupply, errRepository)
			},
			wantErr: true,
		},
		{
			name: "error uow",
			args: args{
				ctx:           ctx,
				productSupply: &productSupplyModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				productSupplyR *productSupplymock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSupplyR, stockR, errRepository)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			productSupplyR := productSupplymock.NewMockRepository(t)
			stockR := stockmock.NewMockRepository(t)
			uow := transactionsnmock.NewMockUnitOfWork(t)
			s := newTestService(productSupplyR, uow)

			tc.mockFn(uow, productSupplyR, stockR, tc.args)

			err := s.UpdateProductSupply(tc.args.ctx, tc.args.productSupply)

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
	productSupplyR productsupply.Repository,
	stockR stock.Repository,
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
				ProductSupply: productSupplyR,
				Stock:         stockR,
			}

			return fn(ctx, repos)
		})
}

// setupGetByParams configures the expected behavior method GetByParams() of a productSupply repository.
func setupGetByParams(
	ctx context.Context,
	productSupplyR *productSupplymock.MockRepository,
	productID *uuid.UUID,
	supplierID *uuid.UUID,
	dateFrom *time.Time,
	dateTo *time.Time,
	productSupplyModels []models.ProductSupply,
	err error,
) {
	productSupplyR.EXPECT().
		GetByParams(ctx, productID, supplierID, dateFrom, dateTo).
		Return(productSupplyModels, err).
		Once()
}

// setupCreate configures the expected behavior method Create() of a productSupply repository.
func setupCreate(
	ctx context.Context,
	productSupplyR *productSupplymock.MockRepository,
	productSupply *models.ProductSupply,
	err error,
) {
	productSupplyR.EXPECT().
		Create(ctx, productSupply).
		Return(err).
		Once()
}

// setupUpdate configures the expected behavior method Update() of a productSupply repository.
func setupUpdate(
	ctx context.Context,
	productSupplyR *productSupplymock.MockRepository,
	productSupply *models.ProductSupply,
	err error,
) {
	productSupplyR.EXPECT().
		Update(ctx, productSupply).
		Return(err).
		Once()
}

// setupDelete configures the expected behavior method Delete() of a productSupply repository.
func setupDelete(
	ctx context.Context,
	productSupplyR *productSupplymock.MockRepository,

	id uuid.UUID,
	err error,
) {
	productSupplyR.EXPECT().
		Delete(ctx, id).
		Return(err).
		Once()
}

// setupSetCountByProductID configures the expected behavior method SetCountByProductID() of a stock repository.
func setupSetCountByProductID(
	ctx context.Context,
	stockR *stockmock.MockRepository,
	productSupply *models.ProductSupply,
	err error,
) {
	stockR.EXPECT().
		SetCountByProductID(
			ctx,
			productSupply.ProductID,
			productSupply.Quantity,
		).
		Return(err).
		Once()
}

// setupUpdateStockByProductSupply sets the expected behavior method UpdateStockByProductSupply() of a stock repository.
func setupUpdateStockByProductSupply(
	ctx context.Context,
	stockR *stockmock.MockRepository,
	productSupplyID uuid.UUID,
	quantity decimal.Decimal,
	err error,
) {
	stockR.EXPECT().
		UpdateStockByProductSupply(ctx, productSupplyID, quantity).
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

func requireProductSupplyDTO(t *testing.T, want models.ProductSupply, got dto.ProductSupplyDTO) {
	t.Helper()

	require.Equal(t, want.ID, got.ID)
	require.Equal(t, want.ProductID, got.ProductID)
	require.Equal(t, want.DeliveryDate.Format("02.01.2006"), got.DeliveryDate)
	require.Equal(t, want.Price, got.Price)
	require.Equal(t, want.Quantity, got.Quantity)
}

// newProductSupplyModel creates and returns productSupply model with arbitrary values for tests.
func newProductSupplyModel() models.ProductSupply {
	return models.ProductSupply{
		ID:           uuid.New(),
		ProductID:    uuid.New(),
		SupplierID:   uuid.New(),
		Price:        decimal.NewFromInt(200),
		Quantity:     decimal.NewFromInt(10),
		DeliveryDate: testhelper.Ptr(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
	}
}
