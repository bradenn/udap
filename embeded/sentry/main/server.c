//
// Created by Braden Nicholson on 6/24/22.
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
#include "indicator.h"
#include "sentry.h"
#include <cJSON.h>
#include <esp_timer.h>

static const char *TAG = "HTTP";

#define HTTPD_401      "401 UNAUTHORIZED"           /*!< HTTP Response 401 */

#define USERNAME CONFIG_SENTRY_AUTH_USER
#define PASSWORD CONFIG_SENTRY_AUTH_PASS

static Sentry sentry;

/* An HTTP GET handler */
static esp_err_t getStatusHandler(httpd_req_t *req) {

    httpd_resp_set_type(req, "application/json");

    char *resp_str = formatJSON(sentry);

    httpd_resp_send(req, resp_str, HTTPD_RESP_USE_STRLEN);

    free(resp_str);
    return ESP_OK;
}

static const httpd_uri_t status = {
        .uri       = "/status",
        .method    = HTTP_GET,
        .handler   = getStatusHandler,
        /* Let's pass response string in user
         * context to demonstrate it's usage */
        .user_ctx  = "Hello World!"
};

void respondError(httpd_req_t *req, char *msg) {
    httpd_resp_send(req, msg, (int) strlen(msg));
}


/* An HTTP PUT handler */
static esp_err_t positionHandler(httpd_req_t *req) {

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

    if (cJSON_HasObjectItem(request, "token")) {
        char *token = cJSON_GetObjectItem(request, "token")->valuestring;
        if (strncmp(token, PASSWORD, strlen(PASSWORD)) != 0) {
            respondError(req, "invalid token, this even has been logged.");
            cJSON_Delete(request);
            return ESP_FAIL;
        }
    } else {
        respondError(req, "a security token must be provided");
        cJSON_Delete(request);
        return ESP_FAIL;
    }

    int panPos;
    int tiltPos;

    if (cJSON_HasObjectItem(request, "pan")) {
        panPos = cJSON_GetObjectItem(request, "pan")->valueint;
        moveTo(&sentry.pan, panPos);
    }

    if (cJSON_HasObjectItem(request, "tilt")) {
        tiltPos = cJSON_GetObjectItem(request, "tilt")->valueint;
        moveTo(&sentry.tilt, tiltPos);
    }


    cJSON_Delete(request);

    char *resp_str = formatJSON(sentry);
    httpd_resp_set_type(req, "application/json");
    httpd_resp_send(req, resp_str, HTTPD_RESP_USE_STRLEN);
    free(resp_str);

    return ESP_OK;
}

static const httpd_uri_t position = {
        .uri       = "/position",
        .method    = HTTP_POST,
        .handler   = positionHandler,
};


/* An HTTP PUT handler */
static esp_err_t beamHandler(httpd_req_t *req) {

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

    if (!cJSON_HasObjectItem(request, "target") || !cJSON_HasObjectItem(request, "active") ||
        !cJSON_HasObjectItem(request, "power")) {
        httpd_resp_send_404(req);
        cJSON_Delete(request);
        return ESP_FAIL;
    }

    char *target = cJSON_GetObjectItem(request, "target")->valuestring;

    int pos = cJSON_GetObjectItem(request, "active")->valueint;
    int power = cJSON_GetObjectItem(request, "power")->valueint;

    if (strncmp(target, "primary", 7) == 0) {
        if (pos == 1) {
            activateBeam(&sentry.primary);
            setBeamOpticalOutput(sentry.primary, power);
        } else {
            deactivateBeam(&sentry.primary);
        }
    } else if (strncmp(target, "secondary", 9) == 0) {
        if (pos == 1) {
            activateBeam(&sentry.secondary);
            setBeamOpticalOutput(sentry.secondary, power);
        } else {
            deactivateBeam(&sentry.secondary);
        }
    } else {
        respondError(req, "target beam does not exist");
        cJSON_Delete(request);
        return ESP_FAIL;
    }

    char *resp_str = formatJSON(sentry);

    httpd_resp_send(req, resp_str, HTTPD_RESP_USE_STRLEN);

    free(resp_str);

    cJSON_Delete(request);

    return ESP_OK;
}

static const httpd_uri_t beam = {
        .uri       = "/beam",
        .method    = HTTP_POST,
        .handler   = beamHandler,
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
    httpd_register_uri_handler(server, &position);
    httpd_register_uri_handler(server, &beam);

    return server;
}

static esp_err_t stop_webserver(httpd_handle_t server) {
    return httpd_stop(server);
}

static void disconnect_handler(void *arg, esp_event_base_t event_base,
                               int32_t event_id, void *event_data) {
    httpd_handle_t *server = (httpd_handle_t *) arg;
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
    httpd_handle_t *server = (httpd_handle_t *) arg;
    if (*server == NULL) {
        ESP_LOGI(TAG, "Starting webserver");
        *server = start_webserver();
    }
}

void setupServer(Sentry sen) {
    sentry = sen;
    moveTo(&sentry.pan, 0);
    moveTo(&sentry.tilt, 0);
    setIndicator(RED);

    ESP_ERROR_CHECK(nvs_flash_init());
    ESP_ERROR_CHECK(esp_netif_init());
    ESP_ERROR_CHECK(esp_event_loop_create_default());

    /* Register event handlers to stop the server when Wi-Fi or Ethernet is disconnected,
     * and re-start it upon connection.
     */

    int result = initWifi();
    setIndicator(result);

    static httpd_handle_t server = NULL;
    ESP_ERROR_CHECK(esp_event_handler_register(IP_EVENT, IP_EVENT_STA_GOT_IP, &connect_handler, &server));
    ESP_ERROR_CHECK(esp_event_handler_register(WIFI_EVENT, WIFI_EVENT_STA_DISCONNECTED, &disconnect_handler, &server));
    server = start_webserver();

    /* Start the server for the first time */

}