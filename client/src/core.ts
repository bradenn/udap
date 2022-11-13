// Copyright (c) 2022 Braden Nicholson

import {inject} from "vue";
import type {Haptics} from "@/haptics";
import {useRouter} from "vue-router";
import type {Notify} from "@/notifications";
import type {Screensaver} from "@/screensaver";
import type {Remote} from "@/remote";


export default {
    router: () => useRouter(),
    remote: () => inject("remote") as Remote,
    screensaver: () => inject("screens") as Screensaver,
    haptics: () => inject("haptics") as Haptics,
    notify: () => inject("notify") as Notify,
}

