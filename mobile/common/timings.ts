// Copyright (c) 2022 Braden Nicholson

import {Remote} from "./remote";
import {onMounted, reactive, watchEffect} from "vue";
import {Timing} from "./types";

export interface ModuleTiming {
    loaded: boolean,
    run: Timing
    update: Timing
}

export interface RemoteTimings {
    timings: Timing[]
    loaded: false,
    getModuleTimings: (moduleUuid: string) => ModuleTiming
}

export function convertNanosecondsToString(nanoseconds: number): string {
    if (nanoseconds < 0) {
        throw new Error("Input must be a non-negative number.");
    }

    if (nanoseconds === 0) {
        return "0ns";
    }

    const units = ["ns", "µs", "ms", "s", "ks"];
    let value = nanoseconds;
    let unitIndex = 0;

    while (value >= 1000 && unitIndex < units.length - 1) {
        if (unitIndex >= 3) {
            let minutes = value >= 60
            return `${minutes}`;
        }
        value /= 1000;
        unitIndex++;
    }

    return `${value.toFixed(2)}${units[unitIndex]}`;
}

export default function useTimings(remote: Remote): RemoteTimings {

    const state: RemoteTimings = reactive({
        timings: [] as Timing[],
        loaded: false,
        // getModuleTimings: getModuleTimings
    } as RemoteTimings)

    onMounted(() => {
        state.timings = remote.timings
    })

    watchEffect(() => {
        state.timings = remote.timings

        return remote.timings
    })

    // function getModuleTimings(moduleUuid: string): ModuleTiming {
    //     let moduleRun = `module.${moduleUuid}.run`
    //     let moduleUpdate = `module.${moduleUuid}.update`
    //
    //     let stub: Timing = {
    //         pointer: "",
    //         name: "",
    //         start: new Date().toString(),
    //         startNano: new Date().valueOf() * 1000,
    //         stop: new Date().toString(),
    //         stopNano: new Date().valueOf() * 1000,
    //         delta: 0,
    //         frequency: 0,
    //         complete: true,
    //         depth: 0,
    //         id: "id"
    //     }
    //
    //     return {
    //         loaded: false,
    //         run: computed(() => {
    //             return state.timings.find(t => t.pointer == moduleRun)
    //         }),
    //         update: computed(() => {
    //             return state.timings.find(t => t.pointer == moduleUpdate)
    //         }),
    //     } as ModuleTiming
    //
    // }

    return state
}

