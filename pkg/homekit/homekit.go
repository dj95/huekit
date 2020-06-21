// Package homekit Create a homekit bridge based on hue lights
package homekit

import (
	"math"
	"strconv"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	log "github.com/sirupsen/logrus"

	"github.com/dj95/huekit/pkg/hue"
)

// StartBridge Create a bridge, required accessories and start the bridge
func StartBridge(pin string, lights []*hue.Light, bridge hue.Bridger) {
	// create the bridge accessory
	bridgeAccessory := accessory.NewBridge(accessory.Info{
		ID:   1,
		Name: "HueKit Bridge",
	})

	// create the lights based on the hue lights without a matching
	// modelID
	accessories := configureLights(lights, bridge)

	// create the ip transport, that publishes homekit functionality
	// and acts as the bridge
	t, err := hc.NewIPTransport(
		hc.Config{Pin: pin},
		bridgeAccessory.Accessory,
		accessories[:]...,
	)

	// error handling
	if err != nil {
		log.Fatal(err)
	}

	// enable graceful exit for the homekit bridge
	hc.OnTermination(func() {
		<-t.Stop()
	})

	// start the communication
	t.Start()
}

func configureLights(lights []*hue.Light, bridge hue.Bridger) []*accessory.Accessory {
	// initialize the accessories
	var accessories []*accessory.Accessory

	// iterate through all hue lights
	for _, light := range lights {
		// check, if the light has a model is from hue
		if hue.ModelIDIsFromHue(light.ModelID) {
			continue
		}

		var acc *accessory.Accessory

		switch light.Type {
		case "On/Off plug-in unit":
			acc = createUnitAccessory(light, bridge)
		case "Dimmable light":
			acc = createDimmableLightAccessory(light, bridge)
		case "Extended color light":
			acc = createExtendedColorLightAccessory(light, bridge)
		default:
			acc = nil
		}

		if acc == nil {
			continue
		}

		// create, configure and save the accessory
		accessories = append(accessories, acc)
	}

	// return all configured accessories
	return accessories
}

func createExtendedColorLightAccessory(light *hue.Light, bridge hue.Bridger) *accessory.Accessory {
	log.Debugf("creating accessory for: %s - %s", light.ID, light.Name)

	// convert the id to an int. As hue's ids are integers, omit the error
	// handling
	id, _ := strconv.Atoi(light.ID)

	// create the lightbulb accessory
	ac := accessory.NewColoredLightbulb(accessory.Info{
		ID:               uint64(id + 1),
		Name:             light.Name,
		Model:            light.ModelID,
		FirmwareRevision: light.SoftwareVersion,
	})

	// configure what do to, when the home app changes the state
	// of the light
	ac.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
		// send a toggle request
		err := bridge.LightUpdateState(light, &hue.State{On: on})

		// if an error occurred...
		if err != nil {
			// ...log it
			log.WithFields(log.Fields{
				"id":    uint64(id + 1),
				"name":  light.Name,
				"state": on,
				"on":    "ValueRemoteUpdate",
			}).Errorf("%v", err)
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

	// configure what do to, when the home app changes the brightness
	// of the light
	ac.Lightbulb.Brightness.OnValueRemoteUpdate(func(bri int) {
		bri = int(math.Floor(float64(bri) * 2.54))

		// send a toggle request
		err := bridge.LightUpdateState(light, &hue.State{On: true, Brightness: bri})

		// if an error occurred...
		if err != nil {
			// ...log it
			log.WithFields(log.Fields{
				"id":   uint64(id + 1),
				"name": light.Name,
				"bri":  bri,
				"on":   "ValueRemoteUpdate",
			}).Errorf("%v", err)
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
		return int(math.Floor(float64(l.State.Brightness) / 2.54))
	})

	// return the configured accessory
	return ac.Accessory
}

func createDimmableLightAccessory(light *hue.Light, bridge hue.Bridger) *accessory.Accessory {
	log.Debugf("creating accessory for: %s - %s", light.ID, light.Name)

	// convert the id to an int. As hue's ids are integers, omit the error
	// handling
	id, _ := strconv.Atoi(light.ID)

	// create the lightbulb accessory
	ac := NewDimmableLightbulb(accessory.Info{
		ID:               uint64(id + 1),
		Name:             light.Name,
		Model:            light.ModelID,
		FirmwareRevision: light.SoftwareVersion,
	})

	// configure what do to, when the home app changes the state
	// of the light
	ac.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
		// send a toggle request
		err := bridge.LightUpdateState(light, &hue.State{On: on})

		// if an error occurred...
		if err != nil {
			// ...log it
			log.WithFields(log.Fields{
				"id":    uint64(id + 1),
				"name":  light.Name,
				"state": on,
				"on":    "ValueRemoteUpdate",
			}).Errorf("%v", err)
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

	// configure what do to, when the home app changes the brightness
	// of the light
	ac.Lightbulb.Brightness.OnValueRemoteUpdate(func(bri int) {
		bri = int(math.Floor(float64(bri) * 2.54))

		// send a toggle request
		err := bridge.LightUpdateState(light, &hue.State{On: true, Brightness: bri})

		// if an error occurred...
		if err != nil {
			// ...log it
			log.WithFields(log.Fields{
				"id":   uint64(id + 1),
				"name": light.Name,
				"bri":  bri,
				"on":   "ValueRemoteUpdate",
			}).Errorf("%v", err)
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
		return int(math.Floor(float64(l.State.Brightness) / 2.54))
	})

	// return the configured accessory
	return ac.Accessory
}

func createUnitAccessory(light *hue.Light, bridge hue.Bridger) *accessory.Accessory {
	log.Debugf("creating accessory for: %s - %s", light.ID, light.Name)

	// convert the id to an int. As hue's ids are integers, omit the error
	// handling
	id, _ := strconv.Atoi(light.ID)

	// create the lightbulb accessory
	ac := accessory.NewLightbulb(accessory.Info{
		ID:               uint64(id + 1),
		Name:             light.Name,
		Model:            light.ModelID,
		FirmwareRevision: light.SoftwareVersion,
	})

	// configure what do to, when the home app changes the state
	// of the light
	ac.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
		// send a toggle request
		err := bridge.LightUpdateState(light, &hue.State{On: on})

		// if an error occurred...
		if err != nil {
			// ...log it
			log.WithFields(log.Fields{
				"id":    uint64(id + 1),
				"name":  light.Name,
				"state": on,
				"on":    "ValueRemoteUpdate",
			}).Errorf("%v", err)
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
