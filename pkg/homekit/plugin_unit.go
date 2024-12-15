package homekit

import (
	"strconv"

	"github.com/brutella/hc/accessory"
	log "github.com/sirupsen/logrus"

	"github.com/dj95/huekit/pkg/hue"
)

func createUnitAccessory(light *hue.Light, bridge hue.Bridger) *accessory.Accessory {
	log.Debugf("creating lightbulb accessory for: %s - %s", light.ID, light.Name)

	// convert the id to an int. As hue's ids are integers, omit the error
	// handling
	id, _ := strconv.Atoi(light.ID)

	// create the lightbulb accessory
	ac := accessory.NewLightbulb(accessory.Info{
		ID:               uint64(id + 1), // #nosec G115 IDs will always be smaller
		Name:             light.Name,
		Model:            light.ModelID,
		FirmwareRevision: light.SoftwareVersion,
	})

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

	// return the configured accessory
	return ac.Accessory
}
