package productsale

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"supermarket/internal/http/dto"
	slogdiscard "supermarket/internal/lib/logger"
	productsalemock "supermarket/internal/mocks/product_sale"
	stockmock "supermarket/internal/mocks/stock"
	transactionsnmock "supermarket/internal/mocks/transactions"
	"supermarket/internal/models"
	productsale "supermarket/internal/repository/product_sale"
	"supermarket/internal/repository/stock"
	"supermarket/internal/repository/transactions"
)

var errRepository = errors.New("error in repository")

// newTestService creates new productSale service for tests.
func newTestService(
	productSaleR *productsalemock.MockRepository,
	uow *transactionsnmock.MockUnitOfWork) Service {
	return New(
		slogdiscard.NewDiscardLogger(),
		productSaleR,
		uow,
	)
}

// TestGetProductsInSales tests method GetProductsInSales().
//
//nolint:funlen
func TestGetProductsInSales(t *testing.T) {
	type args struct {
		ctx    context.Context
		saleID uuid.UUID
	}

	productSaleModels := []models.ProductSale{
		newProductSaleModel(),
	}

	ctx := context.Background()

	testCases := []struct {
		name    string
		args    args
		mockFn  func(productSaleR *productsalemock.MockRepository, args args)
		wantLen int
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:    ctx,
				saleID: uuid.New(),
			},
			mockFn: func(productSaleR *productsalemock.MockRepository, args args) {
				setupGetBySale(
					args.ctx,
					productSaleR,
					args.saleID,
					productSaleModels,
					nil,
				)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "error get in repository",
			args: args{
				ctx:    ctx,
				saleID: uuid.New(),
			},
			mockFn: func(productSaleR *productsalemock.MockRepository, args args) {
				setupGetBySale(
					args.ctx,
					productSaleR,
					args.saleID,
					productSaleModels,
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

			productSaleR := productsalemock.NewMockRepository(t)
			s := newTestService(productSaleR, nil)

			tc.mockFn(productSaleR, tc.args)
			got, err := s.GetProductsInSale(
				tc.args.ctx,
				tc.args.saleID,
			)

			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, got)

				return
			}

			require.NoError(t, err)
			require.Len(t, got, tc.wantLen)

			if len(got) > 0 {
				requireProductSaleDTO(t, productSaleModels[0], got[0])
			}
		})
	}
}

// TestCreateProductSale tests method CreateProductSale().
//
//nolint:funlen
func TestCreateProductSale(t *testing.T) {
	type args struct {
		ctx         context.Context
		productSale *models.ProductSale
	}

	productSaleModel := newProductSaleModel()

	ctx := context.Background()

	testCases := []struct {
		name   string
		args   args
		mockFn func(
			uow *transactionsnmock.MockUnitOfWork,
			productSaleR *productsalemock.MockRepository,
			stockR *stockmock.MockRepository,
			args args,
		)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:         ctx,
				productSale: &productSaleModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				productSaleR *productsalemock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSaleR, stockR, nil)
				setupCreate(args.ctx, productSaleR, args.productSale, nil)
				setupFirstOrCreateByProductID(args.ctx, stockR, args.productSale.ProductID, nil)
				setupSetCountByProductID(args.ctx, stockR, args.productSale, nil)
			},
			wantErr: false,
		},
		{
			name: "error create productSale",
			args: args{
				ctx:         ctx,
				productSale: &productSaleModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				productSaleR *productsalemock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSaleR, stockR, nil)
				setupCreate(args.ctx, productSaleR, args.productSale, errRepository)
			},
			wantErr: true,
		},
		{
			name: "error get stock for set stock",
			args: args{
				ctx:         ctx,
				productSale: &productSaleModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				productSaleR *productsalemock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSaleR, stockR, nil)
				setupCreate(args.ctx, productSaleR, args.productSale, nil)
				setupFirstOrCreateByProductID(args.ctx, stockR, args.productSale.ProductID, errRepository)
			},
			wantErr: true,
		},
		{
			name: "error set stock",
			args: args{
				ctx:         ctx,
				productSale: &productSaleModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				productSaleR *productsalemock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSaleR, stockR, nil)
				setupCreate(args.ctx, productSaleR, args.productSale, nil)
				setupFirstOrCreateByProductID(args.ctx, stockR, args.productSale.ProductID, nil)
				setupSetCountByProductID(args.ctx, stockR, args.productSale, errRepository)
			},
			wantErr: true,
		},
		{
			name: "error uow",
			args: args{
				ctx:         ctx,
				productSale: &productSaleModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				productSaleR *productsalemock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSaleR, stockR, errRepository)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			productSaleR := productsalemock.NewMockRepository(t)
			stockR := stockmock.NewMockRepository(t)
			uow := transactionsnmock.NewMockUnitOfWork(t)
			s := newTestService(productSaleR, uow)

			tc.mockFn(uow, productSaleR, stockR, tc.args)

			err := s.CreateProductInSale(tc.args.ctx, tc.args.productSale)

			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.NotEqual(t, uuid.Nil, tc.args.productSale.ID)
		})
	}
}

