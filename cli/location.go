package cli

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/model"
	"github.com/NubeIO/rubix-edge-wires/clients/nresty"
)

func (inst *Client) GetLocations() ([]*model.Location, error) {
	url := fmt.Sprintf("/api/locations?with_groups=true&with_hosts=true")
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&[]*model.Location{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	var out []*model.Location
	out = *resp.Result().(*[]*model.Location)
	return out, nil
}
