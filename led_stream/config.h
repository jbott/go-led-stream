#ifndef __CONFIG_H__
#define __CONFIG_H__

// Number of pixels in strand
#define PIXEL_NUM 240

// Pin strand is connected on
#define PIXEL_PIN 6

// Pixel led type
#define PIXEL_LED_TYPE WS2812B

// Color order
#define PIXEL_COLOR_ORDER GRB

// Pixel brightness
#define PIXEL_BRIGHTNESS 128

// Framerate
#define PIXEL_FRAMERATE 60

// Serial baud rate
#define SERIAL_BAUD 57600

// Ethernet mac address
#define ETH_MAC_ADDR { 0xDE, 0xED, 0xBA, 0xFE, 0xFE, 0xED }

// Static ip address
#define IP_STATIC_ADDR 192, 168, 1, 76

// MQTT Server IP
#define IP_MQTT_ADDR 192, 168, 1, 21

// MQTT Server port
#define IP_MQTT_PORT 1883

#endif