// TestDeleteProductInSale tests method DeleteProductInSale().
//
//nolint:funlen
func TestDeleteProductInSale(t *testing.T) {
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
			productSaleR *productsalemock.MockRepository,
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
				productSaleR *productsalemock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSaleR, stockR, nil)
				setupUpdateStockByProductSale(args.ctx, stockR, args.id, decimal.Zero, nil)
				setupDelete(args.ctx, productSaleR, args.id, nil)
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
				productSaleR *productsalemock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSaleR, stockR, nil)
				setupUpdateStockByProductSale(args.ctx, stockR, args.id, decimal.Zero, errRepository)
			},
			wantErr: true,
		},
		{
			name: "error delete productSale",
			args: args{
				ctx: ctx,
				id:  uuid.New(),
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				productSaleR *productsalemock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSaleR, stockR, nil)
				setupUpdateStockByProductSale(args.ctx, stockR, args.id, decimal.Zero, nil)
				setupDelete(args.ctx, productSaleR, args.id, errRepository)
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
				productSaleR *productsalemock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSaleR, stockR, errRepository)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			productSaleR := productsalemock.NewMockRepository(t)
			stockR := stockmock.NewMockRepository(t)
			uow := transactionsnmock.NewMockUnitOfWork(t)
			s := newTestService(productSaleR, uow)

			tc.mockFn(uow, productSaleR, stockR, tc.args)

			err := s.DeleteProductInSale(tc.args.ctx, tc.args.id)

			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}

// TestDeleteProductsBySaleID tests method DeleteProductsBySaleID().
//
//nolint:funlen
func TestDeleteProductsBySaleID(t *testing.T) {
	type args struct {
		ctx    context.Context
		saleID uuid.UUID
	}

	ctx := context.Background()

	testCases := []struct {
		name   string
		args   args
		mockFn func(
			uow *transactionsnmock.MockUnitOfWork,
			productSaleR *productsalemock.MockRepository,
			stockR *stockmock.MockRepository,
			args args,
		)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:    ctx,
				saleID: uuid.New(),
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				productSaleR *productsalemock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSaleR, stockR, nil)
				setupUpdateStockOnDeleteSale(args.ctx, stockR, args.saleID, nil)
				setupDeleteBySale(args.ctx, productSaleR, args.saleID, nil)
			},
			wantErr: false,
		},
		{
			name: "error set stock",
			args: args{
				ctx:    ctx,
				saleID: uuid.New(),
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				productSaleR *productsalemock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSaleR, stockR, nil)
				setupUpdateStockOnDeleteSale(args.ctx, stockR, args.saleID, errRepository)
			},
			wantErr: true,
		},
		{
			name: "error delete productSale",
			args: args{
				ctx:    ctx,
				saleID: uuid.New(),
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				productSaleR *productsalemock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSaleR, stockR, nil)
				setupUpdateStockOnDeleteSale(args.ctx, stockR, args.saleID, nil)
				setupDeleteBySale(args.ctx, productSaleR, args.saleID, errRepository)
			},
			wantErr: true,
		},
		{
			name: "error uow",
			args: args{
				ctx:    ctx,
				saleID: uuid.New(),
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				productSaleR *productsalemock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSaleR, stockR, errRepository)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			productSaleR := productsalemock.NewMockRepository(t)
			stockR := stockmock.NewMockRepository(t)
			uow := transactionsnmock.NewMockUnitOfWork(t)
			s := newTestService(productSaleR, uow)

			tc.mockFn(uow, productSaleR, stockR, tc.args)

			err := s.DeleteProductsBySaleID(tc.args.ctx, tc.args.saleID)

			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}

// TestUpdateProductInSale tests method UpdateProductInSale().
//
//nolint:funlen
func TestUpdateProductInSale(t *testing.T) {
	type args struct {
		ctx         context.Context
		productSale *models.ProductSale
	}

	productSaleModel := newProductSaleModel()

	ctx := context.Background()

	testCases := []struct {
		name   string
		args   args
		mockFn func(
			uow *transactionsnmock.MockUnitOfWork,
			productSaleR *productsalemock.MockRepository,
			stockR *stockmock.MockRepository,
			args args,
		)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:         ctx,
				productSale: &productSaleModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				productSaleR *productsalemock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSaleR, stockR, nil)
				setupUpdateStockByProductSale(args.ctx, stockR, args.productSale.ID, productSaleModel.Quantity, nil)
				setupUpdate(args.ctx, productSaleR, args.productSale, nil)
			},
			wantErr: false,
		},
		{
			name: "error set stock",
			args: args{
				ctx:         ctx,
				productSale: &productSaleModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				productSaleR *productsalemock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSaleR, stockR, nil)
				setupUpdateStockByProductSale(
					args.ctx,
					stockR,
					args.productSale.ID,
					productSaleModel.Quantity,
					errRepository,
				)
			},
			wantErr: true,
		},
		{
			name: "error update productSale",
			args: args{
				ctx:         ctx,
				productSale: &productSaleModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				productSaleR *productsalemock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSaleR, stockR, nil)
				setupUpdateStockByProductSale(args.ctx, stockR, args.productSale.ID, productSaleModel.Quantity, nil)
				setupUpdate(args.ctx, productSaleR, args.productSale, errRepository)
			},
			wantErr: true,
		},
		{
			name: "error uow",
			args: args{
				ctx:         ctx,
				productSale: &productSaleModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				productSaleR *productsalemock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, productSaleR, stockR, errRepository)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			productSaleR := productsalemock.NewMockRepository(t)
			stockR := stockmock.NewMockRepository(t)
			uow := transactionsnmock.NewMockUnitOfWork(t)
			s := newTestService(productSaleR, uow)

			tc.mockFn(uow, productSaleR, stockR, tc.args)

			err := s.UpdateProductInSale(tc.args.ctx, tc.args.productSale)

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
	productSaleR productsale.Repository,
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
				ProductSale: productSaleR,
				Stock:       stockR,
			}

			return fn(ctx, repos)
		})
}

