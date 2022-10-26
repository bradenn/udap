//
// Created by Braden Nicholson on 6/14/22.
//

#ifndef MAIN_INDICATOR_H
#define MAIN_INDICATOR_H
#include <driver/gpio.h>

typedef enum indicatorColor {
    GREEN = 0,
    RED = 1,
    BLUE = 2,
} indicatorColor;


void initIndicator();

void setIndicator(indicatorColor color);

#endif //MAIN_INDICATOR_H
