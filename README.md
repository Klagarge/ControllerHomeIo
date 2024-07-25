# Home IO alarm Controller
This program is a very simple controller for the alarm system on Home I/O.
It checks the state of the main door, the garage door and the main motion sensor.
If any of these sensors are triggered, the program will send an alarm message.

## Prerequisites
- Home I/O
- The modbus interface for Home I/O : [Modbus2HomeIO](https://github.com/Klagarge/Modbus2HomeIO)

## Usage
1. Run the Home I/O simulation
2. Run the Modbus2HomeIO program
3. Modifie IP on `main.go`
4. Run the program

## Authors
- **Rémi Heredero** - _Initial work_ - [Klagarge](https://github.com/Klagarge)