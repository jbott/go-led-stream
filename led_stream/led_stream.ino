#include "FastLED.h"
#include "commands.h"
#include "config.h"

FASTLED_USING_NAMESPACE

#if FASTLED_VERSION < 3001000
#error "Requires FastLED 3.1 or later; check github for latest code."
#endif


// Helper macros
#define INTER_FRAME_MILLIS int(1000 / PIXEL_FRAMERATE)

#define CONVERT_16_TO_8_INLINE(value) value & 0xff, (value >> 8) & 0xff

#define IS_FLAG_SET(f) (gFlags & f)

#define FLAG_NONE         (0)
#define FLAG_PAUSE_GHUE   (1 << 0)
#define FLAG_REVERSE_GHUE (1 << 1)

CRGB leds[PIXEL_NUM];
uint8_t gFlags = FLAG_NONE;
uint8_t gHue = 0; // rotating "base color" used by many of the patterns

uint8_t command_buffer[256];
int command_buffer_size = 0;
void *exec_pos = NULL;
long cmd_start_millis = 0;
long last_frame_millis = 0;

enum FUNC_INDEX : uint8_t {
  SET_FLAGS = 0,
  OFF,
  FILL_SOLID_RGB,
  RAINBOW,
  FADE_RGB,
  FADE_TO_BLACK,
  CONFETTI,
  COMMAND_LENGTH
};

funcList gFuncs = {
  set_flags,
  off,
  fill_solid_rgb,
  rainbow,
  fade_rgb,
  fade_to_black,
  confetti
};

void setup() {
  delay(3000);

  Serial.begin(SERIAL_BAUD);
  Serial.setTimeout(SERIAL_TIMEOUT);

  FastLED.addLeds<PIXEL_LED_TYPE,PIXEL_PIN,PIXEL_COLOR_ORDER>(leds, PIXEL_NUM).setCorrection(TypicalLEDStrip);

  Serial.println("led_stream (" __DATE__ " " __TIME__ ")");
}

void loop() {
	if (command_buffer_size == 0) return;

  if (exec_pos == NULL || exec_pos >= command_buffer + command_buffer_size) {
    // Set starting configuration
    FastLED.setBrightness(PIXEL_BRIGHTNESS);
    exec_pos = command_buffer;
    gFlags = FLAG_NONE;
  }

  COMMAND *c = (struct COMMAND*) exec_pos;

  if (cmd_start_millis == 0) {
    // First time running command
    cmd_start_millis = millis();

    // Debug
    Serial.println("COMMAND");
    Serial.println(c->index);
    Serial.println(c->duration);
    Serial.println(c->params_size);
  }

  if (c->params_size == 0) {
    gFuncs[c->index](c, NULL);
  } else {
    gFuncs[c->index](c, exec_pos + sizeof(COMMAND));
  }

  if ((millis() - cmd_start_millis) >= c->duration) {
    exec_pos += sizeof(COMMAND);
    exec_pos += c->params_size;
    cmd_start_millis = 0;

    if (c->duration == 0) return; // Break out early if a zero length command
  }

  while (millis() - last_frame_millis < INTER_FRAME_MILLIS) {} // Spin until time to write frame
  last_frame_millis = millis();

  // Write out
  FastLED.show();

  // Other stuff we need to do
  EVERY_N_MILLISECONDS( 20 ) { if (!IS_FLAG_SET(FLAG_PAUSE_GHUE)) gHue += IS_FLAG_SET(FLAG_REVERSE_GHUE) ? -1 : 1; } // slowly cycle the "base color" through the rainbow
}

uint8_t header[] = { 0xDE, 0xAD, 0xBE, 0xEF };

void serialEvent() {
  // Check for header
  uint8_t read[6];
  Serial.readBytes(read, 6);
  for (int i = 0; i < sizeof(header); i++) {
    if (read[i] != header[i]) {
      return;
    }
  }
  Serial.println("HEADER_FOUND");

  // Header found!
  int packet_type = read[4];
  int packet_length = read[5];
  Serial.print("TYPE: "); Serial.println(packet_type, HEX);
  Serial.print("LENGTH: "); Serial.println(packet_length, HEX);

  uint8_t data[packet_length];
  int n = Serial.readBytes(data, packet_length);
  if (n != packet_length) {
    while (Serial.available()) {
      Serial.read();
    }
    return;
  }
  for (int i = 0; i < sizeof(data); i++) {
    Serial.print(data[i], HEX);
    if (i + 1 != sizeof(data))
      Serial.print(" ");
  }
  Serial.println();

  // Checksum
  uint8_t checksum[4];
  Serial.readBytes(checksum, 4);
  long given_crc = (long)checksum[0] | ((long)checksum[1] << 8) | ((long)checksum[2] << 16) | ((long)checksum[3] << 24);
  long our_crc = crc32(data, packet_length);
  if (given_crc != our_crc) {
    Serial.print("GIVEN: "); Serial.println(given_crc, HEX);
    Serial.print("OURS:  "); Serial.println(our_crc, HEX);
    while (Serial.available()) {
      Serial.read();
    }
    return;
  }

  memcpy(command_buffer, data, packet_length);
  command_buffer_size = packet_length;
  exec_pos = NULL;

  while (Serial.available()) {
    Serial.read();
  }
}

const static PROGMEM prog_uint32_t crc_table[16] = {
    0x00000000, 0x1db71064, 0x3b6e20c8, 0x26d930ac,
    0x76dc4190, 0x6b6b51f4, 0x4db26158, 0x5005713c,
    0xedb88320, 0xf00f9344, 0xd6d6a3e8, 0xcb61b38c,
    0x9b64c2b0, 0x86d3d2d4, 0xa00ae278, 0xbdbdf21c
};

unsigned long crc_update(unsigned long crc, byte data)
{
    byte tbl_idx;
    tbl_idx = crc ^ (data >> (0 * 4));
    crc = pgm_read_dword_near(crc_table + (tbl_idx & 0x0f)) ^ (crc >> 4);
    tbl_idx = crc ^ (data >> (1 * 4));
    crc = pgm_read_dword_near(crc_table + (tbl_idx & 0x0f)) ^ (crc >> 4);
    return crc;
}

unsigned long crc32(uint8_t *data, size_t size)
{
  unsigned long crc = ~0L;
  for (int i=0; i < size; i++) {
    crc = crc_update(crc, data[i]);
  }
  crc = ~crc;
  return crc;
}
