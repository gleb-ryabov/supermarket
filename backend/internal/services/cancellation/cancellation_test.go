package cancellation

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
	cancellationnmock "supermarket/internal/mocks/cancellation"
	stockmock "supermarket/internal/mocks/stock"
	transactionsnmock "supermarket/internal/mocks/transactions"
	"supermarket/internal/models"
	"supermarket/internal/repository/cancellation"
	"supermarket/internal/repository/stock"
	"supermarket/internal/repository/transactions"
)

var errRepository = errors.New("error in repository")

// newTestService creates new cancellation service for tests.
func newTestService(
	cancellationR *cancellationnmock.MockRepository,
	uow *transactionsnmock.MockUnitOfWork) Service {
	return New(
		slogdiscard.NewDiscardLogger(),
		cancellationR,
		uow,
	)
}

// TestGetCancellations tests method GetCancellations().
//
//nolint:funlen
func TestGetCancellations(t *testing.T) {
	type args struct {
		ctx       context.Context
		productID *uuid.UUID
		dateFrom  *time.Time
		dateTo    *time.Time
	}

	cancellationModels := []models.Cancellation{
		newCancellationModel(),
	}

	ctx := context.Background()

	testCases := []struct {
		name    string
		args    args
		mockFn  func(cancellationR *cancellationnmock.MockRepository, args args)
		wantLen int
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:       ctx,
				productID: testhelper.Ptr(uuid.New()),
				dateFrom:  testhelper.Ptr(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
				dateTo:    testhelper.Ptr(time.Date(2027, 12, 31, 0, 0, 0, 0, time.UTC)),
			},
			mockFn: func(cancellationR *cancellationnmock.MockRepository, args args) {
				setupGetByParams(
					args.ctx,
					cancellationR,
					args.productID,
					args.dateFrom,
					args.dateTo,
					cancellationModels,
					nil,
				)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "nil arguments",
			args: args{
				ctx:       ctx,
				productID: nil,
				dateFrom:  nil,
				dateTo:    nil,
			},
			mockFn: func(cancellationR *cancellationnmock.MockRepository, args args) {
				setupGetByParams(
					args.ctx,
					cancellationR,
					args.productID,
					args.dateFrom,
					args.dateTo,
					cancellationModels,
					nil,
				)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "error get in repository",
			args: args{
				ctx:       ctx,
				productID: testhelper.Ptr(uuid.New()),
				dateFrom:  testhelper.Ptr(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
				dateTo:    testhelper.Ptr(time.Date(2027, 12, 31, 0, 0, 0, 0, time.UTC)),
			},
			mockFn: func(cancellationR *cancellationnmock.MockRepository, args args) {
				setupGetByParams(
					args.ctx,
					cancellationR,
					args.productID,
					args.dateFrom,
					args.dateTo,
					cancellationModels,
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

			cancellationR := cancellationnmock.NewMockRepository(t)
			s := newTestService(cancellationR, nil)

			tc.mockFn(cancellationR, tc.args)
			got, err := s.GetCancellations(
				tc.args.ctx,
				tc.args.productID,
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
				requireCancellationDTO(t, cancellationModels[0], got[0])
			}
		})
	}
}

// TestCreateCancellation tests method CreateCancellation().
//
//nolint:funlen
func TestCreateCancellation(t *testing.T) {
	type args struct {
		ctx          context.Context
		cancellation *models.Cancellation
	}

	cancellationModel := newCancellationModel()

	ctx := context.Background()

	testCases := []struct {
		name   string
		args   args
		mockFn func(
			uow *transactionsnmock.MockUnitOfWork,
			cancellationR *cancellationnmock.MockRepository,
			stockR *stockmock.MockRepository,
			args args,
		)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:          ctx,
				cancellation: &cancellationModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				cancellationR *cancellationnmock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, cancellationR, stockR, nil)
				setupCreate(args.ctx, cancellationR, args.cancellation, nil)
				setupFirstOrCreateByProductID(args.ctx, stockR, args.cancellation.ProductID, nil)
				setupSetCountByProductID(args.ctx, stockR, args.cancellation, nil)
			},
			wantErr: false,
		},
		{
			name: "error create cancellation",
			args: args{
				ctx:          ctx,
				cancellation: &cancellationModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				cancellationR *cancellationnmock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, cancellationR, stockR, nil)
				setupCreate(args.ctx, cancellationR, args.cancellation, errRepository)
			},
			wantErr: true,
		},
		{
			name: "error get stock for set stock",
			args: args{
				ctx:          ctx,
				cancellation: &cancellationModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				cancellationR *cancellationnmock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, cancellationR, stockR, nil)
				setupCreate(args.ctx, cancellationR, args.cancellation, nil)
				setupFirstOrCreateByProductID(args.ctx, stockR, args.cancellation.ProductID, errRepository)
			},
			wantErr: true,
		},
		{
			name: "error set stock",
			args: args{
				ctx:          ctx,
				cancellation: &cancellationModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				cancellationR *cancellationnmock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, cancellationR, stockR, nil)
				setupCreate(args.ctx, cancellationR, args.cancellation, nil)
				setupFirstOrCreateByProductID(args.ctx, stockR, args.cancellation.ProductID, nil)
				setupSetCountByProductID(args.ctx, stockR, args.cancellation, errRepository)
			},
			wantErr: true,
		},
		{
			name: "error uow",
			args: args{
				ctx:          ctx,
				cancellation: &cancellationModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				cancellationR *cancellationnmock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, cancellationR, stockR, errRepository)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cancellationR := cancellationnmock.NewMockRepository(t)
			stockR := stockmock.NewMockRepository(t)
			uow := transactionsnmock.NewMockUnitOfWork(t)
			s := newTestService(cancellationR, uow)

			tc.mockFn(uow, cancellationR, stockR, tc.args)

			err := s.CreateCancellation(tc.args.ctx, tc.args.cancellation)

			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.NotEqual(t, uuid.Nil, tc.args.cancellation.ID)
		})
	}
}

