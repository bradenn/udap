// Copyright (c) 2023 Braden Nicholson

import {onMounted, reactive, watchEffect} from "vue";
import {Remote} from "../remote";
import core from "../core";
import {Attribute, Entity, Zone} from "../types";

export interface Light {
    name: string
    icon: string
    type: string
    on: boolean
    dim: number
    cct: number
    hue: number
}

export interface Lights {
    zones: Zone[]
    lights: Light[]
    entities: Entity[]
    attributes: Attribute[]
    loaded: boolean
}

export function useLights(): Lights {
    const state: Lights = reactive({
        zones: [] as Zone[],
        lights: [] as Light[],
        entities: [] as Entity[],
        attributes: [] as Attribute[],
        loaded: false
    } as Lights)

    const remote: Remote = core.remote() as Remote


    onMounted(() => {
        prepareLights()
    })

    watchEffect(() => {
        prepareLights()
        return remote.attributes
    })


    function prepareLights() {
        state.zones = remote.zones.filter(z => z.pinned).filter(z => {
            return z.entities.find(e => e.type === "spectrum")
        })
        state.loaded = true
        state.entities = remote.entities.filter(e => e.type === "spectrum")
        state.attributes = remote.attributes.filter(a => state.entities.find(e => e.id === a.entity))
        // state.lights = state.entities.map(e => {
        //     return entityToLight(e) as Light
        // })
    }

    function entityToLight(entity: Entity): Light {

        return {} as Light
    }

    return state
}

