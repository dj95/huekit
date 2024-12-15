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
	// instantiate the service and register it
	svc := DimmableLightbulbService{}
	svc.Service = service.New(service.TypeLightbulb)

	// register the On characteristic for the power state
	svc.On = characteristic.NewOn()
	svc.AddCharacteristic(svc.On.Characteristic)

	// register the brightness characteristic
	svc.Brightness = characteristic.NewBrightness()
	svc.AddCharacteristic(svc.Brightness.Characteristic)

	// return the custom service
	return &svc
}

func createDimmableLightAccessory(light *hue.Light, bridge hue.Bridger) *accessory.Accessory {
	log.Debugf("creating dimmable light accessory for: %s - %s", light.ID, light.Name)

	// convert the id to an int. As hue's ids are integers, omit the error
	// handling
	id, _ := strconv.Atoi(light.ID)

	// create the lightbulb accessory
	ac := NewDimmableLightbulb(accessory.Info{
		ID:               uint64(id + 1), // #nosec G115 IDs will always be smaller
		Name:             light.Name,
		Model:            light.ModelID,
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
			"id":   id + 1,
			"name": light.Name,
			"type": light.Type,
		}).Debugf("trigger state: %t", on)

		// if an error occurred...
		if err != nil {
			// ...log it
			log.WithFields(log.Fields{
				"id":    id + 1,
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
			"id":   id + 1,
			"name": light.Name,
			"type": light.Type,
		}).Debugf("change brightness: %d", bri)

		// if an error occurred...
		if err != nil {
			// ...log it
			log.WithFields(log.Fields{
				"id":   id + 1,
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

	// return the configured accessory
	return ac.Accessory
}
