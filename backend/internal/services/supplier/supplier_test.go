package supplier

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"supermarket/internal/http/dto"
	slogdiscard "supermarket/internal/lib/logger"
	suppliermock "supermarket/internal/mocks/supplier"
	"supermarket/internal/models"
)

var errRepository = errors.New("error in repository")

// newTestService creates new supplier service for tests.
func newTestService(supplierR *suppliermock.MockRepository) Service {
	return New(
		slogdiscard.NewDiscardLogger(),
		supplierR,
	)
}

// TestGetSuppliers tests method GetSuppliers().
//
//nolint:funlen
func TestGetSuppliers(t *testing.T) {
	type args struct {
		ctx    context.Context
		search string
	}

	supplierModels := []models.Supplier{
		newSupplierModel(),
	}

	testCases := []struct {
		name    string
		args    args
		mockFn  func(repo *suppliermock.MockRepository, args args)
		wantLen int
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				search: "наименование",
			},
			mockFn: func(supplierR *suppliermock.MockRepository, args args) {
				supplierR.EXPECT().
					GetByParams(
						args.ctx,
						args.search,
					).
					Return(supplierModels, nil).
					Once()
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "error get from repository",
			args: args{
				ctx:    context.Background(),
				search: "наименование",
			},
			mockFn: func(supplierR *suppliermock.MockRepository, args args) {
				supplierR.EXPECT().
					GetByParams(
						args.ctx,
						args.search,
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

			supplierR := suppliermock.NewMockRepository(t)
			s := newTestService(supplierR)

			tc.mockFn(supplierR, tc.args)

			got, err := s.GetSuppliers(
				tc.args.ctx,
				tc.args.search,
			)

			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, got)

				return
			}

			require.NoError(t, err)
			require.Len(t, got, tc.wantLen)

			if len(got) > 0 {
				requireSupplierDTO(t, supplierModels[0], got[0])
			}
		})
	}
}

// TestCreateSupplier tests method CreateSupplier().
//
//nolint:funlen
func TestCreateSupplier(t *testing.T) {
	type args struct {
		ctx      context.Context
		supplier *models.Supplier
	}

	supplierModel := newSupplierModel()

	testCases := []struct {
		name    string
		args    args
		mockFn  func(supplierR *suppliermock.MockRepository, args args)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:      context.Background(),
				supplier: &supplierModel,
			},
			mockFn: func(supplierR *suppliermock.MockRepository, args args) {
				supplierR.EXPECT().
					Create(args.ctx, &supplierModel).
					Return(nil).
					Once()
			},
			wantErr: false,
		},
		{
			name: "error create in repository",
			args: args{
				ctx:      context.Background(),
				supplier: &supplierModel,
			},
			mockFn: func(supplierR *suppliermock.MockRepository, args args) {
				supplierR.EXPECT().
					Create(args.ctx, &supplierModel).
					Return(errRepository).
					Once()
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			supplierR := suppliermock.NewMockRepository(t)
			s := newTestService(supplierR)

			tc.mockFn(supplierR, tc.args)

			err := s.CreateSupplier(tc.args.ctx, tc.args.supplier)

			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.NotEqual(t, uuid.Nil, tc.args.supplier.ID)
		})
	}
}

// TestDeleteSupplier tests method DeleteSupplier().
//
//nolint:funlen
func TestDeleteSupplier(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}

	testCases := []struct {
		name    string
		args    args
		mockFn  func(supplierR *suppliermock.MockRepository, args args)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				id:  uuid.New(),
			},
			mockFn: func(supplierR *suppliermock.MockRepository, args args) {
				supplierR.EXPECT().
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
			mockFn: func(supplierR *suppliermock.MockRepository, args args) {
				supplierR.EXPECT().
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

			supplierR := suppliermock.NewMockRepository(t)
			s := newTestService(supplierR)

			tc.mockFn(supplierR, tc.args)

			err := s.DeleteSupplier(tc.args.ctx, tc.args.id)

			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}

// TestUpdateSupplier tests method UpdateSupplier().
//
//nolint:funlen
func TestUpdateSupplier(t *testing.T) {
	type args struct {
		ctx      context.Context
		supplier *models.Supplier
	}

	supplierModel := newSupplierModel()

	testCases := []struct {
		name    string
		args    args
		mockFn  func(supplierR *suppliermock.MockRepository, args args)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:      context.Background(),
				supplier: &supplierModel,
			},
			mockFn: func(supplierR *suppliermock.MockRepository, args args) {
				supplierR.EXPECT().
					Update(args.ctx, args.supplier).
					Return(nil).
					Once()
			},
			wantErr: false,
		},
		{
			name: "error update in repository",
			args: args{
				ctx:      context.Background(),
				supplier: &supplierModel,
			},
			mockFn: func(supplierR *suppliermock.MockRepository, args args) {
				supplierR.EXPECT().
					Update(args.ctx, args.supplier).
					Return(errRepository).
					Once()
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			supplierR := suppliermock.NewMockRepository(t)
			s := newTestService(supplierR)

			tc.mockFn(supplierR, tc.args)

			err := s.UpdateSupplier(tc.args.ctx, tc.args.supplier)

			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}

func requireSupplierDTO(t *testing.T, want models.Supplier, got dto.SupplierDTO) {
	t.Helper()

	require.Equal(t, want.ID, got.ID)
	require.Equal(t, want.Name, got.Name)
	require.Equal(t, want.INN, got.INN)
	require.Equal(t, want.KPP, got.KPP)
	require.Equal(t, want.OGRN, got.OGRN)
	require.Equal(t, want.Phone, got.Phone)
	require.Equal(t, want.Email, got.Email)
}

// newSupplierModel creates and returns supplier model with arbitrary values for tests.
func newSupplierModel() models.Supplier {
	return models.Supplier{
		ID:    uuid.New(),
		Name:  "Наименование",
		INN:   "123456789012",
		KPP:   "123456789",
		OGRN:  "1234567890123",
		Phone: "89991110101",
		Email: "test@email.ru",
	}
}
