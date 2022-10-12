//
// Created by Braden Nicholson on 10/12/22.
//

#ifndef HAPTIC_HAPTIC_H
#define HAPTIC_HAPTIC_H


class Haptic {
public:
    static Haptic &instance();

    Haptic(const Haptic &) = default;

    Haptic &operator=(const Haptic &) = delete;

    void pulse();
    void lightPulse();
    void sinPulse();

private:

    static void allocateGpio();

    Haptic() {
        allocateGpio();
    }

};


#endif //HAPTIC_HAPTIC_H
