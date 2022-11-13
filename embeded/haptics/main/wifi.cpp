//
// Created by Braden Nicholson on 6/14/22.
//

#include "include/wifi.h"
#include "freertos/FreeRTOS.h"
#include "freertos/event_groups.h"
#include "esp_system.h"
#include "esp_wifi.h"
#include "esp_event.h"
#include "esp_log.h"

// Basic Wi-Fi Connection Status Bits
#define WIFI_CONNECTED_BIT BIT0
#define WIFI_FAIL_BIT      BIT1

// Compiler Config Variables
#define WIFI_SSID CONFIG_ESP_WIFI_SSID
#define WIFI_PASSWORD CONFIG_ESP_WIFI_PASSWORD
#define WIFI_MAX_RETRY CONFIG_ESP_MAXIMUM_RETRY

// Log Output Title
static const char *TAG = "WIFI";

// Reconnect retry attempts
static int s_retry_num = 0;

// Wifi Handler
static EventGroupHandle_t s_wifi_event_group;

static void eventHandler(void *arg, esp_event_base_t event_base,
                         int32_t event_id, void *event_data) {

    if (event_base == WIFI_EVENT && event_id == WIFI_EVENT_STA_START) {

        // Attempt first contact
        esp_wifi_connect();

    } else if (event_base == WIFI_EVENT
               && event_id == WIFI_EVENT_STA_DISCONNECTED) {

        if (s_retry_num < WIFI_MAX_RETRY) {

            // Try to reconnect
            esp_wifi_connect();

            // Increase the retry counter
            s_retry_num++;

            // Announce the attempt to reconnect
            ESP_LOGI(TAG, "Attempting to reconnect...");

        } else {
            // Send the failed bit
            xEventGroupSetBits(s_wifi_event_group, WIFI_FAIL_BIT);
        }

        ESP_LOGI(TAG, "Connection Failed");

    } else if (event_base == IP_EVENT && event_id == IP_EVENT_STA_GOT_IP) {

        auto *event = (ip_event_got_ip_t *) event_data;

        // Print the DHCP IP to the console
        ESP_LOGI(TAG, "IPv4: " IPSTR, IP2STR(&event->ip_info.ip));

        // Reset the retry counter on a good connection
        s_retry_num = 0;

        // Send the success bit
        xEventGroupSetBits(s_wifi_event_group, WIFI_CONNECTED_BIT);

    }
}


int wifiInit() {

    s_wifi_event_group = xEventGroupCreate();

    esp_netif_create_default_wifi_sta();

    wifi_init_config_t cfg = WIFI_INIT_CONFIG_DEFAULT();

    ESP_ERROR_CHECK(esp_wifi_init(&cfg));
    ESP_ERROR_CHECK(esp_wifi_set_ps(WIFI_PS_NONE));

    esp_event_handler_instance_t instance_any_id, instance_got_ip;

    ESP_ERROR_CHECK(esp_event_handler_instance_register(
            WIFI_EVENT, ESP_EVENT_ANY_ID, &eventHandler,
            nullptr, &instance_any_id));

    ESP_ERROR_CHECK(esp_event_handler_instance_register(
            IP_EVENT, IP_EVENT_STA_GOT_IP, &eventHandler,
            nullptr, &instance_got_ip));

    wifi_config_t wifi_config = {
            .sta = {
                    .ssid = WIFI_SSID,
                    .password = WIFI_PASSWORD,
                    .threshold = {
                            .authmode =WIFI_AUTH_WPA2_PSK,
                    },
            },
    };


    ESP_ERROR_CHECK(esp_wifi_set_mode(WIFI_MODE_STA));
    ESP_ERROR_CHECK(esp_wifi_set_config(WIFI_IF_STA, &wifi_config));
    ESP_ERROR_CHECK(esp_wifi_start());

    ESP_LOGI(TAG, "Operating in station mode.");

    EventBits_t bits = xEventGroupWaitBits(s_wifi_event_group,
                                           WIFI_CONNECTED_BIT | WIFI_FAIL_BIT,
                                           pdFALSE, pdFALSE, portMAX_DELAY);

    int ret_value;
    if (bits & WIFI_CONNECTED_BIT) {

        ret_value = CONNECTED;
        ESP_LOGI(TAG, "Connected to AP With SSID: %s", WIFI_SSID);

    } else if (bits & WIFI_FAIL_BIT) {

        ret_value = DISCONNECTED;
        ESP_LOGI(TAG, "Failed to connect to AP With SSID: %s", WIFI_SSID);

    } else {

        ret_value = ERROR;
        ESP_LOGE(TAG, "An unexpected event occurred");

    }

    vEventGroupDelete(s_wifi_event_group);

    return ret_value;
}
