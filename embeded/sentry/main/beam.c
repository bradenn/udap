//
// Created by Braden Nicholson on 6/14/22.
//

#include "beam.h"

#include <driver/gpio.h>
#include <driver/ledc.h>
#include <math.h>


Beam configureBeam(gpio_num_t gpio, ledc_channel_t channel, ledc_timer_t timer) {
    // Enable the PWM timer
    ledc_timer_config_t ledc_timer = {
            .speed_mode       = BEAM_SPEED_MODE,
            .timer_num        = timer,
            .duty_resolution  = LEDC_TIMER_13_BIT,
            .freq_hz          = 5000,  // 2.5khz
            .clk_cfg          = LEDC_AUTO_CLK
    };

    ESP_ERROR_CHECK(ledc_timer_config(&ledc_timer));

    // Enable the channel
    ledc_channel_config_t ledc_channel = {
            .speed_mode     = BEAM_SPEED_MODE,
            .channel        = channel,
            .timer_sel      = timer,
            .intr_type      = LEDC_INTR_DISABLE,
            .gpio_num       = gpio,
            .duty           = 0, // Set duty to 0%
            .hpoint         = 0
    };

    ESP_ERROR_CHECK(ledc_channel_config(&ledc_channel));

    Beam target;
    target.gpio = gpio;
    target.channel = channel;
    target.opticalOutput = 15; // 15mw
    target.active = 0;

    return target;
}

double map_range(double value, double low1, double high1, double low2, double high2) {
    return low2 + (high2 - low2) * (value - low1) / (high1 - low1);
}

void setBeamOpticalOutput(Beam target, int mw) {
    double mapped = map_range(mw, 0, target.opticalOutput, 0, (1 << 13) - 1);
    setBeamOutput(target, (uint32_t)round(mapped));
}

void activateBeam(Beam *target) {
    target->active = 1;
}

void deactivateBeam(Beam *target) {
    ESP_ERROR_CHECK(ledc_set_duty(BEAM_SPEED_MODE, target->channel, 0));
    ESP_ERROR_CHECK(ledc_update_duty(BEAM_SPEED_MODE, target->channel));
    target->active = 0;
}

void setBeamOutput(Beam target, uint32_t duty) {
    ESP_ERROR_CHECK(ledc_set_duty(BEAM_SPEED_MODE, target.channel, duty));
    if (target.active) {
        ESP_ERROR_CHECK(ledc_update_duty(BEAM_SPEED_MODE, target.channel));
    }
}
