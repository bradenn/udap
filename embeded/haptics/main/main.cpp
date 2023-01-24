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

int millis() {
    return (int) esp_timer_get_time() / 1000;
}

static QueueHandle_t gpioQueue = nullptr;
static int lastRequest = millis() - 750;

// gpio_isr_handler handles interrupts, it pushes the gpio number to the queue
static void IRAM_ATTR gpio_isr_handler(void *arg) {
    // Cast the input to a GPIO pin
    auto gpio_num = (uint32_t) arg;
    // Check if the time since the last request is greater than 750ms to debounce
    if (millis() - lastRequest > 750) {
        // Push the pin-number to the queue
        xQueueSendFromISR(gpioQueue, &gpio_num, nullptr);
    }
    // Set the last request time, after to loop to reset within the debounce period
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

    // Initialize a queue to process each input interrupt
    gpioQueue = xQueueCreate(10, sizeof(gpio_num_t));

    // Configure GPIO for the input pin
    gpio_config_t gpioConfig = {};
    // Select the input pin
    gpioConfig.pin_bit_mask = 1ULL << GPIO_INPUT_PIN_1;
    // Set the gpio direction to input
    gpioConfig.mode = GPIO_MODE_INPUT;
    // Disable the pull-up resistor
    gpioConfig.pull_up_en = GPIO_PULLUP_DISABLE;
    // Enable the pull-down resistor to prevent noise
    gpioConfig.pull_down_en = GPIO_PULLDOWN_ENABLE;
    // Set the interrupt to be called on the leading edge
    gpioConfig.intr_type = GPIO_INTR_POSEDGE;
    // Apply this config to the runtime
    gpio_config(&gpioConfig);

    // Install the interrupt service
    gpio_install_isr_service(0);
    // Add the pin to the interrupt service
    gpio_isr_handler_add((gpio_num_t) GPIO_INPUT_PIN_1, gpio_isr_handler, (void *) GPIO_INPUT_PIN_1);

    // Create a thread to manage the queue
    xTaskCreate(gpioTriggerQueue, "gpioQueue", 4000, nullptr, 1, nullptr);

    // Initialize the onboard flash
    ESP_ERROR_CHECK(nvs_flash_init());
    // Initialize TCP/IP stack for the runtime
    ESP_ERROR_CHECK(esp_netif_init());
    // Initialize the event loop
    ESP_ERROR_CHECK(esp_event_loop_create_default());

    // Initialize the Wi-Fi connection
    wifiInit();

    // Initialize the HTTP server instance
    Server::instance();

}
