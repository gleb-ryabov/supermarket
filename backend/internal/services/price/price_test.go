package price

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"supermarket/internal/http/dto"
	slogdiscard "supermarket/internal/lib/logger"
	testhelper "supermarket/internal/lib/test"
	pricemock "supermarket/internal/mocks/price"
	"supermarket/internal/models"
)

var errRepository error = errors.New("error in repository")

// newTestService creates new price service for tests.
func newTestService(priceR *pricemock.MockRepository) Service {
	return New(
		slogdiscard.NewDiscardLogger(),
		priceR,
	)
}

// TestGetPrices tests method GetPrices().
//
//nolint:funlen
func TestGetPrices(t *testing.T) {
	type args struct {
		ctx      context.Context
		typeID   *uuid.UUID
		dateFrom *time.Time
		dateTo   *time.Time
	}

	priceModels := []models.Price{
		newPriceModel(),
	}

	testCases := []struct {
		name    string
		args    args
		mockFn  func(repo *pricemock.MockRepository, args args)
		wantLen int
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:      context.Background(),
				typeID:   testhelper.Ptr(uuid.New()),
				dateFrom: testhelper.Ptr(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
				dateTo:   testhelper.Ptr(time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)),
			},
			mockFn: func(priceR *pricemock.MockRepository, args args) {
				priceR.EXPECT().
					GetByParams(
						args.ctx,
						args.typeID,
						args.dateFrom,
						args.dateTo,
					).
					Return(priceModels, nil).
					Once()
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "error get from repository",
			args: args{
				ctx:      context.Background(),
				typeID:   testhelper.Ptr(uuid.New()),
				dateFrom: testhelper.Ptr(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
				dateTo:   testhelper.Ptr(time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)),
			},
			mockFn: func(priceR *pricemock.MockRepository, args args) {
				priceR.EXPECT().
					GetByParams(
						args.ctx,
						args.typeID,
						args.dateFrom,
						args.dateTo,
					).
					Return(nil, errRepository).
					Once()
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			priceR := pricemock.NewMockRepository(t)
			s := newTestService(priceR)

			tc.mockFn(priceR, tc.args)

			got, err := s.GetPrices(
				tc.args.ctx,
				tc.args.typeID,
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
				requirePriceDTO(t, priceModels[0], got[0])
			}
		})
	}
}

// TestCreatePrice tests method CreatePrice().
//
//nolint:funlen
func TestCreatePrice(t *testing.T) {
	type args struct {
		ctx   context.Context
		price *models.Price
	}

	priceModel := newPriceModel()

	testCases := []struct {
		name    string
		args    args
		mockFn  func(priceR *pricemock.MockRepository, args args)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:   context.Background(),
				price: &priceModel,
			},
			mockFn: func(priceR *pricemock.MockRepository, args args) {
				priceR.EXPECT().
					Create(args.ctx, &priceModel).
					Return(nil).
					Once()
			},
			wantErr: false,
		},
		{
			name: "error create in repository",
			args: args{
				ctx:   context.Background(),
				price: &priceModel,
			},
			mockFn: func(priceR *pricemock.MockRepository, args args) {
				priceR.EXPECT().
					Create(args.ctx, &priceModel).
					Return(errRepository).
					Once()
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			priceR := pricemock.NewMockRepository(t)
			s := newTestService(priceR)

			tc.mockFn(priceR, tc.args)

			err := s.CreatePrice(tc.args.ctx, tc.args.price)

			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotEqual(t, uuid.Nil, tc.args.price.ID)
		})
	}
}

// TestDeletePrice tests method DeletePrice().
//
//nolint:funlen
func TestDeletePrice(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}

	testCases := []struct {
		name    string
		args    args
		mockFn  func(priceR *pricemock.MockRepository, args args)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				id:  uuid.New(),
			},
			mockFn: func(priceR *pricemock.MockRepository, args args) {
				priceR.EXPECT().
					Delete(args.ctx, args.id).
					Return(nil).
					Once()
			},
			wantErr: false,
		},
		{
			name: "error delete in repository",
			args: args{
				ctx: context.Background(),
				id:  uuid.New(),
			},
			mockFn: func(priceR *pricemock.MockRepository, args args) {
				priceR.EXPECT().
					Delete(args.ctx, args.id).
					Return(errRepository).
					Once()
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			priceR := pricemock.NewMockRepository(t)
			s := newTestService(priceR)

			tc.mockFn(priceR, tc.args)

			err := s.DeletePrice(tc.args.ctx, tc.args.id)

			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

// TestUpdatePrice tests method UpdatePrice().
//
//nolint:funlen
func TestUpdatePrice(t *testing.T) {
	type args struct {
		ctx   context.Context
		price *models.Price
	}

	priceModel := newPriceModel()

	testCases := []struct {
		name    string
		args    args
		mockFn  func(priceR *pricemock.MockRepository, args args)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:   context.Background(),
				price: &priceModel,
			},
			mockFn: func(priceR *pricemock.MockRepository, args args) {
				priceR.EXPECT().
					Update(args.ctx, args.price).
					Return(nil).
					Once()
			},
			wantErr: false,
		},
		{
			name: "error update in repository",
			args: args{
				ctx:   context.Background(),
				price: &priceModel,
			},
			mockFn: func(priceR *pricemock.MockRepository, args args) {
				priceR.EXPECT().
					Update(args.ctx, args.price).
					Return(errRepository).
					Once()
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			priceR := pricemock.NewMockRepository(t)
			s := newTestService(priceR)

			tc.mockFn(priceR, tc.args)

			err := s.UpdatePrice(tc.args.ctx, tc.args.price)

			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func requirePriceDTO(t *testing.T, want models.Price, got dto.PriceDTO) {
	t.Helper()

	require.Equal(t, want.ID, got.ID)
	require.Equal(t, want.DateStart.Format("02.01.2006"), got.DateStart)
	require.True(t, want.Discount.Equal(got.Discount))
	require.True(t, want.FullPrice.Equal(got.FullPrice))
	require.True(t, want.TotalPrice.Equal(got.TotalPrice))

	var dateEnd string
	if want.DateEnd != nil {
		dateEnd = want.DateEnd.Format("02.01.2006")
	}

	require.Equal(t, dateEnd, got.DateEnd)
}

// newPriceModel creates and returns price model with arbitrary values for tests.
func newPriceModel() models.Price {
	return models.Price{
		ID:         uuid.New(),
		ProductID:  uuid.New(),
		DateStart:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		DateEnd:    testhelper.Ptr(time.Date(2027, 12, 31, 0, 0, 0, 0, time.UTC)),
		Discount:   decimal.NewFromInt(10),
		FullPrice:  decimal.NewFromInt(100),
		TotalPrice: decimal.NewFromInt(90),
	}
}
