//
// Created by Braden Nicholson on 6/14/22.
//

#ifndef MAIN_WIFI_H
#define MAIN_WIFI_H


typedef enum wifi_state {
    CONNECTED = 0,
    DISCONNECTED = 1,
    ERROR = 2,
} wifi_state;

wifi_state initWifi();

#endif //MAIN_WIFI_H
