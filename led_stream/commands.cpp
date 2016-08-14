#include "commands.h"

FASTLED_USING_NAMESPACE

void set_flags(COMMAND *c, void *ptr) {
  SET_FLAGS_PARAMS *params = (ptr != NULL) ? (struct SET_FLAGS_PARAMS*)ptr : new SET_FLAGS_PARAMS{};

  gFlags = params->flags;
}

void rainbow(COMMAND *c, void *ptr) {
  RAINBOW_PARAMS *params = (ptr != NULL) ? (struct RAINBOW_PARAMS*)ptr : new RAINBOW_PARAMS{};

  fill_rainbow( leds, PIXEL_NUM, gHue, params->delta );
}

void fade_rgb(COMMAND *c, void *ptr) {
  FADE_PARAMS *params = (ptr != NULL) ? (struct FADE_PARAMS*)ptr : new FADE_PARAMS{};

  CRGB c1 = CRGB(params->start_color[0],
      params->start_color[1],
      params->start_color[2]);
  CRGB c2 = CRGB(params->end_color[0],
      params->end_color[1],
      params->end_color[2]);

  double ratio = (millis() - cmd_start_millis) / double(c->duration);

  fill_solid(leds, PIXEL_NUM,
      blend(c1, c2, min(ratio, 1.0) * 255));
}
