// Package homekit Create a homekit bridge based on hue lights
package homekit

import (
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	log "github.com/sirupsen/logrus"

	"github.com/dj95/huekit/pkg/hue"
)

// StartBridge Create a bridge, required accessories and start the bridge
func StartBridge(pin string, port string, lights []*hue.Light, bridge hue.Bridger) {
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
		hc.Config{Port: port, Pin: pin},
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

		// create the accessory based on the type
		switch light.Type {
		case "On/Off plug-in unit":
			acc = createUnitAccessory(light, bridge)
		case "Dimmable light":
			acc = createDimmableLightAccessory(light, bridge)
		case "Color temperature light":
			acc = createColorTemperatureLightAccessory(light, bridge)
		default:
			acc = nil
			log.Infof("currently type: '%s' is not supported. Please create an issue, if you need support for it: https://github.com/dj95/huekit/issues", light.Type)
		}

		// if the type does not match, continue
		if acc == nil {
			continue
		}

		// create, configure and save the accessory
		accessories = append(accessories, acc)
	}

	// return all configured accessories
	return accessories
}
