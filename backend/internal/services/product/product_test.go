package product

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
	productmock "supermarket/internal/mocks/product"
	"supermarket/internal/models"
)

var errRepository = errors.New("error in repository")

// newTestService creates new product service for tests.
func newTestService(productR *productmock.MockRepository) Service {
	return New(
		slogdiscard.NewDiscardLogger(),
		productR,
	)
}

// TestGetProducts tests method GetProducts().
//
//nolint:funlen
func TestGetProducts(t *testing.T) {
	type args struct {
		ctx    context.Context
		name   string
		typeID *uuid.UUID
	}

	productModels := []models.Product{
		newProductModel(),
	}

	testCases := []struct {
		name    string
		args    args
		mockFn  func(productR *productmock.MockRepository, args args)
		wantLen int
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				name:   "milk",
				typeID: testhelper.Ptr(uuid.New()),
			},
			mockFn: func(productR *productmock.MockRepository, args args) {
				productR.EXPECT().
					GetByParams(args.ctx, args.name, args.typeID).
					Return(productModels, nil).
					Once()
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "nil product type ID",
			args: args{
				ctx:    context.Background(),
				name:   "milk",
				typeID: nil,
			},
			mockFn: func(productR *productmock.MockRepository, args args) {
				productR.EXPECT().
					GetByParams(args.ctx, args.name, args.typeID).
					Return(productModels, nil).
					Once()
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "error get from repository",
			args: args{
				ctx:    context.Background(),
				name:   "milk",
				typeID: testhelper.Ptr(uuid.New()),
			},
			mockFn: func(productR *productmock.MockRepository, args args) {
				productR.EXPECT().
					GetByParams(args.ctx, args.name, args.typeID).
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

			productR := productmock.NewMockRepository(t)
			s := newTestService(productR)

			tc.mockFn(productR, tc.args)

			got, err := s.GetProducts(tc.args.ctx, tc.args.name, tc.args.typeID)

			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, got)

				return
			}

			require.NoError(t, err)
			require.Len(t, got, tc.wantLen)

			if len(got) > 0 {
				requireProductDTO(t, productModels[0], got[0])
			}
		})
	}
}

// TestCreateProduct tests method CreateProduct().
//
//nolint:funlen
func TestCreateProduct(t *testing.T) {
	type args struct {
		ctx     context.Context
		product *models.Product
	}

	productModel := newProductModel()

	testCases := []struct {
		name    string
		args    args
		mockFn  func(productR *productmock.MockRepository, args args)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:     context.Background(),
				product: &productModel,
			},
			mockFn: func(productR *productmock.MockRepository, args args) {
				productR.EXPECT().
					Create(args.ctx, args.product).
					Return(nil).
					Once()
			},
			wantErr: false,
		},
		{
			name: "error create in repository",
			args: args{
				ctx:     context.Background(),
				product: &productModel,
			},
			mockFn: func(productR *productmock.MockRepository, args args) {
				productR.EXPECT().
					Create(args.ctx, args.product).
					Return(errRepository).
					Once()
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			productR := productmock.NewMockRepository(t)
			s := newTestService(productR)

			tc.mockFn(productR, tc.args)

			err := s.CreateProduct(tc.args.ctx, tc.args.product)

			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.NotEqual(t, uuid.Nil, tc.args.product.ID)
		})
	}
}

// TestDeleteProduct tests method DeleteProduct().
//

func TestDeleteProduct(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}

	testCases := []struct {
		name    string
		args    args
		mockFn  func(productR *productmock.MockRepository, args args)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				id:  uuid.New(),
			},
			mockFn: func(productR *productmock.MockRepository, args args) {
				productR.EXPECT().
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
			mockFn: func(productR *productmock.MockRepository, args args) {
				productR.EXPECT().
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

			productR := productmock.NewMockRepository(t)
			s := newTestService(productR)

			tc.mockFn(productR, tc.args)
			err := s.DeleteProduct(tc.args.ctx, tc.args.id)

			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}

// TestUpdateProduct tests method UpdateProduct().
//
//nolint:funlen
func TestUpdateProduct(t *testing.T) {
	type args struct {
		ctx context.Context
		pt  *models.Product
	}

	productModel := newProductModel()

	testCases := []struct {
		name    string
		args    args
		mockFn  func(productR *productmock.MockRepository, args args)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				pt:  &productModel,
			},
			mockFn: func(productR *productmock.MockRepository, args args) {
				productR.EXPECT().
					Update(args.ctx, args.pt).
					Return(nil).
					Once()
			},
			wantErr: false,
		},
		{
			name: "error update in repository",
			args: args{
				ctx: context.Background(),
				pt:  &productModel,
			},
			mockFn: func(productR *productmock.MockRepository, args args) {
				productR.EXPECT().
					Update(args.ctx, args.pt).
					Return(errRepository).
					Once()
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			productR := productmock.NewMockRepository(t)
			s := newTestService(productR)

			tc.mockFn(productR, tc.args)
			err := s.UpdateProduct(tc.args.ctx, tc.args.pt)

			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}

func requireProductDTO(t *testing.T, want models.Product, got dto.ProductDTO) {
	t.Helper()

	require.Equal(t, want.ID, got.ID)
	require.Equal(t, want.Name, got.Name)
	require.Equal(t, want.Manufacturer, got.Manufacturer)
	require.Equal(t, want.Weight, got.Weight)
	require.Equal(t, want.Volume, got.Volume)
}

// newProductModel creates and returns product model with arbitrary values for tests.
func newProductModel() models.Product {
	return models.Product{
		ID:           uuid.New(),
		TypeID:       uuid.New(),
		Name:         "Milk",
		Manufacturer: "MilkManufacturer",
		Weight:       decimal.NewFromInt(1),
		Volume:       decimal.NewFromInt(1),
	}
}
