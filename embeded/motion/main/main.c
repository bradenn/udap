//
// Created by Braden Nicholson on 6/13/22.
//

#include <math.h>
#include <esp_log.h>
#include <driver/ledc.h>
#include <esp_timer.h>
#include "sdkconfig.h"
#include "driver/gpio.h"
#include "client.c"
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "freertos/queue.h"
#include "wifi.h"
#include "server.h"


static QueueHandle_t gpio_evt_queue = NULL;

static void IRAM_ATTR gpio_isr_handler(void *arg) {
    uint32_t gpio_num = (uint32_t) arg;
    xQueueSendFromISR(gpio_evt_queue, &gpio_num, NULL);
}

static void gpio_task_example(void *arg) {
    uint32_t io_num;
    int cnt = 0;
    while (1) {
        if (xQueueReceive(gpio_evt_queue, &io_num, portMAX_DELAY)) {
            cnt++;
            if(io_num == 2) {
                request("http://10.0.1.2:5058/motion-1");
            }else if(io_num == 4){
                request("http://10.0.1.2:5058/motion-2");
            }

            ESP_LOGI("MOTION", "Detected (%d)", io_num);
        }
    }
}

#define GPIO_INPUT_PIN_1 2
#define GPIO_INPUT_PIN_2 4

void app_main(void) {

    initWifi();


    gpio_config_t gpioConfig = {};
    gpioConfig.intr_type = GPIO_INTR_POSEDGE;
    gpioConfig.mode = GPIO_MODE_INPUT;
    gpioConfig.pin_bit_mask = 1ULL << GPIO_INPUT_PIN_1;
    gpioConfig.pull_up_en = 0;
    gpioConfig.pull_down_en = 0;
    gpio_config(&gpioConfig);


    gpioConfig.intr_type = GPIO_INTR_POSEDGE;
    gpioConfig.mode = GPIO_MODE_INPUT;
    gpioConfig.pin_bit_mask = 1ULL << GPIO_INPUT_PIN_2;
    gpioConfig.pull_up_en = 0;
    gpioConfig.pull_down_en = 0;
    gpio_config(&gpioConfig);

    ledc_timer_config_t ledc_timer = {
            .speed_mode       = LEDC_LOW_SPEED_MODE,
            .duty_resolution = LEDC_TIMER_12_BIT,
            .timer_num        = LEDC_TIMER_0,
            .freq_hz          = 10000,  // 160hz
            .clk_cfg          = LEDC_AUTO_CLK
    };
    ESP_ERROR_CHECK(ledc_timer_config(&ledc_timer));

    ledc_channel_config_t light = {
            .speed_mode     = LEDC_LOW_SPEED_MODE,
            .channel        = LEDC_CHANNEL_0,
            .timer_sel      = LEDC_TIMER_0,
            .duty           = 0, // Set duty to 0%
            .hpoint         = 0,
    };
    light.gpio_num = 12;
    ESP_ERROR_CHECK(ledc_channel_config(&light));



    //create a queue to handle gpio event from isr
    gpio_evt_queue = xQueueCreate(10, sizeof(uint32_t));
    //start gpio task
    gpio_install_isr_service(ESP_INTR_FLAG_LEVEL3);
    gpio_isr_handler_add(GPIO_INPUT_PIN_1, gpio_isr_handler, (void*) GPIO_INPUT_PIN_1);
    gpio_isr_handler_add(GPIO_INPUT_PIN_2, gpio_isr_handler, (void*) GPIO_INPUT_PIN_2);
    xTaskCreate(gpio_task_example, "gpio_task_example", 5000, NULL, 3, NULL);
    setupServer();
}
