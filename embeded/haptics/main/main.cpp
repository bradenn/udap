//
// Created by Braden Nicholson on 6/13/22.
//


#include <esp_timer.h>
#include <driver/ledc.h>
#include "server.h"
#include "haptic.h"

extern "C" void app_main(void) {

    auto haptics = Haptic::instance();

    haptics.lightPulse();

    setupServer();
}
