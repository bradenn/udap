//
// Created by Braden Nicholson on 6/13/22.
//

#include "server.h"
#include "haptic.h"

extern "C" void app_main(void) {
    Haptic::instance();
    setupServer();
}
