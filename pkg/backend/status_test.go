package backend_test

import (
	"github.com/puppetlabs/prm/pkg/backend"
	"github.com/puppetlabs/prm/pkg/config"
	"reflect"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/puppetlabs/prm/internal/pkg/mock"
	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/stretchr/testify/assert"
)

func TestPrm_GetStatus(t *testing.T) {
	tests := []struct {
		name       string
		p          *prm.Prm
		wantStatus backend.Status
	}{
		{
			name: "Returns a correct Status object",
			p: &prm.Prm{
				RunningConfig: config.Config{
					PuppetVersion: semver.MustParse("7.15.0"),
					Backend:       DOCKER,
				},
				Backend: &mock.MockBackend{
					StatusIsAvailable:   true,
					StatusMessageString: "Running just fine!",
				},
			},
			wantStatus: backend.Status{
				PuppetVersion: semver.MustParse("7.15.0"),
				Backend:       DOCKER,
				BackendStatus: BackendStatus{
					IsAvailable: true,
					StatusMsg:   "Running just fine!",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotStatus := tt.p.GetStatus(); !reflect.DeepEqual(gotStatus, tt.wantStatus) {
				t.Errorf("Prm.GetStatus() = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}

func TestFormatStatus(t *testing.T) {
	type args struct {
		status     backend.Status
		outputType string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		matches []string
	}{
		{
			name: "human format running backend",
			args: args{
				outputType: "human",
				status: backend.Status{
					PuppetVersion: semver.MustParse("7.15.0"),
					Backend:       DOCKER,
					BackendStatus: BackendStatus{
						IsAvailable: true,
						StatusMsg:   "Running just fine",
					},
				},
			},
			matches: []string{
				"> Puppet version: 7.15.0",
				"> Backend: docker (running)",
			},
		},
		{
			name: "human format errored backend",
			args: args{
				outputType: "human",
				status: backend.Status{
					PuppetVersion: semver.MustParse("7.15.0"),
					Backend:       DOCKER,
					BackendStatus: BackendStatus{
						IsAvailable: false,
						StatusMsg:   "Descriptive error!",
					},
				},
			},
			matches: []string{
				"> Puppet version: 7.15.0",
				"> Backend: docker (error)",
				"> Descriptive error!",
			},
		},
		{
			name: "json format running backend",
			args: args{
				outputType: "json",
				status: backend.Status{
					PuppetVersion: semver.MustParse("7.15.0"),
					Backend:       DOCKER,
					BackendStatus: BackendStatus{
						IsAvailable: true,
						StatusMsg:   "Running just fine",
					},
				},
			},
			matches: []string{
				`"PuppetVersion":"7.15.0"`,
				`"Backend":"docker"`,
				`"IsAvailable":true`,
			},
		},
		{
			name: "json format errored backend",
			args: args{
				outputType: "json",
				status: backend.Status{
					PuppetVersion: semver.MustParse("7.15.0"),
					Backend:       DOCKER,
					BackendStatus: BackendStatus{
						IsAvailable: false,
						StatusMsg:   "Descriptive error!",
					},
				},
			},
			matches: []string{
				`"PuppetVersion":"7.15.0"`,
				`"Backend":"docker"`,
				`"IsAvailable":false`,
				`"StatusMsg":"Descriptive error!"`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStatusMessage, err := backend.FormatStatus(tt.args.status, tt.args.outputType)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, match := range tt.matches {
				assert.Contains(t, gotStatusMessage, match)
			}
		})
	}
}
