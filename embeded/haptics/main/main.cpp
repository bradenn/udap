//
// Created by Braden Nicholson on 6/13/22.
//

#include <esp_timer.h>
#include "driver/gpio.h"
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "freertos/queue.h"
#include "client.cpp"
#include "server.h"
#include "haptic.h"

#define GPIO_INPUT_PIN_1 2


int millis() {
    return (int) esp_timer_get_time() / 1000;
}

void setupInput() {
    gpio_config_t gpioConfig = {};
    gpioConfig.intr_type = GPIO_INTR_POSEDGE;
    gpioConfig.mode = GPIO_MODE_INPUT;
    gpioConfig.pin_bit_mask = 1ULL << GPIO_INPUT_PIN_1;
    gpioConfig.pull_up_en = GPIO_PULLUP_DISABLE;
    gpioConfig.pull_down_en = GPIO_PULLDOWN_ENABLE;
    gpio_config(&gpioConfig);
}

static QueueHandle_t gpioQueue = nullptr;

static void IRAM_ATTR gpio_isr_handler(void *arg) {
    auto gpio_num = (uint32_t) arg;
    xQueueSendFromISR(gpioQueue, &gpio_num, nullptr);
}

static void gpioTriggerQueue(void *arg) {
    uint32_t io_num;
    int lastRequest = millis() - 500;
    while (true) {
        if (xQueueReceive(gpioQueue, &io_num, portMAX_DELAY/1000)) {
            if(io_num != GPIO_INPUT_PIN_1) continue;
            if (millis() - lastRequest >= 500) {
                request("http://10.0.1.2:5058/haptic-1");
                ESP_LOGI("MOTION", "TRIGGER (%d)", io_num);
            }

            lastRequest = millis();
        }
    }
}

extern "C" void app_main(void) {
    // Initialize input GPIO pins
    setupInput();
    // Instantiate the haptic controller
    Haptic::instance();


    gpioQueue = xQueueCreate(10, sizeof(uint32_t));
    gpio_install_isr_service(ESP_INTR_FLAG_LEVEL3);
    gpio_isr_handler_add(static_cast<gpio_num_t>(GPIO_INPUT_PIN_1),
                         gpio_isr_handler, (void *) GPIO_INPUT_PIN_1);

    xTaskCreate(gpioTriggerQueue, "gpioTriggerQueue", 5000, nullptr, 3,
                nullptr);

    // Load the web server
    setupServer();

}
