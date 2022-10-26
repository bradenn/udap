//
// Created by Braden Nicholson on 6/14/22.
//
#ifndef MAIN_SERVO_H
#define MAIN_SERVO_H
#include <driver/gpio.h>
#include "driver/mcpwm.h"

#define MIN_PULSE_US  (500)
#define MAX_PULSE_US  (2400)

#define MAX_POSITION (90)


typedef struct Servo {
    gpio_num_t gpio;
    mcpwm_timer_t timer;
    mcpwm_unit_t unit;
    int position;
} Servo;

Servo configureServo(gpio_num_t gpio, mcpwm_unit_t channel, mcpwm_timer_t timer);

void moveTo(Servo *servo, int position);

#endif //MAIN_SERVO_H
