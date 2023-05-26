package countrycodes

import (
	"reflect"
	"testing"
)

func TestFindByName(t *testing.T) {
	cc, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	matches := cc.FindByName("United States Minor")

	if len(matches) != 1 {
		t.Fatalf("Extra matches found")
	}

	um, _ := cc.GetByAlpha2("UM")

	if matches[0] != um {
		t.Fatalf("Match for United States Minor Outlying Islands failed")
	}
}

func TestGetByNumeric(t *testing.T) {
	cc, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	code, _ := cc.GetByNumeric(840)

	if code.Name != "United States of America" {
		t.Fatalf("GetByNumeric failed, got: %s", code.Name)
	}
}

func TestNewClient(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		t.Fatal(err)
	}
}

func Test_validate(t *testing.T) {
	type args struct {
		alpha2 map[string]CountryCode
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "simple passing case",
			args: args{
				alpha2: map[string]CountryCode{
					"MU": {
						Name:       "Made Up",
						Alpha2:     "MU",
						Alpha3:     "MDU",
						Numeric:    999,
						Assignment: 3,
					}},
			},
			wantErr: false,
		},
		{
			name: "alpha2 and code mismatch",
			args: args{
				alpha2: map[string]CountryCode{
					"MU": {
						Name:       "Made Up",
						Alpha2:     "MB",
						Alpha3:     "MDU",
						Numeric:    999,
						Assignment: 3,
					}},
			},
			wantErr: true,
		},
		{
			name: "alpha2 len",
			args: args{
				alpha2: map[string]CountryCode{
					"MU": {
						Name:       "Made Up",
						Alpha2:     "MB1",
						Alpha3:     "MDU",
						Numeric:    999,
						Assignment: 3,
					}},
			},
			wantErr: true,
		},
		{
			name: "alpha3 len",
			args: args{
				alpha2: map[string]CountryCode{
					"MU": {
						Name:       "Made Up",
						Alpha2:     "MBU",
						Alpha3:     "MDU",
						Numeric:    999,
						Assignment: 3,
					}},
			},
			wantErr: true,
		},
		{
			name: "invalid assignment",
			args: args{
				alpha2: map[string]CountryCode{
					"MU": {
						Name:       "Made Up",
						Alpha2:     "MU",
						Alpha3:     "MDU",
						Numeric:    999,
						Assignment: 99,
					}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validate(tt.args.alpha2); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_GetByAlpha2(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		a2 string
	}
	tests := []struct {
		name   string
		fields *Client
		args   args
		want   CountryCode
		want1  bool
	}{
		{
			name:   "simple passing case",
			fields: client,
			args:   args{},
			want:   CountryCode{},
			want1:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := &Client{
				byAlpha2:  tt.fields.byAlpha2,
				byName:    tt.fields.byName,
				byAlpha3:  tt.fields.byAlpha3,
				byNumeric: tt.fields.byNumeric,
				nameTrie:  tt.fields.nameTrie,
			}
			got, got1 := cc.GetByAlpha2(tt.args.a2)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByAlpha2() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetByAlpha2() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
