package greetings

import "testing"

func TestHello(t *testing.T) {
	type args struct {
		buildVersion string
		buildDate    string
		buildCommit  string
	}
	tests := []struct {
		name string
		args args
		want *Greetings
	}{
		{
			name: "check emty",
			args: args{
				buildVersion: "",
				buildDate:    "",
				buildCommit:  "",
			},
			want: &Greetings{
				BuildVersion: "N/A",
				BuildDate:    "N/A",
				BuildCommit:  "N/A",
			},
		},
		{
			name: "check positive test",
			args: args{
				buildVersion: "1",
				buildDate:    "2",
				buildCommit:  "3",
			},
			want: &Greetings{
				BuildVersion: "1",
				BuildDate:    "2",
				BuildCommit:  "3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Hello(tt.args.buildVersion, tt.args.buildDate, tt.args.buildCommit); err != nil {
				t.Errorf("Hello() error = %v, wantErr %v", err, tt.want)
			}
		})
	}
}
