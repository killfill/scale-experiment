package appscaler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"scale-experiment/cf"
)

type AppSummary struct {
	Instances int
	CpuAvg    float32
}

type AppScaler struct {
	api           cf.Api
	scaleDownFrom float32
	scaleUpFrom   float32
	minInstances  int
}

func New(a cf.Api) AppScaler {
	return AppScaler{api: a, scaleDownFrom: 0.4, scaleUpFrom: 0.9, minInstances: 2}
}

type AppStat struct {
	State string
	Stats struct {
		DiskQuota int   `json:"disk_quota"`
		MemQuota  int32 `json:"mem_quota"`
		Usage     struct {
			Cpu  float32
			Disk int
			Mem  int
		}
	}
}

func (s *AppScaler) GetSummary(appGuid string) (summary AppSummary, err error) {

	summary = AppSummary{}

	res, err := s.api.Get(fmt.Sprintf("/v2/apps/%s/stats", appGuid))
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("Could not get stats for app. Status: %v\n", res.StatusCode))
		return
	}

	var response map[string]AppStat
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return
	}

	res.Body.Close()

	summary.Instances = len(response)

	var cpuSum float32
	for idx, instance := range response {
		if instance.State != "RUNNING" {
			err = errors.New(fmt.Sprintf("Instance %s is not in RUNNING state: %s", idx, instance.State))
			return
		}
		cpuSum += instance.Stats.Usage.Cpu
	}

	summary.CpuAvg = cpuSum / float32(summary.Instances)
	return
}

func (s *AppScaler) ProposedInstances(summary AppSummary) int {

	//Scale down: one by one
	if summary.Instances > s.minInstances && summary.CpuAvg < s.scaleDownFrom {
		return summary.Instances - 1
	}

	//Scale up:
	if summary.CpuAvg > s.scaleUpFrom {

		//If cpu are 100%, return instances*2

		busyLevel := 100 - summary.CpuAvg
		maxBusy := 100 - s.scaleUpFrom
		return int(float32(summary.Instances)*2.0*busyLevel/maxBusy) + 1
		// return summary.Instances + int(summary.Instances/2) + 1
	}

	return summary.Instances
}

func (s *AppScaler) ScaleAppTo(app string, instances int) error {

	js := []byte(fmt.Sprintf("{\"instances\": %d}", instances))
	b := bytes.NewBuffer(js)

	r, err := s.api.Put(fmt.Sprintf("/v2/apps/%s", app), b)
	if err != nil {
		return err
	}

	if r.StatusCode != 201 {
		return errors.New(fmt.Sprintf("Could not scale to %s. Status: %s\n", instances, r.StatusCode))
	}

	r.Body.Close()
	return nil
}
