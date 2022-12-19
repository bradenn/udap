//
// Created by Braden Nicholson on 8/16/22.
//

#include "server.h"
#include <esp_event.h>
#include <esp_log.h>
#include "esp_netif.h"
#include "esp_eth.h"
#include <esp_http_server.h>
#include "haptic.h"
#include <cJSON.h>
#include <esp_timer.h>
#include <hal/ledc_types.h>


static const char *TAG = "HTTP";

// Endpoint handler for `GET /status`
static esp_err_t status_get_handler(httpd_req_t *req) {
    // Send payload type header
    httpd_resp_set_type(req, "application/json");
    // Allow Cross-Origin Access
    httpd_resp_set_hdr(req, "Access-Control-Allow-Headers", "*");
    httpd_resp_set_hdr(req, "Access-Control-Allow-Origin", "*");
    // Format the status into JSON
    char *resp_str = status_format_json();
    httpd_resp_send(req, resp_str, HTTPD_RESP_USE_STRLEN);
    free(resp_str);
    return ESP_OK;
}

// Endpoint handler for `OPTIONS /status`
// Handles all preflight checks for browser related requests
static esp_err_t status_options_handler(httpd_req_t *req) {
    // Send payload type header
    httpd_resp_set_type(req, "application/json");
    // Allow Cross-Origin Access
    httpd_resp_set_hdr(req, "Access-Control-Allow-Headers", "*");
    httpd_resp_set_hdr(req, "Access-Control-Allow-Origin", "*");
    // Send the status OK message
    httpd_resp_set_status(req, HTTPD_200);
    // Return without error
    return ESP_OK;
}

static const httpd_uri_t status_get = {
        .uri       = "/status",
        .method    = HTTP_GET,
        .handler   = status_get_handler,
        .user_ctx = nullptr,
};

static const httpd_uri_t status_options = {
        .uri       = "/status",
        .method    = HTTP_OPTIONS,
        .handler   = status_options_handler,
        .user_ctx = nullptr,
};


// Socket handler is the http method handler for requests made to the /ws endpoint
static esp_err_t socket_get_handler(httpd_req_t *req) {
    // If the connection is a http GET request, initialize a new connection
    if (req->method == HTTP_GET) {
        // Log the connection
        ESP_LOGI(TAG, "Handshake done, the new connection was opened");
        // Return from the handler
        return ESP_OK;
    }
    // Instantiate a websocket frame
    httpd_ws_frame_t ws_pkt;
    // Clear the packet
    memset(&ws_pkt, 0, sizeof(httpd_ws_frame_t));
    // Set the type to websocket text
    ws_pkt.type = HTTPD_WS_TYPE_TEXT;
    // Check to make sure the frame is ready by receiving zero bytes
    esp_err_t ret = httpd_ws_recv_frame(req, &ws_pkt, 0);
    if (ret != ESP_OK) {
        ESP_LOGE(TAG, "httpd_ws_recv_frame failed with %d", ret);
        return ret;
    }
    // If the frame is not empty
    if (ws_pkt.len) {
        // Allocate memory for the incoming buffer
        auto buffer = (uint8_t *) calloc(1, ws_pkt.len + 1);;
        // Set the packet payload pointer to the buffer
        ws_pkt.payload = buffer;
        // Receive the rest of the data
        ret = httpd_ws_recv_frame(req, &ws_pkt, ws_pkt.len);
        if (ret != ESP_OK) {
            ESP_LOGE(TAG, "httpd_ws_recv_frame failed with %d", ret);
            free(buffer);
            return ret;
        }
        // If the packet type remains text, parse it.
        if (ws_pkt.type == HTTPD_WS_TYPE_TEXT) {
            // Parse the text as JSON data
            cJSON *request = cJSON_Parse((char *) (ws_pkt.payload));
            // Extract the frequency
            int f = cJSON_GetObjectItem(request, "freq")->valueint;
            // Extract the amplitude
            int a = cJSON_GetObjectItem(request, "amplitude")->valueint;
            // Extract the power
            int c = cJSON_GetObjectItem(request, "power")->valueint;
            // Get an instance of the haptic engine
            auto h = Haptic::instance();
            // Send the pulse
            h.pulseCustom(f, a, c);
            // Free the JSON memory
            cJSON_Delete(request);
        }
        // Free the buffer allocated earlier
        free(buffer);
    }
    // Return normally
    return 0;
}


static const httpd_uri_t socket_get = {
        .uri       = "/ws",
        .method    = HTTP_GET,
        .handler   = socket_get_handler,
        .is_websocket = true,
};

// Handle invalid http endpoint requests
esp_err_t http_404_error_handler(httpd_req_t *req, httpd_err_code_t err) {
    // Send a message to the client
    httpd_resp_send_err(req, HTTPD_404_NOT_FOUND, "Unauthorized request. This incident has been logged.");
    // Return an error code
    return ESP_FAIL;
}

esp_err_t Server::start_webserver() {
    // Make sure the server is not running
    if (server != nullptr) {
        return ESP_FAIL;
    }

    server = nullptr;
    config = HTTPD_DEFAULT_CONFIG();

    esp_err_t ret = httpd_start(&server, &config);
    if (ESP_OK != ret) {
        ESP_LOGI(TAG, "Error starting server!");
        return ret;
    }


    // Start the httpd server
    ESP_LOGI(TAG, "Starting server on port: '%d'", config.server_port);

    httpd_register_uri_handler(server, &status_get);
    httpd_register_uri_handler(server, &socket_get);

    httpd_register_uri_handler(server, &status_options);

    httpd_register_err_handler(server, HTTPD_404_NOT_FOUND, &http_404_error_handler);

    running = true;

    return ESP_OK;
}

esp_err_t Server::stop_webserver() {
    auto err = httpd_stop(server);
    if (err == ESP_OK) {
        server = nullptr;
        running = false;
    }
    return err;
}

static void connect_handler(void *arg, esp_event_base_t event_base, int32_t event_id, void *srv) {
    ESP_LOGI(TAG, "Wi-Fi Connected:");
    auto server = (Server *) srv;
    if (!server->is_running()) {
        ESP_LOGI(TAG, "Starting webserver");
        server->start_webserver();
    }
}

static void disconnect_handler(void *arg, esp_event_base_t event_base, int32_t event_id, void *srv) {
    ESP_LOGI(TAG, "Wi-Fi Disconnected:");
    auto server = (Server *) srv;
    if (server->is_running()) {
        ESP_LOGI(TAG, "Stopping webserver");
        if (server->stop_webserver() != ESP_OK) {
            ESP_LOGE(TAG, "Failed to stop http server");
        }
    }
}

// Returns a string containing the system status. The string must be freed after use to prevent memory leaks.
char *status_format_json() {
    // Instantiate an object
    auto obj = cJSON_CreateObject();
    // Add a status key to the object
    cJSON_AddItemToObject(obj, "status", cJSON_CreateString("OK"));
    // Write the object to a string
    auto val = cJSON_Print(obj);
    // Delete the  object
    cJSON_Delete(obj);
    // Return the string
    return val;
}

// Create an instance of the server when first called
Server &Server::instance() {
    static Server the_instance;
    return the_instance;
}

// Initialize the server object
Server::Server() {
    // Set up the event handlers
    esp_event_handler_register(IP_EVENT, IP_EVENT_STA_GOT_IP, &connect_handler, this);
    esp_event_handler_register(WIFI_EVENT, WIFI_EVENT_STA_DISCONNECTED, &disconnect_handler, this);
    // Start the webserver
    start_webserver();
}

// Returns true if the server is running
bool Server::is_running() const {
    return running;
}
