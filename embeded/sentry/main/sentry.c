//
// Created by Braden Nicholson on 6/25/22.
//

#include "sentry.h"
#include <cJSON.h>
#include <esp_mac.h>


//double readTemp(Sentry sentry) {
//    ESP_ERROR_CHECK(temperature_sensor_enable(sentry.temp));
//    float tsens_out;
//    ESP_ERROR_CHECK(temperature_sensor_get_celsius(sentry.temp, &tsens_out));
//    ESP_ERROR_CHECK(temperature_sensor_disable(sentry.temp));
//    return (double) tsens_out;
//}

char *formatJSON(Sentry sentry) {

    // Parent Status Object
    cJSON *status = cJSON_CreateObject();

    // System Info
    cJSON *system = cJSON_CreateObject();

    // Get the mac address

    cJSON_AddItemToObject(system, "mac", cJSON_CreateString(sentry.macAddress));
//    cJSON_AddItemToObject(system, "temp", cJSON_CreateNumber(readTemp(sentry)));

    cJSON_AddItemToObject(status, "system", system);
    // Servo Object
    cJSON *servos = cJSON_CreateObject();

    cJSON *panServo = cJSON_CreateNumber(sentry.pan.position);
    cJSON_AddItemToObject(servos, "pan", panServo);

    cJSON *tiltServo = cJSON_CreateNumber(sentry.tilt.position);
    cJSON_AddItemToObject(servos, "tilt", tiltServo);

    // Beam Object
    cJSON *beams = cJSON_CreateObject();

    cJSON *beamPrimary = cJSON_CreateBool(sentry.primary.active == 1);
    cJSON_AddItemToObject(beams, "primary", beamPrimary);

    cJSON *beamSecondary = cJSON_CreateBool(sentry.secondary.active == 1);
    cJSON_AddItemToObject(beams, "secondary", beamSecondary);

    cJSON_AddItemToObject(status, "servos", servos);
    cJSON_AddItemToObject(status, "beams", beams);

    char *formatted = cJSON_Print(status);

    cJSON_Delete(status);

    return formatted;
}

void init(Sentry *sentry) {
//    sentry->temp = NULL;
//    temperature_sensor_config_t temp_sensor = {
//            .range_min = 20,
//            .range_max = 50,
//    };
//    ESP_ERROR_CHECK(temperature_sensor_install(&temp_sensor, &sentry->temp));

    uint8_t mac[6];
    esp_read_mac(mac, ESP_MAC_WIFI_STA);
    sprintf(sentry->macAddress, "%x:%x:%x:%x:%x:%x", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5]);
}
