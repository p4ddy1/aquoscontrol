# Aquos Control

This is a tool to control a Sharp Aquos TV via the serial port over network. It creates a http server listening for requests containing commands and sends them to the TV. Useful to integrate the device in a home automation system like [Domoticz](http://www.domoticz.com/). For further information on which serial commands to use please take a look at the user manual of your model.

## Usage
### Installing
```
# go get github.com/p4ddy1/aquoscontrol
# $GOPATH/bin/aquoscontrol -config=/path/to/config.toml
```
### Example config.toml
```
# Logging to console
logging = false

# Configure the http server
[server]
address = ":8000" 
path = "/ctl"       

# Configure the serial port
[serial]
port = "/dev/ttyUSB1"
baud = 9600

# Define commands
[commands]

    [commands.poweron]
    command = "POWR1" # Command sent to the TV 
    hasParameter = false # Wont accept parameter in GET request

    [commands.poweroff]
    command = "POWR0"
    hasParameter = false
    
    [commands.input]
    command = "IAVD"
    hasParameter = true # Parameter value in GET request will be added
```

### Example requests
**Poweron**
```
http://192.168.178.20:8000/ctl?cmd=poweron
```

**Poweroff**
```
http://192.168.178.20:8000/ctl?cmd=poweroff
```

**Switch input to HDMI 1**
```
http://192.168.178.20:8000/ctl?cmd=input&value=4
```