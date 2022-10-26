//
// Created by Braden Nicholson on 6/14/22.
//


#include <driver/gpio.h>
#include <driver/ledc.h>
#include "indicator.h"

#define INDICATOR_RED CONFIG_GPIO_INDICATOR_RED
#define INDICATOR_GREEN CONFIG_GPIO_INDICATOR_GREEN
#define INDICATOR_BLUE CONFIG_GPIO_INDICATOR_BLUE

void assignGPIOPin(int gpio, ledc_timer_t timer, ledc_channel_t channel) {
    // Enable the PWM timer
    ledc_timer_config_t ledc_timer = {
            .speed_mode       = LEDC_HIGH_SPEED_MODE,
            .timer_num        = timer,
            .duty_resolution  = LEDC_TIMER_8_BIT,
            .freq_hz          = 5000,  // 2.5khz
            .clk_cfg          = LEDC_AUTO_CLK
    };

    ESP_ERROR_CHECK(ledc_timer_config(&ledc_timer));

    // Enable the channel
    ledc_channel_config_t ledc_channel = {
            .speed_mode     = LEDC_HIGH_SPEED_MODE,
            .channel        = channel,
            .timer_sel      = timer,
            .intr_type      = LEDC_INTR_DISABLE,
            .gpio_num       = gpio,
            .duty           = 0, // Set duty to 0%
            .hpoint         = 0
    };

    ESP_ERROR_CHECK(ledc_channel_config(&ledc_channel));
}


void setIndicator(indicatorColor color) {
    switch (color) {
        case RED:
            ESP_ERROR_CHECK(ledc_set_duty(LEDC_HIGH_SPEED_MODE, LEDC_CHANNEL_5, 64));
            ESP_ERROR_CHECK(ledc_update_duty(LEDC_HIGH_SPEED_MODE, LEDC_CHANNEL_5));

            ESP_ERROR_CHECK(ledc_set_duty(LEDC_HIGH_SPEED_MODE, LEDC_CHANNEL_6, 0));
            ESP_ERROR_CHECK(ledc_update_duty(LEDC_HIGH_SPEED_MODE, LEDC_CHANNEL_6));

            ESP_ERROR_CHECK(ledc_set_duty(LEDC_HIGH_SPEED_MODE, LEDC_CHANNEL_7, 0));
            ESP_ERROR_CHECK(ledc_update_duty(LEDC_HIGH_SPEED_MODE, LEDC_CHANNEL_7));
            break;
        case GREEN:
            ESP_ERROR_CHECK(ledc_set_duty(LEDC_HIGH_SPEED_MODE, LEDC_CHANNEL_5, 0));
            ESP_ERROR_CHECK(ledc_update_duty(LEDC_HIGH_SPEED_MODE, LEDC_CHANNEL_5));

            ESP_ERROR_CHECK(ledc_set_duty(LEDC_HIGH_SPEED_MODE, LEDC_CHANNEL_6, 64));
            ESP_ERROR_CHECK(ledc_update_duty(LEDC_HIGH_SPEED_MODE, LEDC_CHANNEL_6));

            ESP_ERROR_CHECK(ledc_set_duty(LEDC_HIGH_SPEED_MODE, LEDC_CHANNEL_7, 0));
            ESP_ERROR_CHECK(ledc_update_duty(LEDC_HIGH_SPEED_MODE, LEDC_CHANNEL_7));
            break;
        case BLUE:
            ESP_ERROR_CHECK(ledc_set_duty(LEDC_HIGH_SPEED_MODE, LEDC_CHANNEL_5, 0));
            ESP_ERROR_CHECK(ledc_update_duty(LEDC_HIGH_SPEED_MODE, LEDC_CHANNEL_5));

            ESP_ERROR_CHECK(ledc_set_duty(LEDC_HIGH_SPEED_MODE, LEDC_CHANNEL_6, 0));
            ESP_ERROR_CHECK(ledc_update_duty(LEDC_HIGH_SPEED_MODE, LEDC_CHANNEL_6));

            ESP_ERROR_CHECK(ledc_set_duty(LEDC_HIGH_SPEED_MODE, LEDC_CHANNEL_7, 64));
            ESP_ERROR_CHECK(ledc_update_duty(LEDC_HIGH_SPEED_MODE, LEDC_CHANNEL_7));
            break;
        default:
            break;
    }
}

void initIndicator() {

    assignGPIOPin(INDICATOR_RED, LEDC_TIMER_2, LEDC_CHANNEL_5);
    assignGPIOPin(INDICATOR_GREEN, LEDC_TIMER_2, LEDC_CHANNEL_6);
    assignGPIOPin(INDICATOR_BLUE, LEDC_TIMER_2, LEDC_CHANNEL_7);

}