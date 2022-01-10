package tesla

import "testing"

func Test_isTokenExpired(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name        string
		args        args
		wantExpired bool
	}{
		{
			name:        "Test access token expired",
			args:        args{token: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE1MTYyMzkwMjJ9.4Adcj3UFYzPUVaVF43FmMab6RlaQD8A9V8wFzzht-KQ`},
			wantExpired: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotExpired := isTokenExpired(tt.args.token); gotExpired != tt.wantExpired {
				t.Errorf("isTokenExpired() = %v, want %v", gotExpired, tt.wantExpired)
			}
		})
	}
}

func Test_joinPath(t *testing.T) {
	type args struct {
		baseUrl string
		paths   []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test join path",
			args: args{
				baseUrl: "https://www.tesla.com",
				paths:   []string{"/api/v1", "vehicleData"},
			},
			want: "https://www.tesla.com/api/v1/vehicleData",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := joinPath(tt.args.baseUrl, tt.args.paths...); got != tt.want {
				t.Errorf("joinPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
