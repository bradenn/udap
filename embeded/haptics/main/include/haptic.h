//
// Created by Braden Nicholson on 10/12/22.
//

#ifndef HAPTIC_HAPTIC_H
#define HAPTIC_HAPTIC_H
#include <driver/ledc.h>

class Haptic {
public:
    static Haptic &instance();

    Haptic(const Haptic &) = default;

    Haptic &operator=(const Haptic &) = delete;

    void sinPulseHigh();
    void sinPulseLow();

    void pulseCustom(int freq, int amp, int max);

private:

    void allocateGpio();

    ledc_channel_config_t highFrequency, lowFrequency;

    Haptic() {
        allocateGpio();
    }

    static void pulseHigh(int value);

    static void pulseLow(int value);
};


#endif //HAPTIC_HAPTIC_H
