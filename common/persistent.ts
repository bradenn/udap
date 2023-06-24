// Copyright (c) 2022 Braden Nicholson

import {provide, reactive, watch} from "vue";

export interface PreferencesRemote {
    accent: string
    background: string
    text: string
    blur: number
    pattern: {
        opacity: number
        scale: number
        name: string
        svg: string
    }
}

export function usePersistent(): PreferencesRemote {

    // preferenceDefaults are the default preferences for a new terminal install
    const preferenceDefaults: PreferencesRemote = {
        background: "rgba(21,21,24,1)",
        accent: "rgba(21,21,24,1)",
        text: "rgba(255,255,255,1)",
        pattern: {
            opacity: 1,
            scale: 20,
            name: "",
            svg: ""
        },
        blur: 0
    }

    // Create a reactive object to contain the preferences
    let preferences = reactive<PreferencesRemote>(<PreferencesRemote>restore())

    // Restores values from a previous save, or sets defaults
    function restore(): PreferencesRemote {
        // Get the preferences string from localStorage
        let stored = localStorage.getItem("preferences")
        // If the retrieval was successful, parse it
        if (stored) {
            // Parse the string to the Preferences object
            // Return the parsed Preferences object
            let out = JSON.parse(stored) as PreferencesRemote
            return out
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
    function save(preferences: PreferencesRemote) {
        // Convert the object to a string
        let payload = JSON.stringify(preferences)
        // Save the string to localStorage
        localStorage.setItem("preferences", payload)
    }

    // Provide the reactive Preferences object to all components
    provide('preferences', preferences)

    return preferences
}

