//
// Created by Braden Nicholson on 10/12/22.
//

#include <driver/ledc.h>
#include <thread>
#include <cmath>
#include "haptic.h"

Haptic &Haptic::instance() {
    static Haptic the_instance;
    return the_instance;
}

#define PX GPIO_NUM_12
#define PY GPIO_NUM_14

void Haptic::allocateGpio() {
    ledc_timer_config_t ledc_timer = {
            .speed_mode       = LEDC_LOW_SPEED_MODE,
            .duty_resolution = LEDC_TIMER_13_BIT,
            .timer_num        = LEDC_TIMER_0,
            .freq_hz          = 160,  // 160hz
            .clk_cfg          = LEDC_AUTO_CLK
    };
    ESP_ERROR_CHECK(ledc_timer_config(&ledc_timer));

    ledc_channel_config_t ledc_channel = {
            .speed_mode     = LEDC_LOW_SPEED_MODE,
            .channel        = LEDC_CHANNEL_0,
            .timer_sel      = LEDC_TIMER_0,
            .duty           = 0, // Set duty to 0%
            .hpoint         = 0,
    };
    ledc_channel.gpio_num = PX;
    ledc_channel.intr_type = LEDC_INTR_DISABLE;
    ESP_ERROR_CHECK(ledc_channel_config(&ledc_channel));

    ledc_timer = {
            .speed_mode       = LEDC_LOW_SPEED_MODE,
            .duty_resolution = LEDC_TIMER_13_BIT,
            .timer_num        = LEDC_TIMER_0,
            .freq_hz          = 320,  // 160hz
            .clk_cfg          = LEDC_AUTO_CLK
    };

    ESP_ERROR_CHECK(ledc_timer_config(&ledc_timer));

    // Enable the channel
    ledc_channel = {
            .speed_mode     = LEDC_LOW_SPEED_MODE,
            .channel        = LEDC_CHANNEL_1,
            .timer_sel      = LEDC_TIMER_0,
            .duty           = 0, // Set duty to 0%
            .hpoint         = 0
    };
    ledc_channel.gpio_num = PY;
    ledc_channel.intr_type = LEDC_INTR_DISABLE;
    ESP_ERROR_CHECK(ledc_channel_config(&ledc_channel));
}

void Haptic::sinPulse() {
    int duration = 1875;

    for (int j = 0; j < 2; j++) {
        for (int i = 0; i < duration / 100; i++) {
            double s = pow(sin((double) i / 2), 2);
            ledc_set_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_1,
                          ceil(8191.0 * s));
            ledc_update_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_1);
            usleep(100);
        }
        ledc_set_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_1, 0);
        ledc_update_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_1);
        usleep(3125);
    }

}

void Haptic::pulse() {
    for (int i = 0; i < 2; i++) {
        ledc_set_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_1, 8191);
        ledc_update_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_1);
        usleep(3125);
        ledc_set_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_1, 0);
        ledc_update_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_1);
        usleep(3125);
    }
}

void Haptic::lightPulse() {
    for (int i = 0; i < 2; i++) {
        ledc_set_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_0, 8191);
        ledc_update_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_0);
        usleep(1563);
        ledc_set_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_0, 0);
        ledc_update_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_0);
        usleep(1563);
    }
}

