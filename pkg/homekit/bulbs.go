package homekit

import (
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
)

// DimmableLightbulb Represent a light bulb, that is only dimmable
type DimmableLightbulb struct {
	*accessory.Accessory
	Lightbulb *DimmableLightbulbService
}

// NewDimmableLightbulb Create a new accessory for the dimmable lightbulb
func NewDimmableLightbulb(info accessory.Info) *DimmableLightbulb {
	// initialize the accessory
	acc := DimmableLightbulb{}

	// set the base accessory with given information
	acc.Accessory = accessory.New(info, accessory.TypeLightbulb)

	// set the dimmable service
	acc.Lightbulb = newDimmableLightbulbService()

	// register all services to the accessory itself
	acc.AddService(acc.Lightbulb.Service)

	return &acc
}

// DimmableLightbulbService Represent the services behind the dimmable light
// bulb, e.g. power state and brightness
type DimmableLightbulbService struct {
	*service.Service

	On         *characteristic.On
	Brightness *characteristic.Brightness
}

func newDimmableLightbulbService() *DimmableLightbulbService {
	svc := DimmableLightbulbService{}
	svc.Service = service.New(service.TypeLightbulb)

	svc.On = characteristic.NewOn()
	svc.AddCharacteristic(svc.On.Characteristic)

	svc.Brightness = characteristic.NewBrightness()
	svc.AddCharacteristic(svc.Brightness.Characteristic)

	return &svc
}
