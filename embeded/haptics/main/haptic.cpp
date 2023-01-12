//
// Created by Braden Nicholson on 10/12/22.
//

#include <driver/ledc.h>
#include <thread>
#include <cmath>
#include <esp_timer.h>
#include "haptic.h"

Haptic &Haptic::instance() {
    static Haptic the_instance;
    return the_instance;
}

#define GPIO_HIGH_FREQ   GPIO_NUM_12
#define GPIO_LOW_FREQ    GPIO_NUM_14

void Haptic::allocateGpio() {

    ledc_timer_config_t ledc_timer = {
            .speed_mode       = LEDC_LOW_SPEED_MODE,
            .duty_resolution = LEDC_TIMER_12_BIT,
            .timer_num        = LEDC_TIMER_0,
            .freq_hz          = 10000,  // 160hz
            .clk_cfg          = LEDC_AUTO_CLK
    };
    ESP_ERROR_CHECK(ledc_timer_config(&ledc_timer));

    lowFrequency = {
            .speed_mode     = LEDC_LOW_SPEED_MODE,
            .channel        = LEDC_CHANNEL_0,
            .timer_sel      = LEDC_TIMER_0,
            .duty           = 0, // Set duty to 0%
            .hpoint         = 0,
    };
    lowFrequency.gpio_num = GPIO_HIGH_FREQ;
    lowFrequency.intr_type = LEDC_INTR_DISABLE;
    ESP_ERROR_CHECK(ledc_channel_config(&lowFrequency));

    ledc_timer = {
            .speed_mode       = LEDC_LOW_SPEED_MODE,
            .duty_resolution = LEDC_TIMER_12_BIT,
            .timer_num        = LEDC_TIMER_1,
            .freq_hz          = 10000,  // 160hz
            .clk_cfg          = LEDC_AUTO_CLK
    };

    ESP_ERROR_CHECK(ledc_timer_config(&ledc_timer));

    // Enable the channel
    highFrequency = {
            .speed_mode     = LEDC_LOW_SPEED_MODE,
            .channel        = LEDC_CHANNEL_1,
            .timer_sel      = LEDC_TIMER_1,
            .duty           = 0, // Set duty to 0%
            .hpoint         = 0
    };
    highFrequency.gpio_num = GPIO_LOW_FREQ;
    highFrequency.intr_type = LEDC_INTR_DISABLE;
    ESP_ERROR_CHECK(ledc_channel_config(&highFrequency));
}

void Haptic::sinPulseLow() {

    int steps = 20;
    double dr = (M_PI * 2) / steps;
    int dt = 1563 / steps;
    for (int j = 0; j < steps; j++) {
        double r = pow(cos((j * dr)), 2) * 4095.0;
        ledc_set_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_1, floor(r));
        ledc_update_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_1);
        usleep(dt);
    }
    ledc_set_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_1, 0);
    ledc_update_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_1);

}

void Haptic::sinPulseHigh() {

    int steps = 20;
    double dr = (M_PI * 2) / steps;
    int dt = 3125 / steps;
    for (int j = 0; j < steps; j++) {
        double r = pow(cos((j * dr)), 2) * 4095.0;
        ledc_set_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_0, floor(r));
        ledc_update_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_0);
        usleep(dt);
    }
    ledc_set_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_0, 0);
    ledc_update_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_0);

}

void Haptic::pulseCustom(int freq, int amp, int max) {
    if (freq == 0) {
        for (int j = 0; j < max; j++) {
            pulseLow(amp);
            pulseLow(0);
        }
    } else if (freq == 1) {
        for (int j = 0; j < max; j++) {
            pulseHigh(amp);
            pulseHigh(0);
        }
    } else if (freq == 2) {
        sinPulseLow();
    }else if (freq == 3) {
        sinPulseHigh();
    }
}

void Haptic::pulseLow(int value) {
    ledc_set_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_0, value);
    ledc_update_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_0);
    usleep(3125);
}

void Haptic::pulseHigh(int value) {
    ledc_set_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_1, value);
    ledc_update_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_1);
    usleep(1563);
}


