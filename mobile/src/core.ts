// Copyright (c) 2022 Braden Nicholson

import {inject} from "vue";
import type {Router} from "vue-router";
import {useRouter} from "vue-router";
import type {Remote} from "udap-ui/remote";
import type {PreferencesRemote} from "udap-ui/persistent";


export interface Core {
    preferences: () => PreferencesRemote,
    router: () => Router,
    remote: () => Remote,
}

export default {
    preferences: () => inject("preferences"),
    router: () => useRouter(),
    remote: () => inject("remote") as Remote,
} as Core

