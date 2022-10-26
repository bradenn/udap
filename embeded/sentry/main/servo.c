#include <driver/gpio.h>
#include "servo.h"
//
// Created by Braden Nicholson on 6/14/22.
//

static uint32_t angleToDuty(int angle)
{
    return (angle + MAX_POSITION) * (MAX_PULSE_US - MIN_PULSE_US) / (2 * MAX_POSITION) + MIN_PULSE_US;
}


Servo configureServo(gpio_num_t gpio, mcpwm_unit_t unit, mcpwm_timer_t timer) {
    ESP_ERROR_CHECK(mcpwm_gpio_init(unit, MCPWM0A, gpio));
    mcpwm_config_t pwm_config = {
            .frequency = 50, // frequency = 50Hz, i.e. for every servo motor time period should be 20ms
            .cmpr_a = 0,     // duty cycle of PWMxA = 0
            .counter_mode = MCPWM_UP_COUNTER,
            .duty_mode = MCPWM_DUTY_MODE_0,
    };

    ESP_ERROR_CHECK(mcpwm_init(unit, timer, &pwm_config));
    Servo servo;
    servo.timer = timer;
    servo.unit = unit;
    servo.gpio = gpio;
    servo.position = 0;
    return servo;
}

void moveTo(Servo *servo, int position) {
    if(position < -MAX_POSITION || position > MAX_POSITION){
        return;
    }
    ESP_ERROR_CHECK(mcpwm_set_duty_in_us(servo->unit, servo->timer, MCPWM_OPR_A, angleToDuty(position)));
    servo->position = position;
//    vTaskDelay(pdMS_TO_TICKS(10));
}