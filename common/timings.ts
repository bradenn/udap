// Copyright (c) 2022 Braden Nicholson

import {Remote} from "./remote";
import {computed, ComputedRef, onMounted, reactive, watchEffect} from "vue";
import {Timing} from "./types";

export interface ModuleTiming {
    run: ComputedRef<Timing>
    update: ComputedRef<Timing>
}

export interface RemoteTimings {
    timings: Timing[]
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
        getModuleTimings: getModuleTimings
    } as RemoteTimings)

    onMounted(() => {
        state.timings = remote.timings
    })

    watchEffect(() => {
        state.timings = remote.timings
        return remote.timings
    })

    function getModuleTimings(moduleUuid: string): ModuleTiming {
        let moduleRun = `module.${moduleUuid}.run`
        let moduleUpdate = `module.${moduleUuid}.update`

        return {
            run: computed(() => {
                return state.timings.find(t => t.pointer == moduleRun)
            }),
            update: computed(() => {
                return state.timings.find(t => t.pointer == moduleUpdate)
            }),
        } as ModuleTiming

    }

    return state
}

