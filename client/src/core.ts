// Copyright (c) 2022 Braden Nicholson

import {inject} from "vue";
import type {Haptics} from "@/haptics";
import type {Router} from "vue-router";
import {useRouter} from "vue-router";
import type {Context} from "@/context";
import type {Notify} from "@/notifications";
import type {Screensaver} from "@/screensaver";
import type {Remote} from "@/remote";


export interface Core {
    router: () => Router,
    remote: () => Remote,
    context: () => Context,
    screensaver: () => Screensaver,
    haptics: () => Haptics,
    notify: () => Notify,
}

export default {
    router: () => useRouter(),
    remote: () => inject("remote") as Remote,
    context: () => inject("context") as Context,
    screensaver: () => inject("screens") as Screensaver,
    haptics: () => inject("haptics") as Haptics,
    notify: () => inject("notify") as Notify,
} as Core