// TestDeleteCancellation tests method DeleteCancellation().
//
//nolint:funlen
func TestDeleteCancellation(t *testing.T) {
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
			cancellationR *cancellationnmock.MockRepository,
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
				cancellationR *cancellationnmock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, cancellationR, stockR, nil)
				setupUpdateStockByCancellation(args.ctx, stockR, args.id, decimal.Zero, nil)
				setupDelete(args.ctx, cancellationR, args.id, nil)
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
				cancellationR *cancellationnmock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, cancellationR, stockR, nil)
				setupUpdateStockByCancellation(args.ctx, stockR, args.id, decimal.Zero, errRepository)
			},
			wantErr: true,
		},
		{
			name: "error delete cancellation",
			args: args{
				ctx: ctx,
				id:  uuid.New(),
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				cancellationR *cancellationnmock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, cancellationR, stockR, nil)
				setupUpdateStockByCancellation(args.ctx, stockR, args.id, decimal.Zero, nil)
				setupDelete(args.ctx, cancellationR, args.id, errRepository)
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
				cancellationR *cancellationnmock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, cancellationR, stockR, errRepository)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cancellationR := cancellationnmock.NewMockRepository(t)
			stockR := stockmock.NewMockRepository(t)
			uow := transactionsnmock.NewMockUnitOfWork(t)
			s := newTestService(cancellationR, uow)

			tc.mockFn(uow, cancellationR, stockR, tc.args)

			err := s.DeleteCancellation(tc.args.ctx, tc.args.id)

			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}

