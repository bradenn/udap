// Copyright (c) 2022 Braden Nicholson

import type {Ref} from "vue";
import {ref} from "vue";


export interface Context {
    showContext(): void

    hideContext(): void

    isActive(): boolean
}

const show: Ref<boolean> = ref<boolean>(false)

function showContext() {
    show.value = true
}

function hideContext() {
    show.value = false
}

function isActive(): boolean {
    return show.value
}

export default {
    showContext,
    hideContext,
    isActive,
};



