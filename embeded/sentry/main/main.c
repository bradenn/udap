//
// Created by Braden Nicholson on 6/13/22.
//

#include <math.h>

#include "sdkconfig.h"

#include "beam.h"
#include "servo.h"
#include "indicator.h"
#include "server.h"

// Beam Assignments
#define BEAM_PRIMARY_PIN CONFIG_SENTRY_BEAM_PRIMARY_GPIO
#define BEAM_SECONDARY_PIN CONFIG_SENTRY_BEAM_SECONDARY_GPIO

// Servo Assignments
#define SERVO_PAN CONFIG_SENTRY_PAN_SERVO_GPIO
#define SERVO_TILT CONFIG_SENTRY_TILT_SERVO_GPIO


int map_rangeI(double value, double low1, double high1, double low2, double high2) {
    return (int)round(low2 + (high2 - low2) * (value - low1) / (high1 - low1));
}

void app_main(void) {



    Sentry sentry;

    init(&sentry);

    sentry.pan = configureServo(SERVO_PAN, MCPWM_UNIT_0, MCPWM_TIMER_0);
    sentry.tilt = configureServo(SERVO_TILT, MCPWM_UNIT_1, MCPWM_TIMER_0);

    moveTo(&sentry.pan, 0);
    moveTo(&sentry.tilt, 0);
    initIndicator();
    setIndicator(RED);
    sentry.primary = configureBeam(BEAM_PRIMARY_PIN, LEDC_CHANNEL_1, LEDC_TIMER_2);
    sentry.secondary = configureBeam(BEAM_SECONDARY_PIN, LEDC_CHANNEL_0, LEDC_TIMER_2);

    setupServer(sentry);

//    int deg = 0;
//    int dir = 0;
//
//    double x = 0;
//    double y = 0;
//    double z = 0;
//    while (1) {
//
//
//        if(deg >= 512) {
//            dir = 1;
//        }else if(deg <= 0) {
//            dir = 0;
//        }
//        deg += dir == 1?-1:1;
//
//        moveTo(&sentry.pan, 256-deg);
//        moveTo(&sentry.tilt, 256-deg);
//
////        duty = (duty + 100) % 8191;
////
////        setBeamOutput(target, duty);
////        setBeamOutput(primary, 8191 - duty);
//
////        applyBeamOutput(target);
////        applyBeamOutput(primary);
//
////        moveTo(&sentry.pan, deg);
//
//
//        vTaskDelay(pdMS_TO_TICKS(50));
//    }
}
