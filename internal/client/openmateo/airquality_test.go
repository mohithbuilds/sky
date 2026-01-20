package openmateo

import (
	"reflect"
	"testing"
)

func TestAirQualityClient_GetAirQuality(t *testing.T) {
	type fields struct {
		baseClient *baseClient
		BaseURL    string
	}
	type args struct {
		longitude                  float64
		latitude                   float64
		hourlyAirQualityParameters []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Result
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aqc := &AirQualityClient{
				baseClient: tt.fields.baseClient,
				BaseURL:    tt.fields.BaseURL,
			}
			got, err := aqc.GetAirQuality(
				tt.args.longitude,
				tt.args.latitude,
				tt.args.hourlyAirQualityParameters,
			)
			if (err != nil) != tt.wantErr {
				t.Fatalf("AirQualityClient.GetAirQuality() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AirQualityClient.GetAirQuality() = %v, want %v", got, tt.want)
			}
		})
	}
}
