<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {Attribute, Entity} from "@/types";
import {onBeforeMount, reactive, watchEffect} from "vue";
import core from "@/core";
import Slider from "@/components/Slider.vue";
import attributeService from "@/services/attributeService";


let state = reactive({
    entity: {} as Entity,
    attributes: [] as Attribute[],
    entityId: "",
    spectrum: {
        on: false,
        dim: 0,
        cct: 0,
        hue: 0,
    },
})

const router = core.router()
const remote = core.remote()
onBeforeMount(() => {
    state.entityId = router.currentRoute.value.params["entityId"] as string
    const entity = remote.entities.find(e => e.id === state.entityId)
    if (!entity) return
    state.entity = entity
    updateAttribute()
})
watchEffect(() => {
    updateAttribute()
    return remote.attributes
})

function updateAttribute() {
    state.attributes = remote.attributes.filter(a => a.entity === state.entity.id)

    let found = state.attributes.find(a => a.key === 'on')
    state.spectrum.on = found?.value === "true"

    found = state.attributes.find(a => a.key === 'dim')
    state.spectrum.dim = parseInt(found?.value || "0")

    found = state.attributes.find(a => a.key === 'cct')
    state.spectrum.cct = parseInt(found?.value || "0")

    found = state.attributes.find(a => a.key === 'hue')
    state.spectrum.hue = parseInt(found?.value || "0")
}

function sendRequest(key: string, value: string) {
    let found = state.attributes.find(a => a.key === key)
    if (!found) return
    found.request = value
    attributeService.request(found).then(_ => {

    }).catch(err => {
        console.log(err)
    })
}

function setDim(value: number) {
    sendRequest("dim", `${Math.max(value, 1)}`)
}

</script>

<template>
    <div class="d-flex gap-1 flex-column">
        <div class="d-flex align-items-center justify-content-start gap-2 header">
            <div class=" label-c2 sf-icon px-3">{{ state.entity.icon }}</div>
            <div>
                <div class="label-c5 label-w700 label-o5 ">{{
                    state.entity.alias ? state.entity.alias : state.entity.name
                    }}
                </div>
            </div>
        </div>
        <div>
            {{ state.spectrum.dim }}
            <Slider :change="(value) => sendRequest('dim', `${Math.max(value, 1)}`)" :max="100" :min="1"
                    :step="2" :value="state.spectrum.dim"
                    bg="dim"></Slider>
            {{ state.spectrum.cct }}
            <Slider :change="(value) => sendRequest('cct', `${Math.max(value, 2000)}`)" :max="6000" :min="2000"
                    :step="1" :value="state.spectrum.cct"
                    bg="cct"></Slider>
            {{ state.spectrum.hue }}
            <Slider :change="(value) => sendRequest('hue', `${Math.max(value, 0)}`)" :max="360" :min="1" :step="1"
                    :value="state.spectrum.hue" bg="hue"></Slider>
        </div>
    </div>
</template>

<style scoped>
.header {

    padding: 1rem 0.25rem;
    backdrop-filter: blur(40px);
    box-shadow: inset 0 0 1px 1.5px rgba(37, 37, 37, 0.6), 0 0 3px 1px rgba(22, 22, 22, 0.6);
    /* Note: backdrop-filter has minimal browser support */
    border-radius: 11.5px;
    -webkit-backdrop-filter: blur(40px) !important;
}
</style>