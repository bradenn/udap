//
// Created by Braden Nicholson on 6/25/22.
//

#ifndef SENTRY_SENTRY_H
#define SENTRY_SENTRY_H

#include <driver/temperature_sensor.h>
#include "servo.h"
#include "beam.h"

typedef struct Sentry {
    Servo pan;
    Servo tilt;

    Beam primary;
    Beam secondary;

    char macAddress[48];

} Sentry;


void init(Sentry*);
char *formatJSON(Sentry);


#endif //SENTRY_SENTRY_H
