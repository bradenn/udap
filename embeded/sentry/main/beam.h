//
// Created by Braden Nicholson on 6/14/22.
//

#ifndef MAIN_BEAM_H
#define MAIN_BEAM_H

#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include <driver/gpio.h>
#include <driver/ledc.h>

#define BEAM_SPEED_MODE LEDC_HIGH_SPEED_MODE

typedef struct Beam {
    gpio_num_t gpio;
    ledc_channel_t channel;
    int opticalOutput;
    int active;
} Beam;

Beam configureBeam(gpio_num_t gpio, ledc_channel_t channel, ledc_timer_t timer);

void setBeamOutput(Beam target, uint32_t duty);

void setBeamOpticalOutput(Beam target, int mw);

void applyBeamOutput(Beam target);

void activateBeam(Beam *target);

void deactivateBeam(Beam *target);


#endif //MAIN_BEAM_H