// setupGetBySale configures the expected behavior method GetBySale() of a productSale repository.
func setupGetBySale(
	ctx context.Context,
	productSaleR *productsalemock.MockRepository,
	saleID uuid.UUID,
	productSaleModels []models.ProductSale,
	err error,
) {
	productSaleR.EXPECT().
		GetBySale(ctx, saleID).
		Return(productSaleModels, err).
		Once()
}

// setupCreate configures the expected behavior method Create() of a productSale repository.
func setupCreate(
	ctx context.Context,
	productSaleR *productsalemock.MockRepository,
	productSale *models.ProductSale,
	err error,
) {
	productSaleR.EXPECT().
		Create(ctx, productSale).
		Return(err).
		Once()
}

// setupUpdate configures the expected behavior method Update() of a productSale repository.
func setupUpdate(
	ctx context.Context,
	productSaleR *productsalemock.MockRepository,
	productSale *models.ProductSale,
	err error,
) {
	productSaleR.EXPECT().
		Update(ctx, productSale).
		Return(err).
		Once()
}

// setupDelete configures the expected behavior method Delete() of a productSale repository.
func setupDelete(
	ctx context.Context,
	productSaleR *productsalemock.MockRepository,
	id uuid.UUID,
	err error,
) {
	productSaleR.EXPECT().
		Delete(ctx, id).
		Return(err).
		Once()
}

// setupDeleteBySale configures the expected behavior method DeleteBySale() of a productSale repository.
func setupDeleteBySale(
	ctx context.Context,
	productSaleR *productsalemock.MockRepository,
	saleID uuid.UUID,
	err error,
) {
	productSaleR.EXPECT().
		DeleteBySale(ctx, saleID).
		Return(err).
		Once()
}

// setupUpdateStockOnDeleteSale configures the expected behavior method UpdateStockOnDeleteSale() of a stock repository.
func setupUpdateStockOnDeleteSale(
	ctx context.Context,
	stockR *stockmock.MockRepository,
	saleID uuid.UUID,
	err error,
) {
	stockR.EXPECT().
		UpdateStockOnDeleteSale(
			ctx,
			saleID,
		).
		Return(err).
		Once()
}

// setupUpdateStockByProductSale sets the expected behavior method UpdateStockByProductSale() of a stock repository.
func setupUpdateStockByProductSale(
	ctx context.Context,
	stockR *stockmock.MockRepository,
	productSaleID uuid.UUID,
	quantity decimal.Decimal,
	err error,
) {
	stockR.EXPECT().
		UpdateStockByProductSale(ctx, productSaleID, quantity).
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

// setupSetCountByProductID configures the expected behavior method SetCountByProductID() of a stock repository.
func setupSetCountByProductID(
	ctx context.Context,
	stockR *stockmock.MockRepository,
	productSale *models.ProductSale,
	err error,
) {
	stockR.EXPECT().
		SetCountByProductID(
			ctx,
			productSale.ProductID,
			productSale.Quantity.Neg(),
		).
		Return(err).
		Once()
}

func requireProductSaleDTO(t *testing.T, want models.ProductSale, got dto.ProductSaleDTO) {
	t.Helper()

	require.Equal(t, want.ID, got.ID)
	require.Equal(t, want.SaleID, got.SaleID)
	require.Equal(t, want.ProductID, got.ProductID)
	require.Equal(t, want.SalePrice, got.SalePrice)
	require.Equal(t, want.Quantity, got.Quantity)
}

// newProductSaleModel creates and returns productSale model with arbitrary values for tests.
func newProductSaleModel() models.ProductSale {
	return models.ProductSale{
		ID:        uuid.New(),
		SaleID:    uuid.New(),
		ProductID: uuid.New(),
		SalePrice: decimal.NewFromInt(500),
		Quantity:  decimal.NewFromInt(20),
	}
}
