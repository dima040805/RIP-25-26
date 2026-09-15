package repository

import (
	"math"
	"testing"
)

func TestCalculatePlanetRadius(t *testing.T) {
	tests := []struct {
		name       string
		starRadius int
		date       string
		shine      float64
		want       float64
		wantErr    bool
	}{
		{name: "one percent dip", starRadius: 100000, date: "2025-10-01", shine: 1, want: 10000},
		{name: "four percent dip", starRadius: 100000, date: "2025-10-01", shine: 4, want: 20000},
		{name: "no dip", starRadius: 100000, date: "2025-10-01", shine: 0, want: 0},
		{name: "empty research date", starRadius: 100000, date: "", shine: 1, wantErr: true},
		{name: "negative dip", starRadius: 100000, date: "2025-10-01", shine: -0.5, wantErr: true},
		{name: "dip above seven percent", starRadius: 100000, date: "2025-10-01", shine: 7.5, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculatePlanetRadius(tt.starRadius, tt.date, tt.shine)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got radius %v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(got-tt.want) > 1e-6 {
				t.Fatalf("radius = %v, want %v", got, tt.want)
			}
		})
	}
}
