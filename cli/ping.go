package cli

import "github.com/NubeIO/nubeio-rubix-lib-models-go/dto"

func (inst *Client) Ping() (deviceInfo *dto.DeviceInfo, pingable bool, isValidToken bool) {
	url := "/api/system/device"
	resp, err := inst.Rest.R().
		SetResult(&dto.DeviceInfo{}).
		Get(url)
	if err != nil {
		return nil, false, false
	}
	if resp.StatusCode() == 401 {
		return nil, true, false
	}
	deviceInfoResult := resp.Result().(*dto.DeviceInfo)
	return deviceInfoResult, true, true
}
