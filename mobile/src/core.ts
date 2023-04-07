// Copyright (c) 2022 Braden Nicholson

import {inject} from "vue";
import type {Router} from "vue-router";
import {useRouter} from "vue-router";
import type {Remote} from "@/remote";


export interface Core {
    router: () => Router,
    remote: () => Remote,
}

export default {
    router: () => useRouter(),
    remote: () => inject("remote") as Remote,
} as Core

