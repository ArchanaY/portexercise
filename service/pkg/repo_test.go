package repo

import (
	"context"
	"portexercise/service/domain"
	"reflect"
	"testing"
)

func TestInMemoryStore_Insert(t *testing.T) {
	type args struct {
		ctx context.Context
		pi  domain.Port
	}
	tests := []struct {
		name    string
		is      *InMemoryStore
		args    args
		wantErr bool
	}{
		{
			name: "Happy Days",
			is:   NewInMemoryStore(),
			args: args{
				ctx: context.TODO(),
				pi: domain.Port{
					Name:     "name1",
					City:     "city1",
					Country:  "country1",
					Province: "province1",
					Timezone: "timezone1",
					Code:     "code1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.is.Insert(tt.args.ctx, tt.args.pi); (err != nil) != tt.wantErr {
				t.Errorf("InMemoryStore.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInMemoryStore_Fetch(t *testing.T) {
	pi := domain.Port{
		Key:         "key1",
		Name:        "name1",
		City:        "city1",
		Country:     "country1",
		Alias:       []string{"alias1"},
		Regions:     []string{"region1"},
		Coordinates: []float64{1.1111111},
		Province:    "province1",
		Timezone:    "timezone1",
		Unlocs:      []string{"unlocs1"},
		Code:        "code1",
	}
	st := NewInMemoryStore()
	_ = st.Insert(context.TODO(), pi)
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		is      *InMemoryStore
		args    args
		want    domain.Port
		wantErr bool
	}{
		{
			name: "Happy Days",
			is:   st,
			args: args{
				ctx: context.TODO(),
				key: "key1",
			},
			want: domain.Port{
				Name:        "name1",
				City:        "city1",
				Country:     "country1",
				Alias:       []string{"alias1"},
				Regions:     []string{"region1"},
				Coordinates: []float64{1.1111111},
				Province:    "province1",
				Timezone:    "timezone1",
				Unlocs:      []string{"unlocs1"},
				Code:        "code1",
			},
			wantErr: false,
		},
		{
			name: "Unhappy Days",
			is:   st,
			args: args{
				ctx: context.TODO(),
				key: "unknown",
			},
			want:    domain.Port{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.is.Fetch(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("InMemoryStore.Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InMemoryStore.Fetch() = %v, want %v", got, tt.want)
			}
		})
	}
}
