// Copyright (c) 2022 Braden Nicholson

import type {Preferences} from "@/types";
import {provide, reactive, watch} from "vue";

export function usePersistent() {

    // preferenceDefaults are the default preferences for a new terminal install
    const preferenceDefaults: Preferences = {
        ui: {
            screensaver: {
                enabled: true,
                countdown: 60 * 10,
                selection: "bubbles"
            },
            background: {
                image: "swirl",
                blur: true
            },
            blur: 6,
            mode: "cursor",
            theme: "dark",
            brightness: 100,
            grid: false,
            watermark: false,
            night: false,
            outlines: true,
        }
    }

    // Create a reactive object to contain the preferences
    let preferences = reactive<Preferences>(restore())

    // Restores values from a previous save, or sets defaults
    function restore() {
        // Get the preferences string from localStorage
        let stored = localStorage.getItem("preferences")
        // If the retrieval was successful, parse it
        if (stored) {
            // Parse the string to the Preferences object
            let parsed: Preferences = JSON.parse(stored)
            // Return the parsed Preferences object
            return parsed
        } else {
            // If the retrieval failed, save the default values to localStorage
            save(preferenceDefaults)
            // Return the default parameters
            return preferenceDefaults
        }
    }

    // Watch all changes made to the Preferences reactive object
    watch(preferences, () => {
        // Save any changes to localStorage
        save(preferences)
    })

    // Save the Preferences object to localStorage
    function save(preferences: Preferences) {
        // Convert the object to a string
        let payload = JSON.stringify(preferences)
        // Save the string to localStorage
        localStorage.setItem("preferences", payload)
    }

    // Provide the reactive Preferences object to all components
    provide('preferences', preferences)

    return preferences
}

