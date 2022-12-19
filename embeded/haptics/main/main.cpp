//
// Created by Braden Nicholson on 6/13/22.
//

#include <esp_timer.h>
#include <nvs_flash.h>
#include <esp_netif.h>
#include <esp_event.h>
#include "driver/gpio.h"
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "freertos/queue.h"
#include "client.cpp"
#include "server.h"
#include "haptic.h"
#include "wifi.h"

#define GPIO_INPUT_PIN_1 4


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

    while (true) {
        if (xQueueReceive(gpioQueue, &io_num, portMAX_DELAY)) {
            if (io_num != GPIO_INPUT_PIN_1) continue;
            ESP_LOGI("MOTION", "TRIGGER (%d)", io_num);
            request("http://10.0.1.2:5058/haptic-1");
        }
    }
}

extern "C" void app_main(void) {

    // Instantiate the haptic controller
    Haptic::instance();


    gpioQueue = xQueueCreate(10, sizeof(gpio_num_t));

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
