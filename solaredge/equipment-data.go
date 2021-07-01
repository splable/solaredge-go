package solaredge

import (
	"errors"
	"fmt"
	"time"
)

type EquipmentDataResponse struct {
	Data EquipmentData `json:"data"`
}

type EquipmentData struct {
	Count       int64         `json:"count"`
	Telemetries []Telemetries `json:"telemetries"`
}

type Telemetries struct {
	Date        DateTime `json:"date"`
	Temperature float64  `json:"temperature"`
}

type EquipmentDataRequest struct {
	StartTime DateTime `url:"startTime"`
	EndTime   DateTime `url:"endTime"`
}

func (s *SiteService) InverterData(siteId int64, inverterId string, request EquipmentDataRequest) (EquipmentData, error) {
	// Ensure start and end are defined
	if request.EndTime.IsZero() || request.StartTime.IsZero() {
		return EquipmentData{}, errors.New("startTime and endTime are required")
	}
	// Ensure delta between start and end is one week or less
	if request.EndTime.Sub(request.StartTime.Time) > time.Hour*24*7 {
		return EquipmentData{}, errors.New("duration between StartTime and EndTime should be less than one week")
	}

	u, err := addOptions(fmt.Sprintf("/equipment/%d/%s/data", siteId, inverterId), request)
	if err != nil {
		return EquipmentData{}, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return EquipmentData{}, err
	}
	var equipmentDataResponse EquipmentDataResponse
	_, err = s.client.do(req, &equipmentDataResponse)
	return equipmentDataResponse.Data, err
}
