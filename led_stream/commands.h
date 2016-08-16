#ifndef __COMMANDS_H__
#define __COMMANDS_H__

#include <stdint.h>
#include <stddef.h>

#include "FastLED.h"
#include "config.h"

FASTLED_USING_NAMESPACE

extern CRGB leds[];
extern uint8_t gFlags;
extern uint8_t gHue;
extern long cmd_start_millis;

struct COMMAND {
  uint8_t index;
  uint16_t duration;
  uint8_t params_size;
};

typedef void (*funcList[])(COMMAND *c, void *ptr);


void off(COMMAND *c, void *ptr);


struct FILL_SOLID_PARAMS {
  uint8_t color[3];
};
void fill_solid_rgb(COMMAND *c, void *ptr);


struct SET_FLAGS_PARAMS {
  uint8_t flags = 0;
};
void set_flags(COMMAND *c, void *ptr);


struct RAINBOW_PARAMS {
  uint8_t delta = 1;
};
void rainbow(COMMAND *c, void *ptr);


struct FADE_PARAMS {
  uint8_t start_color[3];
  uint8_t end_color[3];
};
void fade_rgb(COMMAND *c, void *ptr);

struct FADE_TO_BLACK_PARAMS {
  uint8_t fadeBy = 1;
};
void fade_to_black(COMMAND *c, void *ptr);

struct CONFETTI_PARAMS {
  uint8_t count = 1;
  uint8_t fadeBy = 10;
};
void confetti(COMMAND *c, void *ptr);
#endif
