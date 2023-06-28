// Copyright (c) 2023 Braden Nicholson

import {onMounted, reactive, watchEffect} from "vue";
import core from "@/core";
import {Module} from "udap-ui/types";


export interface ModuleController {
    modules: Module[]
    loaded: boolean
}

export default function useListModules(): ModuleController {

    const remote = core.remote();

    const state = reactive({
        modules: {} as Module[],
        loaded: false
    })

    onMounted(() => {
        state.modules = remote.modules
    })

    watchEffect(() => {
        state.modules = remote.modules
        return remote.modules
    })


    return state
}