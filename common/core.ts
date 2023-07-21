// Copyright (c) 2022 Braden Nicholson

import {inject} from "vue";
import type {Router} from "vue-router";
import {useRouter} from "vue-router";
import type {Remote} from "./remote";
import type {PreferencesRemote} from "./persistent";
import {RemoteTimings} from "./timings";

// This interface defines the set of common objects accessible
export interface Core {
    preferences: () => PreferencesRemote,
    router: () => Router,
    remote: () => Remote,
    timings: () => RemoteTimings
}

export default {
    preferences: () => inject("preferences") as PreferencesRemote,
    router: () => useRouter(),
    remote: () => inject("remote") as Remote,
    timings: () => inject("timings") as RemoteTimings
} as Core

