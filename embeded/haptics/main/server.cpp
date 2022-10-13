//
// Created by Braden Nicholson on 8/16/22.
//

#include "server.h"
#include <esp_event.h>
#include <esp_log.h>
#include <nvs_flash.h>
#include <sys/param.h>
#include "esp_netif.h"
#include "esp_eth.h"
#include <esp_http_server.h>
#include "wifi.h"
#include "haptic.h"
#include <cJSON.h>
#include <esp_timer.h>
#include <hal/ledc_types.h>
#include <driver/ledc.h>
#include <thread>

static const char *TAG = "HTTP";

/* An HTTP GET handler */
static esp_err_t getStatusHandler(httpd_req_t *req) {
    httpd_resp_set_type(req, "application/json");
    // CORS Headers
    httpd_resp_set_hdr(req, "Access-Control-Allow-Headers", "*");
    httpd_resp_set_hdr(req, "Access-Control-Allow-Origin", "*");

    char *resp_str = formatJson();
    httpd_resp_send(req, resp_str, HTTPD_RESP_USE_STRLEN);
    free(resp_str);
    return ESP_OK;
}

static const httpd_uri_t status = {
        .uri       = "/status",
        .method    = HTTP_GET,
        .handler   = getStatusHandler,
        .user_ctx = nullptr,
        /* Let's pass response string in user
         * context to demonstrate it's usage */
};

void respondError(httpd_req_t *req, char *msg) {
    httpd_resp_send(req, msg, (int) strlen(msg));
}

static esp_err_t popOptionsHandler(httpd_req_t *req) {

    httpd_resp_set_type(req, "application/json");
    /* Set CORS header */
    httpd_resp_set_hdr(req, "Access-Control-Allow-Headers", "*");
    httpd_resp_set_hdr(req, "Access-Control-Allow-Origin", "*");

    char *resp_str = formatJson();
    httpd_resp_send(req, resp_str, HTTPD_RESP_USE_STRLEN);
    httpd_resp_set_status(req, HTTPD_200);
    free(resp_str);
    return ESP_OK;
}

/* An HTTP PUT handler */
static esp_err_t popHandler(httpd_req_t *req) {

    httpd_resp_set_type(req, "application/json");
    /* Set CORS header */
    httpd_resp_set_hdr(req, "Access-Control-Allow-Headers", "*");
    httpd_resp_set_hdr(req, "Access-Control-Allow-Origin", "*");


    char content[100];

    /* Truncate if content length larger than the buffer */
    size_t recv_size = MIN(req->content_len, sizeof(content));

    int ret = httpd_req_recv(req, content, recv_size);
    if (ret <= 0) {  /* 0 return value indicates connection closed */
        /* Check if timeout occurred */
        if (ret == HTTPD_SOCK_ERR_TIMEOUT) {
            /* In case of timeout one can choose to retry calling
             * httpd_req_recv(), but to keep it simple, here we
             * respond with an HTTP 408 (Request Timeout) error */
            httpd_resp_send_408(req);
        }
        /* In case of error, returning ESP_FAIL will
         * ensure that the underlying socket is closed */
        return ESP_FAIL;
    }

    cJSON *request = cJSON_Parse(content);

    int c = cJSON_GetObjectItem(request, "power")->valueint;
       auto h = Haptic::instance();
    if (c == 0) {
       h.pulse();

    } else if (c == 1) {
        for(int i = 0; i < 20; i++) {
            h.lightPulse();
            usleep(100*1000);
        }
    } else {

        h.sinPulse();
    }

    cJSON_Delete(request);
    char *resp_str = formatJson();
    httpd_resp_send(req, resp_str, HTTPD_RESP_USE_STRLEN);
    free(resp_str);
    return ESP_OK;
}

static const httpd_uri_t pop = {
        .uri       = "/pop",
        .method    = HTTP_POST,
        .handler   = popHandler,
        .user_ctx = nullptr,
};

static const httpd_uri_t popOptions = {
        .uri       = "/pop",
        .method    = HTTP_OPTIONS,
        .handler   = popOptionsHandler,
        .user_ctx = nullptr,
};

esp_err_t http_404_error_handler(httpd_req_t *req, httpd_err_code_t err) {
    /* For any other URI send 404 and close socket */
    httpd_resp_send_err(req, HTTPD_404_NOT_FOUND, "Some 404 error message");
    return ESP_FAIL;
}

static httpd_handle_t start_webserver(void) {
    httpd_handle_t server = NULL;
    httpd_config_t config = HTTPD_DEFAULT_CONFIG();

    esp_err_t ret = httpd_start(&server, &config);
    if (ESP_OK != ret) {
        ESP_LOGI(TAG, "Error starting server!");
        return NULL;
    }
    // Start the httpd server
    ESP_LOGI(TAG, "Starting server on port: '%d'", config.server_port);
    httpd_register_uri_handler(server, &status);
    httpd_register_uri_handler(server, &popOptions);
    httpd_register_uri_handler(server, &pop);

    return server;
}

static esp_err_t stop_webserver(httpd_handle_t server) {
    return httpd_stop(server);
}

static void disconnect_handler(void *arg, esp_event_base_t event_base,
                               int32_t event_id, void *event_data) {
    auto *server = (httpd_handle_t *) arg;
    if (*server) {
        ESP_LOGI(TAG, "Stopping webserver");
        if (stop_webserver(*server) == ESP_OK) {
            *server = NULL;
        } else {
            ESP_LOGE(TAG, "Failed to stop http server");
        }
    }
}

static void connect_handler(void *arg, esp_event_base_t event_base,
                            int32_t event_id, void *event_data) {
    auto *server = (httpd_handle_t *) arg;
    if (*server == NULL) {
        ESP_LOGI(TAG, "Starting webserver");
        *server = start_webserver();
    }
}

void setupServer() {

    ESP_ERROR_CHECK(nvs_flash_init());
    ESP_ERROR_CHECK(esp_netif_init());
    ESP_ERROR_CHECK(esp_event_loop_create_default());

    /* Register event handlers to stop the server when Wi-Fi or Ethernet is disconnected,
     * and re-start it upon connection.
     */

    wifiInit();

    static httpd_handle_t server = NULL;
    ESP_ERROR_CHECK(esp_event_handler_register(IP_EVENT, IP_EVENT_STA_GOT_IP,
                                               &connect_handler, &server));
    ESP_ERROR_CHECK(
            esp_event_handler_register(WIFI_EVENT, WIFI_EVENT_STA_DISCONNECTED,
                                       &disconnect_handler, &server));
    server = start_webserver();

    /* Start the server for the first time */

}

char *formatJson() {
    auto obj = cJSON_CreateObject();
    cJSON_AddItemToObject(obj, "status", cJSON_CreateString("OK"));

    auto val = cJSON_Print(obj);
    cJSON_Delete(obj);
    return val;

}
