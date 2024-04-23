// Copyright (c) 2023 Braden Nicholson

import {onMounted, reactive, watchEffect} from "vue";
import core from "@/core";
import {Module, Timing} from "udap-ui/types";
import {Remote} from "udap-ui/remote";

export interface ModuleMeta {
    id: string
    delta: number
    start: number
    stop: number
}

export interface ModuleMeta {
    id: string
    delta: number
    start: number
    stop: number
}

export interface ModuleTiming {
    module?: Module
    timing: Timing
}

export interface ModuleController {
    modules: Module[]
    timing: Timing[]
    moduleMeta: ModuleMeta[],
    timings: ModuleTiming[],
    loaded: boolean
}

export default function useListModules(): ModuleController {

    const remote: Remote = core.remote() as Remote;

    const state: ModuleController = reactive({
        modules: {} as Module[],
        moduleMeta: [] as ModuleMeta[],
        timing: [],
        timings: [],
        loaded: false
    }) as ModuleController

    onMounted(() => {
        state.modules = remote.modules
        state.timing = remote.timings
        reloadModuleMeta()
    })

    watchEffect(() => {
        state.modules = remote.modules
        return remote.modules
    })

    watchEffect(() => {
        state.timing = remote.timings
        reloadModuleMeta()
        return remote.timings
    })

    function reloadModuleMeta() {
        if (!state.modules || !state.timing) return
        let sorted: Timing[] = state.timing.sort((a: Timing, b: Timing) => a.startNano - b.startNano)
        state.timings = sorted.map((t: Timing) => {
                return {
                    timing: t,
                    module: state.modules.find(m => t.pointer === `module.${m.uuid}.update`)
                } as ModuleTiming
            }
        ) as ModuleTiming[]
        state.loaded = true;
    }


    return state
}