// TestUpdateCancellation tests method UpdateCancellation().
//
//nolint:funlen
func TestUpdateCancellation(t *testing.T) {
	type args struct {
		ctx          context.Context
		cancellation *models.Cancellation
	}

	cancellationModel := newCancellationModel()

	ctx := context.Background()

	testCases := []struct {
		name   string
		args   args
		mockFn func(
			uow *transactionsnmock.MockUnitOfWork,
			cancellationR *cancellationnmock.MockRepository,
			stockR *stockmock.MockRepository,
			args args,
		)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:          ctx,
				cancellation: &cancellationModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				cancellationR *cancellationnmock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, cancellationR, stockR, nil)
				setupUpdateStockByCancellation(args.ctx, stockR, args.cancellation.ID, cancellationModel.Quantity, nil)
				setupUpdate(args.ctx, cancellationR, args.cancellation, nil)
			},
			wantErr: false,
		},
		{
			name: "error set stock",
			args: args{
				ctx:          ctx,
				cancellation: &cancellationModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				cancellationR *cancellationnmock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, cancellationR, stockR, nil)
				setupUpdateStockByCancellation(
					args.ctx,
					stockR,
					args.cancellation.ID,
					cancellationModel.Quantity,
					errRepository,
				)
			},
			wantErr: true,
		},
		{
			name: "error update cancellation",
			args: args{
				ctx:          ctx,
				cancellation: &cancellationModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				cancellationR *cancellationnmock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, cancellationR, stockR, nil)
				setupUpdateStockByCancellation(args.ctx, stockR, args.cancellation.ID, cancellationModel.Quantity, nil)
				setupUpdate(args.ctx, cancellationR, args.cancellation, errRepository)
			},
			wantErr: true,
		},
		{
			name: "error uow",
			args: args{
				ctx:          ctx,
				cancellation: &cancellationModel,
			},
			mockFn: func(
				uow *transactionsnmock.MockUnitOfWork,
				cancellationR *cancellationnmock.MockRepository,
				stockR *stockmock.MockRepository,
				args args,
			) {
				setupDoUOW(args.ctx, uow, cancellationR, stockR, errRepository)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cancellationR := cancellationnmock.NewMockRepository(t)
			stockR := stockmock.NewMockRepository(t)
			uow := transactionsnmock.NewMockUnitOfWork(t)
			s := newTestService(cancellationR, uow)

			tc.mockFn(uow, cancellationR, stockR, tc.args)

			err := s.UpdateCancellation(tc.args.ctx, tc.args.cancellation)

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
	cancellationR cancellation.Repository,
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
				Cancellation: cancellationR,
				Stock:        stockR,
			}

			return fn(ctx, repos)
		})
}

// setupGetByParams configures the expected behavior method GetByParams() of a cancellation repository.
func setupGetByParams(
	ctx context.Context,
	cancellationR *cancellationnmock.MockRepository,
	productID *uuid.UUID,
	dateFrom *time.Time,
	dateTo *time.Time,
	cancellationModels []models.Cancellation,
	err error,
) {
	cancellationR.EXPECT().
		GetByParams(ctx, productID, dateFrom, dateTo).
		Return(cancellationModels, err).
		Once()
}

// setupCreate configures the expected behavior method Create() of a cancellation repository.
func setupCreate(
	ctx context.Context,
	cancellationR *cancellationnmock.MockRepository,
	cancellation *models.Cancellation,
	err error,
) {
	cancellationR.EXPECT().
		Create(ctx, cancellation).
		Return(err).
		Once()
}

// setupUpdate configures the expected behavior method Update() of a cancellation repository.
func setupUpdate(
	ctx context.Context,
	cancellationR *cancellationnmock.MockRepository,
	cancellation *models.Cancellation,
	err error,
) {
	cancellationR.EXPECT().
		Update(ctx, cancellation).
		Return(err).
		Once()
}

// setupDelete configures the expected behavior method Delete() of a cancellation repository.
func setupDelete(
	ctx context.Context,
	cancellationR *cancellationnmock.MockRepository,

	id uuid.UUID,
	err error,
) {
	cancellationR.EXPECT().
		Delete(ctx, id).
		Return(err).
		Once()
}

// setupSetCountByProductID configures the expected behavior method SetCountByProductID() of a stock repository.
func setupSetCountByProductID(
	ctx context.Context,
	stockR *stockmock.MockRepository,
	cancellation *models.Cancellation,
	err error,
) {
	stockR.EXPECT().
		SetCountByProductID(
			ctx,
			cancellation.ProductID,
			cancellation.Quantity.Neg(),
		).
		Return(err).
		Once()
}

// setupUpdateStockByCancellation sets the expected behavior method UpdateStockByCancellation() of a stock repository.
func setupUpdateStockByCancellation(
	ctx context.Context,
	stockR *stockmock.MockRepository,
	cancellationID uuid.UUID,
	quantity decimal.Decimal,
	err error,
) {
	stockR.EXPECT().
		UpdateStockByCancellation(ctx, cancellationID, quantity).
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

func requireCancellationDTO(t *testing.T, want models.Cancellation, got dto.CancellationDTO) {
	t.Helper()

	require.Equal(t, want.ID, got.ID)
	require.Equal(t, want.DateTime.Format("02.01.2006 15:04"), got.DateTime)
	require.Equal(t, want.ProductID, got.ProductID)
	require.Equal(t, want.Quantity, got.Quantity)
}

// newCancellationModel creates and returns cancellation model with arbitrary values for tests.
func newCancellationModel() models.Cancellation {
	return models.Cancellation{
		ID:        uuid.New(),
		DateTime:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		ProductID: uuid.New(),
		Quantity:  decimal.NewFromInt(20),
	}
}
