//
// Created by Braden Nicholson on 7/8/22.
//

#include <esp_log.h>
#include "esp_http_client.h"
#include "esp_crt_bundle.h"
#include "esp_tls.h"

#define TAG "CLIENT"

#define MAX_HTTP_OUTPUT_BUFFER 2048

esp_err_t http_event_handler(esp_http_client_event_t *evt) {
    int mbedtls_err = 0;
    esp_err_t err;

    switch (evt->event_id) {
        case HTTP_EVENT_ERROR:
            ESP_LOGD(TAG, "HTTP_EVENT_ERROR");
            break;
        case HTTP_EVENT_ON_CONNECTED:
            ESP_LOGD(TAG, "HTTP_EVENT_ON_CONNECTED");
            break;
        case HTTP_EVENT_DISCONNECTED:
            ESP_LOGI(TAG, "HTTP_EVENT_DISCONNECTED");
            err = esp_tls_get_and_clear_last_error(
                    (esp_tls_error_handle_t) evt->data, &mbedtls_err, nullptr);
            if (err != 0) {
                ESP_LOGI(TAG, "Last esp error code: 0x%x", err);
                ESP_LOGI(TAG, "Last mbedtls failure: 0x%x", mbedtls_err);
            }
            break;
        case HTTP_EVENT_REDIRECT:
            ESP_LOGD(TAG, "HTTP_EVENT_REDIRECT");
            break;
        default:
            break;
    }
    return ESP_OK;
}

void request(const char *addr) {

    char response[MAX_HTTP_OUTPUT_BUFFER] = {0};

    auto config = esp_http_client_config_t{};
    config.url = addr;
    config.event_handler = http_event_handler;
    config.user_data = response;

    esp_http_client_handle_t client = esp_http_client_init(&config);
    esp_http_client_set_method(client, HTTP_METHOD_POST);

    esp_err_t err = esp_http_client_perform(client);
    if (err == ESP_OK) {
        ESP_LOGI(TAG, "HTTP GET Status = %d, content_length = %lld",
                 esp_http_client_get_status_code(client),
                 esp_http_client_get_content_length(client));
    } else {
        ESP_LOGE(TAG, "HTTP GET request failed: %s", esp_err_to_name(err));
    }

    esp_http_client_cleanup(client);

}