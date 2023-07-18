// Copyright (c) 2023 Braden Nicholson

import {onMounted, reactive, watchEffect} from "vue";
import core from "@/core";
import {Endpoint} from "udap-ui/types";


export interface EndpointController {
    endpoints: Endpoint[]


    loaded: boolean
}

export default function useListEndpoints(): EndpointController {

    const remote = core.remote();

    const state = reactive<EndpointController>({
        endpoints: {} as Endpoint[],

        loaded: false
    })

    onMounted(() => {
        state.endpoints = remote.endpoints.sort((a, b) => new Date(b.updated).valueOf() - new Date(a.updated).valueOf())
        state.loaded = true;
    })

    watchEffect(() => {
        state.endpoints = remote.endpoints
        return remote.modules
    })


    return state
}