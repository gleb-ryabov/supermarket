package producttype

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"supermarket/internal/http/dto"
	slogdiscard "supermarket/internal/lib/logger"
	testhelper "supermarket/internal/lib/test"
	productTypemock "supermarket/internal/mocks/product_type"
	"supermarket/internal/models"
)

var errRepository = errors.New("error in repository")

// newTestService creates new productType service for tests.
func newTestService(productTypeR *productTypemock.MockRepository) Service {
	return New(
		slogdiscard.NewDiscardLogger(),
		productTypeR,
	)
}

// TestGetProductTypes tests method GetProductTypes().
//
//nolint:funlen
func TestGetProductTypes(t *testing.T) {
	type args struct {
		ctx      context.Context
		name     string
		forAdult *bool
	}

	productTypeModels := []models.ProductType{
		newProductTypeModel(),
	}

	testCases := []struct {
		name    string
		args    args
		mockFn  func(repo *productTypemock.MockRepository, args args)
		wantLen int
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:      context.Background(),
				name:     "молочные",
				forAdult: testhelper.Ptr(false),
			},
			mockFn: func(productTypeR *productTypemock.MockRepository, args args) {
				productTypeR.EXPECT().
					GetByParams(
						args.ctx,
						args.name,
						args.forAdult,
					).
					Return(productTypeModels, nil).
					Once()
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "nil arguments",
			args: args{
				ctx:      context.Background(),
				name:     "",
				forAdult: nil,
			},
			mockFn: func(productTypeR *productTypemock.MockRepository, args args) {
				productTypeR.EXPECT().
					GetByParams(
						args.ctx,
						args.name,
						args.forAdult,
					).
					Return(productTypeModels, nil).
					Once()
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "error get from repository",
			args: args{
				ctx:      context.Background(),
				name:     "молочные",
				forAdult: testhelper.Ptr(false),
			},
			mockFn: func(productTypeR *productTypemock.MockRepository, args args) {
				productTypeR.EXPECT().
					GetByParams(
						args.ctx,
						args.name,
						args.forAdult,
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

			productTypeR := productTypemock.NewMockRepository(t)
			s := newTestService(productTypeR)

			tc.mockFn(productTypeR, tc.args)

			got, err := s.GetProductTypes(
				tc.args.ctx,
				tc.args.name,
				tc.args.forAdult,
			)

			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, got)

				return
			}

			require.NoError(t, err)
			require.Len(t, got, tc.wantLen)

			if len(got) > 0 {
				requireProductTypeDTO(t, productTypeModels[0], got[0])
			}
		})
	}
}

// TestCreateProductType tests method CreateProductType().
//
//nolint:funlen
func TestCreateProductType(t *testing.T) {
	type args struct {
		ctx         context.Context
		productType *models.ProductType
	}

	productTypeModel := newProductTypeModel()

	testCases := []struct {
		name    string
		args    args
		mockFn  func(productTypeR *productTypemock.MockRepository, args args)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:         context.Background(),
				productType: &productTypeModel,
			},
			mockFn: func(productTypeR *productTypemock.MockRepository, args args) {
				productTypeR.EXPECT().
					Create(args.ctx, &productTypeModel).
					Return(nil).
					Once()
			},
			wantErr: false,
		},
		{
			name: "error create in repository",
			args: args{
				ctx:         context.Background(),
				productType: &productTypeModel,
			},
			mockFn: func(productTypeR *productTypemock.MockRepository, args args) {
				productTypeR.EXPECT().
					Create(args.ctx, &productTypeModel).
					Return(errRepository).
					Once()
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			productTypeR := productTypemock.NewMockRepository(t)
			s := newTestService(productTypeR)

			tc.mockFn(productTypeR, tc.args)

			err := s.CreateProductType(tc.args.ctx, tc.args.productType)

			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.NotEqual(t, uuid.Nil, tc.args.productType.ID)
		})
	}
}

// TestDeleteProductType tests method DeleteProductType().
//
//nolint:funlen
func TestDeleteProductType(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}

	testCases := []struct {
		name    string
		args    args
		mockFn  func(productTypeR *productTypemock.MockRepository, args args)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				id:  uuid.New(),
			},
			mockFn: func(productTypeR *productTypemock.MockRepository, args args) {
				productTypeR.EXPECT().
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
			mockFn: func(productTypeR *productTypemock.MockRepository, args args) {
				productTypeR.EXPECT().
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

			productTypeR := productTypemock.NewMockRepository(t)
			s := newTestService(productTypeR)

			tc.mockFn(productTypeR, tc.args)

			err := s.DeleteProductType(tc.args.ctx, tc.args.id)

			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}

// TestUpdateProductType tests method UpdateProductType().
//
//nolint:funlen
func TestUpdateProductType(t *testing.T) {
	type args struct {
		ctx         context.Context
		productType *models.ProductType
	}

	productTypeModel := newProductTypeModel()

	testCases := []struct {
		name    string
		args    args
		mockFn  func(productTypeR *productTypemock.MockRepository, args args)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:         context.Background(),
				productType: &productTypeModel,
			},
			mockFn: func(productTypeR *productTypemock.MockRepository, args args) {
				productTypeR.EXPECT().
					Update(args.ctx, args.productType).
					Return(nil).
					Once()
			},
			wantErr: false,
		},
		{
			name: "error update in repository",
			args: args{
				ctx:         context.Background(),
				productType: &productTypeModel,
			},
			mockFn: func(productTypeR *productTypemock.MockRepository, args args) {
				productTypeR.EXPECT().
					Update(args.ctx, args.productType).
					Return(errRepository).
					Once()
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			productTypeR := productTypemock.NewMockRepository(t)
			s := newTestService(productTypeR)

			tc.mockFn(productTypeR, tc.args)

			err := s.UpdateProductType(tc.args.ctx, tc.args.productType)

			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}

func requireProductTypeDTO(t *testing.T, want models.ProductType, got dto.ProductTypeDTO) {
	t.Helper()

	require.Equal(t, want.ID, got.ID)
	require.Equal(t, want.Name, got.Name)
	require.Equal(t, want.ForAdult, got.ForAdult)
}

// newProductTypeModel creates and returns productType model with arbitrary values for tests.
func newProductTypeModel() models.ProductType {
	return models.ProductType{
		ID:       uuid.New(),
		Name:     "Молочная продукция",
		ForAdult: false,
	}
}
