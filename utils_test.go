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
