# go-led-stream
This is *VERY* rough.

The inspiration for this project came from a $40 LED strand from Amazon.

The basic architure is laid out below:

-----------

#### Arduino
- Loops looking for serial data in
- When finding a packet starting with 0xDEADBEEF read in a packet, check it, and load the data
- Run a custom interpreter with that provided data acting as bytecode
- Use that bytecode to determine how to mutate the LED strand
- Send frames to the strand

#### Go
- Act as the controller for the Arduino
- Build bytecode to run different patters
- Wrap that bytecode in a packet and send it over
