#include "commands.h"

FASTLED_USING_NAMESPACE

void off(COMMAND *c, void *ptr) {
  fill_solid(leds, PIXEL_NUM, CRGB(0x0));
}

const FILL_SOLID_PARAMS default_fill_solid_params = FILL_SOLID_PARAMS{};
void fill_solid_rgb(COMMAND *c, void *ptr) {
  const FILL_SOLID_PARAMS *params = (ptr != NULL) ? (struct FILL_SOLID_PARAMS*)ptr : &default_fill_solid_params;

  fill_solid(leds, PIXEL_NUM, CRGB(params->color[0], params->color[1], params->color[2]));
}

const SET_FLAGS_PARAMS default_set_flags_params = SET_FLAGS_PARAMS{};
void set_flags(COMMAND *c, void *ptr) {
  const SET_FLAGS_PARAMS *params = (ptr != NULL) ? (struct SET_FLAGS_PARAMS*)ptr : &default_set_flags_params;

  gFlags = params->flags;
}

const RAINBOW_PARAMS default_rainbow_params = RAINBOW_PARAMS{};
void rainbow(COMMAND *c, void *ptr) {
  const RAINBOW_PARAMS *params = (ptr != NULL) ? (struct RAINBOW_PARAMS*)ptr : &default_rainbow_params;

  fill_rainbow( leds, PIXEL_NUM, gHue, params->delta );
}

const FADE_PARAMS default_fade_params = FADE_PARAMS{};
void fade_rgb(COMMAND *c, void *ptr) {
  const FADE_PARAMS *params = (ptr != NULL) ? (struct FADE_PARAMS*)ptr : &default_fade_params;

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

const FADE_TO_BLACK_PARAMS default_fade_to_black_params = FADE_TO_BLACK_PARAMS{};
void fade_to_black(COMMAND *c, void *ptr) {
  const FADE_TO_BLACK_PARAMS *params = (ptr != NULL) ? (struct FADE_TO_BLACK_PARAMS*)ptr : &default_fade_to_black_params;
  fadeToBlackBy( leds, PIXEL_NUM, params->fadeBy );
}

const CONFETTI_PARAMS default_confetti_params = CONFETTI_PARAMS{};
void confetti(COMMAND *c, void *ptr) {
  const CONFETTI_PARAMS *params = (ptr != NULL) ? (struct CONFETTI_PARAMS*)ptr : &default_confetti_params;

  fadeToBlackBy( leds, PIXEL_NUM, params->fadeBy);
  for (int i = 0; i < params->count; i++) {
    int pos = random16(PIXEL_NUM);
    leds[pos] += CHSV( gHue + random8(64), 200, 255);
  }
}
