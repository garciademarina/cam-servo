package main

import (
	"log"

	"github.com/garciademarina/cam-servo/server"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {

	adaptor := raspi.NewAdaptor()
	servoHorizontal := gpio.NewServoDriver(adaptor, "12")
	servoVertical := gpio.NewServoDriver(adaptor, "16")

	horizontalMoveAngle := make(chan uint8, 100)
	verticalMoveAngle := make(chan uint8, 100)

	// start server
	s := server.New(horizontalMoveAngle, verticalMoveAngle)
	go s.Run()

	work := func() {
		go func() {
			for {
				select {
				case angle := <-horizontalMoveAngle:
					log.Printf("Moving horizontal servo to %d", angle)
					servoHorizontal.Move(angle)
				case angle := <-verticalMoveAngle:
					log.Printf("Moving vertical servo to %d", angle)
					servoVertical.Move(angle)
				}
			}
		}()
	}

	robot := gobot.NewRobot("Servo robot",
		[]gobot.Connection{adaptor},
		[]gobot.Device{servoHorizontal, servoVertical},
		work,
	)

	robot.Start()

}
