package homekit

import (
	"math"
	"strconv"

	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
	log "github.com/sirupsen/logrus"

	"github.com/dj95/huekit/pkg/hue"
)

// ColorTemperatureLight Represent a light, that is dimmable and can white colors
type ColorTemperatureLight struct {
	*accessory.Accessory
	Lightbulb *ColorTemperatureLightService
}

// NewColorTemperatureLight Create a new accessory for the CCT light
func NewColorTemperatureLight(info accessory.Info) *ColorTemperatureLight {
	// initialize the accessory
	acc := ColorTemperatureLight{}

	// set the base accessory with given information
	acc.Accessory = accessory.New(info, accessory.TypeLightbulb)

	// set the dimmable service
	acc.Lightbulb = newColorTemperatureLightService()

	// register all services to the accessory itself
	acc.AddService(acc.Lightbulb.Service)

	return &acc
}

// DimmableLightbulbService Represent the services behind the dimmable light
// bulb, e.g. power state and brightness
type ColorTemperatureLightService struct {
	*service.Service

	On               *characteristic.On
	Brightness       *characteristic.Brightness
	ColorTemperature *characteristic.ColorTemperature
}

func newColorTemperatureLightService() *ColorTemperatureLightService {
	// instantiate the service and register it
	svc := ColorTemperatureLightService{}
	svc.Service = service.New(service.TypeLightbulb)

	// register the On characteristic for the power state
	svc.On = characteristic.NewOn()
	svc.AddCharacteristic(svc.On.Characteristic)

	// register the brightness characteristic
	svc.Brightness = characteristic.NewBrightness()
	svc.AddCharacteristic(svc.Brightness.Characteristic)

	// register the color temperature characteristic
	svc.ColorTemperature = characteristic.NewColorTemperature()
	svc.AddCharacteristic(svc.ColorTemperature.Characteristic)

	// return the custom service
	return &svc
}

func createColorTemperatureLightAccessory(light *hue.Light, bridge hue.Bridger) *accessory.Accessory {
	log.Debugf("creating color temperature light accessory for: %s - %s", light.ID, light.Name)

	// convert the id to an int. As hue's ids are integers, omit the error
	// handling
	id, _ := strconv.Atoi(light.ID)

	// create the lightbulb accessory
	ac := NewColorTemperatureLight(accessory.Info{
		ID:               uint64(id + 1),
		Name:             light.Name,
		Model:            light.ModelID,
		Manufacturer:     light.ManufacturerName,
		FirmwareRevision: light.SoftwareVersion,
	})

	//
	// Power State
	//

	// configure what do to, when the home app changes the state
	// of the light
	ac.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
		// send a toggle request
		err := bridge.LightUpdateState(light, &hue.State{On: on})

		log.WithFields(log.Fields{
			"id":   uint64(id + 1),
			"name": light.Name,
			"type": light.Type,
		}).Debugf("trigger state: %t", on)

		// if an error occurred...
		if err != nil {
			// ...log it
			log.WithFields(log.Fields{
				"id":    uint64(id + 1),
				"name":  light.Name,
				"state": on,
				"on":    "on",
			}).Errorf("%s", err.Error())
		}
	})

	// configure what to do, when the home app fetches the state
	// of the light
	ac.Lightbulb.On.OnValueRemoteGet(func() bool {
		// refetch the light information based on the id
		l, err := bridge.Light(light.ID)

		// return, that the light is of, if an error
		// occurred
		if err != nil {
			return false
		}

		// otherwise return the correct state
		return l.State.On
	})

	//
	// Brightness
	//

	// configure what do to, when the home app changes the brightness
	// of the light
	ac.Lightbulb.Brightness.OnValueRemoteUpdate(func(bri int) {
		bri = int(math.Floor(float64(bri)*254) / 100)

		// send a toggle request
		err := bridge.LightUpdateState(light, &hue.State{On: true, Brightness: bri})

		log.WithFields(log.Fields{
			"id":   uint64(id + 1),
			"name": light.Name,
			"type": light.Type,
		}).Debugf("change brightness: %d", bri)

		// if an error occurred...
		if err != nil {
			// ...log it
			log.WithFields(log.Fields{
				"id":   uint64(id + 1),
				"name": light.Name,
				"bri":  bri,
				"on":   "brightness",
			}).Errorf("%s", err.Error())
		}
	})

	// configure what to do, when the home app fetches the brightness
	// of the light
	ac.Lightbulb.Brightness.OnValueRemoteGet(func() int {
		// refetch the light information based on the id
		l, err := bridge.Light(light.ID)

		// return, that the light is of, if an error
		// occurred
		if err != nil {
			return 0
		}

		// otherwise return the correct state
		return int(math.Floor(float64(l.State.Brightness*100) / 254))
	})

	//
	// ColorTemperature
	//

	// configure what do to, when the home app changes the color temperature
	// of the light
	// homekit range for color temperature 50 - 400 [mired]
	ac.Lightbulb.ColorTemperature.OnValueRemoteUpdate(func(colorTemperature int) {
		colorTemperature = int(math.Min(400, math.Max(50, float64(colorTemperature))))

		// send a toggle request
		err := bridge.LightUpdateState(light, &hue.State{On: true, ColorTemperature: colorTemperature})

		log.WithFields(log.Fields{
			"id":   uint64(id + 1),
			"name": light.Name,
			"type": light.Type,
		}).Debugf("change color-temperature: %d", colorTemperature)

		// if an error occurred...
		if err != nil {
			// ...log it
			log.WithFields(log.Fields{
				"id":               uint64(id + 1),
				"name":             light.Name,
				"colorTemperature": colorTemperature,
				"on":               "color-temperature",
			}).Errorf("%s", err.Error())
		}
	})

	// configure what to do, when the home app fetches the color temperature
	// of the light
	// hue range for color temperature 153 - 500 [mired]
	ac.Lightbulb.Brightness.OnValueRemoteGet(func() int {
		// refetch the light information based on the id
		l, err := bridge.Light(light.ID)

		// return, that the light is of, if an error
		// occurred
		if err != nil {
			return 0
		}

		// otherwise return the correct state
		return int(math.Min(500, math.Max(153, float64(l.State.ColorTemperature))))
	})

	// return the configured accessory
	return ac.Accessory
}
