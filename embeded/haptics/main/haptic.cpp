//
// Created by Braden Nicholson on 10/12/22.
//

#include <driver/ledc.h>
#include <thread>
#include <cmath>
#include <esp_timer.h>
#include <freertos/FreeRTOS.h>
#include <freertos/task.h>
#include "haptic.h"

Haptic &Haptic::instance() {
    static Haptic the_instance;
    return the_instance;
}

#define PX GPIO_NUM_12
#define PY GPIO_NUM_14

#define GPIO_HIGH_FREQ   GPIO_NUM_12
#define GPIO_LOW_FREQ    GPIO_NUM_14
#define GPIO_OUTPUT_PIN_SEL  ((1ULL<<GPIO_HIGH_FREQ) | (1ULL<<GPIO_LOW_FREQ))

void Haptic::allocateGpio() {

//    gpio_config_t io_conf = {};
//    //disable interrupt
//    io_conf.intr_type = GPIO_INTR_DISABLE;
//    //set as output mode
//    io_conf.mode = GPIO_MODE_OUTPUT;
//    //bit mask of the pins that you want to set,e.g.GPIO12/14
//    io_conf.pin_bit_mask = GPIO_OUTPUT_PIN_SEL;
//    //disable pull-down mode
//    io_conf.pull_down_en = GPIO_PULLDOWN_DISABLE;
//    //disable pull-up mode
//    io_conf.pull_up_en = GPIO_PULLUP_DISABLE;
//    //configure GPIO with the given settings
//    gpio_config(&io_conf);



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
    lowFrequency.gpio_num = PX;
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
    highFrequency.gpio_num = PY;
    highFrequency.intr_type = LEDC_INTR_DISABLE;
    ESP_ERROR_CHECK(ledc_channel_config(&highFrequency));
}

void Haptic::sinPulse() {

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
        sinPulse();
    }
}

int64_t micros() {
    return esp_timer_get_time();
}


void Haptic::pulseLow(int value) {
//    gpio_set_level(GPIO_LOW_FREQ, value > 0 ? 1 : 0);
    ledc_set_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_0, value);
    ledc_update_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_0);
    usleep(3125);
}

void Haptic::pulseHigh(int value) {
    ledc_set_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_1, value);
    ledc_update_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_1);
//    gpio_set_level(GPIO_HIGH_FREQ, value > 0 ? 1 : 0);
    usleep(1563);
}

void Haptic::pulse() {
    int ticks = 2;
    for (int i = 0; i < ticks; i++) {
        pulseHigh(8191);
        pulseHigh(0);
    }
}


void Haptic::lightPulse() {
    for (int i = 0; i < 2; i++) {
        pulseLow(8191);
        pulseLow(0);
    }
}

