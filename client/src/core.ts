// Copyright (c) 2022 Braden Nicholson

import {inject} from "vue";
import type {Remote} from "@/types";
import type {Haptics} from "@/views/terminal/haptics";
import {useRouter} from "vue-router";


export default {
    router: () => useRouter(),
    remote: () => inject("remote") as Remote,
    haptics: () => inject("haptics") as Haptics,
    notifications: () => inject("notifications") as any,
}

