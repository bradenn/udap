// Copyright (c) 2022 Braden Nicholson

import config from "./config.js"

enum Attr {
    Controller = "controller",
    Background = "background",
    Screensaver = "screensaver",
    TouchMode = "touchmode",
    Token = "token",
    Theme = "theme",
}

class Preference {
    kind: Attr
    value: any

    // Initialize the class for a given preference
    constructor(kind: Attr) {
        // Set the kind identifier enum
        this.kind = kind
        // Try to retrieve the value form localStorage
        this.get();
    }

    // Set the value to the predefined default from config.json
    setDefault() {
        // Fetch the default value
        let defaultValue = config.defaults[this.kind]
        // Set the default value
        this.set(defaultValue)
        // Return the same value
        return defaultValue
    }

    // Fetch the value stored in localStorage, or set the default
    get() {
        // Get the local value from localStorage
        let result = window.localStorage.getItem(this.kind)
        // Check if the local value is set
        if (result == null) {
            // Set and return the default value
            return this.setDefault()
        }
        // Return the stored value
        return result
    }

    // Set the value stored in localStorage
    set(value: any) {
        // Set the value directly
        return window.localStorage.setItem(this.kind, value)
    }

}

export {
    Preference,
    Attr
};