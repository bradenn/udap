//
// Created by Braden Nicholson on 6/13/22.
//

#include <esp_timer.h>
#include <nvs_flash.h>
#include <esp_netif.h>
#include <esp_event.h>
#include <cmath>
#include "driver/gpio.h"
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "freertos/queue.h"
#include "client.cpp"
#include "server.h"
#include "haptic.h"
#include "wifi.h"

#define GPIO_INPUT_PIN_1 4

#define GPIO_OUT_1 26
#define GPIO_OUT_2 25

static ledc_channel_config_t highFrequency, lowFrequency;

int millis() {
    return (int) esp_timer_get_time() / 1000;
}

static QueueHandle_t gpioQueue = nullptr;

static void IRAM_ATTR gpio_isr_handler(void *arg) {
    auto gpio_num = (uint32_t) arg;
    static int lastRequest = millis() - 750;
    if (millis() - lastRequest > 750) {
        xQueueSendFromISR(gpioQueue, &gpio_num, nullptr);
    }
    lastRequest = millis();
}

static void gpioTriggerQueue(void *arg) {
    uint32_t io_num;
    bool swap = false;

    int val = 0;
    double cycles = 500.0;

    while (true) {
//        val = (val+1) % (int) cycles;
//        if (val >= 1000) {
//            swap = true;
//        } else if (val <= 0) {
//            swap = false;
//        }

        if (xQueueReceive(gpioQueue, &io_num, portMAX_DELAY)) {
            if (io_num != GPIO_INPUT_PIN_1) continue;
            ESP_LOGI("MOTION", "TRIGGER (%d)", io_num);
//        if (swap) {
//        if (swap) {
//            gpio_set_level(GPIO_NUM_25, 0);
//            gpio_set_level(GPIO_NUM_26, 1);
//        } else {
//            gpio_set_level(GPIO_NUM_25, 1);
//            gpio_set_level(GPIO_NUM_26, 0);
//        }

//        }else{
//            ledc_set_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_2, 0);
//            ledc_update_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_2);
//            ledc_set_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_3, 4095);
//            ledc_update_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_3);
//        }
//
//        int x = floor((4096.0/2.0) + sin(val*((M_PI*2.0)/cycles)) * (4096.0/2.0));
//        int y = floor((4096.0/2.0) + cos(val*((M_PI*2.0)/cycles)) * (4096.0/2.0));
//
//
//        ledc_set_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_2, x);
//        ledc_set_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_3, y);
//
//        ledc_update_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_2);
//        ledc_update_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_3);
//
////        swap = !swap;
//        usleep(50);
            request("http://10.0.1.2:5058/haptic-1");
        }
    }
}

extern "C" void app_main(void) {

    // Instantiate the haptic controller
    Haptic::instance();


    gpioQueue = xQueueCreate(10, sizeof(gpio_num_t));

//
//    ledc_timer_config_t ledc_timer = {
//            .speed_mode       = LEDC_LOW_SPEED_MODE,
//            .duty_resolution = LEDC_TIMER_12_BIT,
//            .timer_num        = LEDC_TIMER_2,
//            .freq_hz          = 8000,  // 160hz
//            .clk_cfg          = LEDC_AUTO_CLK
//    };
//    ESP_ERROR_CHECK(ledc_timer_config(&ledc_timer));
//
//    lowFrequency = {
//            .speed_mode     = LEDC_LOW_SPEED_MODE,
//            .channel        = LEDC_CHANNEL_2,
//            .timer_sel      = LEDC_TIMER_2,
//            .duty           = 0, // Set duty to 0%
//            .hpoint         = 0,
//    };
//    lowFrequency.gpio_num = GPIO_NUM_25;
//    lowFrequency.intr_type = LEDC_INTR_DISABLE;
//    ESP_ERROR_CHECK(ledc_channel_config(&lowFrequency));
//
//    ledc_timer = {
//            .speed_mode       = LEDC_LOW_SPEED_MODE,
//            .duty_resolution = LEDC_TIMER_12_BIT,
//            .timer_num        = LEDC_TIMER_3,
//            .freq_hz          = 8000,  // 160hz
//            .clk_cfg          = LEDC_AUTO_CLK
//    };
//
//    ESP_ERROR_CHECK(ledc_timer_config(&ledc_timer));
//
//    // Enable the channel
//    highFrequency = {
//            .speed_mode     = LEDC_LOW_SPEED_MODE,
//            .channel        = LEDC_CHANNEL_3,
//            .timer_sel      = LEDC_TIMER_3,
//            .duty           = 0, // Set duty to 0%
//            .hpoint         = 0
//    };
//    highFrequency.gpio_num = GPIO_NUM_26;
//    highFrequency.intr_type = LEDC_INTR_DISABLE;
//    ESP_ERROR_CHECK(ledc_channel_config(&highFrequency));

    gpio_config_t gpioConfig = {};
    gpioConfig.pin_bit_mask = 1ULL << GPIO_INPUT_PIN_1;
    gpioConfig.mode = GPIO_MODE_INPUT;
    gpioConfig.pull_up_en = GPIO_PULLUP_DISABLE;
    gpioConfig.pull_down_en = GPIO_PULLDOWN_ENABLE;
    gpioConfig.intr_type = GPIO_INTR_POSEDGE;
    gpio_config(&gpioConfig);


    gpio_install_isr_service(0);
    gpio_isr_handler_add((gpio_num_t) GPIO_INPUT_PIN_1, gpio_isr_handler, (void *) GPIO_INPUT_PIN_1);

    xTaskCreate(gpioTriggerQueue, "gpioQueue", 4000, nullptr, 1, nullptr);

    ESP_ERROR_CHECK(nvs_flash_init());
    ESP_ERROR_CHECK(esp_netif_init());
    ESP_ERROR_CHECK(esp_event_loop_create_default());

    wifiInit();

    Server::instance();

}
