//
// Created by Braden Nicholson on 6/24/22.
//

#ifndef SENTRY_SERVER_H
#define SENTRY_SERVER_H

#include <esp_http_server.h>
#include <esp_event_base.h>

char *status_format_json();

class Server {
public:
    static Server &instance();

    Server(const Server &) = default;

    Server &operator=(const Server &) = delete;

    esp_err_t start_webserver();

    esp_err_t stop_webserver();

    bool is_running() const;

private:
    Server();

    bool running = false;

    httpd_handle_t server{};
    httpd_config_t config{};


};


#endif //SENTRY_SERVER_H
