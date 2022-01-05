package tesla

import (
	"testing"
)

func TestTeslaApi_SetVerifier(t1 *testing.T) {
	type fields struct {
		username string
		password string
	}
	type args struct {
		customVerifier string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   func(string) bool
	}{
		{
			name:   "Test Verifier length",
			fields: fields{},
			args: args{
				customVerifier: "",
			},
			want: func(got string) bool {
				if len(got) == 86 {
					return true
				}
				return false
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := TeslaApi{}
			if got := t.getVerifier(tt.args.customVerifier); !tt.want(got) {
				t1.Errorf("getVerifier() = %v, check is %v", got, tt.want(got))
			}
		})
	}
}

func TestTeslaApi_getAuthCode(t1 *testing.T) {
	tests := []struct {
		name string
		ta   *TeslaApi
	}{
		{
			name: "Test auth panics",
			ta: NewTeslaApi(
				"",
				"",
				"eyJhbGciOiJSU6IkpXVCIsImtpZCI6Ilg0RmNua0RCUVBUTnBrZTZiMnNuRi04YmdVUSJ9.eyJpc3MiOudGVzbGEuY29tL29hdXRoMi92MyIsImF1ZCI6Imh0dHBzOi8vYXV0aC50ZXNsYS5jb20vb2F1dGgyL3YzL3Rva2VuIiwiaWF0IjoxNjQxMDExNTcyLCJzY3AiOlsib3BlbmlkIiwib2ZmbGluZV9hY2Nlc3MiXSwiZGF0YSI6eyJ2IjoiMSIsImF1ZCI6Imh0dHBzOi8vb3duZXItYXBpLnRlc2xhbW90b3JzLmNvbS8iLCJzdWIiOiI2NmI5NjMxMS0wZDEyLTRjNjgtYjQ1Mi02MTFmOTZhNDkzNGEiLCJzY3AiOlsib3BlbmlkIiwiZW1haWwiLCJvZmZsaW5lX2FjY2VzcyJdLCJhenAiOiJvd25lcmFwaSIsImFtciI6WyJwd2QiXSwiYXV0aF90aW1lIjoxNjQxMDExNTcyfX0.gVIaapAdziniOhSDc2GBp1UwkCLxGFZaqyipw_3tQxIUMqxuKaGSTZR5R8fzx4X0qk2Q3EA9XuQgdcNHDYCByOC0RZinFju1DEe5HdnY6NqCkRcKxbgFFQyF7VtS-OU7Jf8i-KrvImftVyvho59_RofYsgwB0C2ZOaE7UJJyTsG5pX4weGhsIsIu9IIYJbTxemnwcUJdrZsqbe4OMMEeIq6PoU0otSVmRiiHb0wC0ofo0rYPKVrGKYqAuyQxl6osOTx49-Ampu7kvXPNpmxe5abCeJ36btT4-fn6LTd7r1BnlHHyOGSFWz73XzGI9KSsYmsSmVpePe0N9zbESm7wLuz5MKlb2HqJRzH7uUe5L3Ssi2ywaA6R2ZZolf5_0H7exHjguPPg_kErIiuiX79CrTzYBiw3kDu5L-3V0G8TsnOVocSgfDJNXmvZexpzKAQXtKjKfr7iFNcicFyB6MQLy4Hh7hfR2T5GgqHon43YkBYYP3C0ZKxhhMb1KkJ5mwvowYcWPCxrvzc-99nT3orvbgyK1uURmDkx6rUrHY2MI77kIJXn-eC3M1TyOG_SATK8WvIrAOVMl6nEd9-Z4dvsNgg2U2AwilEgETHnuBoSEbG006CkvgqvDKRgWIAlV1n4i3pofrdN5lPt_9S1CtmYIG_Rz_k86zAc",
				true,
			),
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			authCode := tt.ta.getAuthCode()
			tt.ta.getToken(authCode)
			tt.ta.renewToken()
			tt.ta.getVerifier("")
			tt.ta.getChallenge()
			tt.ta.RefreshToken()
			tt.ta.isAuth()
		})
	}
}
