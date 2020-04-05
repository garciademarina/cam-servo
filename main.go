package main

import (
	"github.com/garciademarina/cam-servo/server"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

// 8 min
// 52 max

func main() {

	adaptor := raspi.NewAdaptor()
	servo := gpio.NewServoDriver(adaptor, "12")
	servo.Start()
	servo.CurrentAngle = 0

	// start server
	s := server.New(servo)
	go s.Run()

	work := func() {
		// log.Printf("Moving servo Max")
		// servo.Move(52)
		// time.Sleep(3 * time.Second)
		// var angle uint8 = 52

		// gobot.Every(1*time.Second, func() {
		// 	angle = uint8(angle - 2)
		// 	log.Printf("Moving servo to %d", angle)
		// 	servo.Move(angle)
		// })

		// log.Printf("Moving servo CENTER")
		// servo.Center()
		// time.Sleep(3 * time.Second)

		// log.Printf("Moving servo Min")
		// servo.Move(0)
		// time.Sleep(3 * time.Second)

	}

	robot := gobot.NewRobot("Servo robot",
		[]gobot.Connection{adaptor},
		[]gobot.Device{servo},
		work,
	)

	robot.Start()

}
