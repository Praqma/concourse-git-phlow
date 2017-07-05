package mwriter

import (
	"github.com/praqma/concourse-git-phlow/models"
	"gopkg.in/zorkian/go-datadog-api.v2"
	"time"
)

const (
	Error   = "error"
	Warning = "warning"
	Info    = "info"
)

//DataDog...
type DataDog struct {
	Name   string
	Active bool
	Client *datadog.Client
}

//SpawnCerberus...
func SpawnCerberus(source models.Source) *DataDog {
	var active = false
	if source.DataDogAppKey != "" && source.DataDogApiKey != "" {
		active = true
	}

	c := datadog.NewClient(source.DataDogApiKey, source.DataDogAppKey)
	return &DataDog{Client: c, Name: source.DataDogMetricName, Active: active}
}

//BarkEvent...
func (d *DataDog) BarkEvent(output, alertType string) error {

	if d.Active {
		e := datadog.Event{Title: &d.Name, Text: &output, AlertType: &alertType}

		_, err := d.Client.PostEvent(&e)
		if err != nil {
			return err
		}
	}
	return nil
}

//WufMetric...
func (d *DataDog) WufMetric() error {

	if d.Active {
		var dType, unit, metric = "c", "s", "concourse_tollgate"

		dp := []datadog.DataPoint{datadog.DataPoint{float64(time.Now().Unix()), 1}}
		m := datadog.Metric{Metric: &metric, Type: &dType, Host: &d.Name, Points: dp, Unit: &unit}

		met := []datadog.Metric{}
		met = append(met, m)

		err := d.Client.PostMetrics(met)
		if err != nil {
			return err
		}
	}

	return nil
}
