<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {ApiRateLimit, Attribute, Entity} from "@/types"
import moment from "moment";
import type {Remote} from "@/remote";
import ControlLight from "@/components/widgets/ControlLight.vue";
// Establish a local reactive state
let state = reactive<{
    loading: boolean,
    active: boolean,
    showMenu: boolean,
    activeColor: string,
    powerAttribute: Attribute,
    shortStatus: string
    levelAttribute: Attribute,
    apiAttribute: Attribute,
    attributes: Attribute[],
    rateLimit: ApiRateLimit
}>({
    loading: true,
    active: false,
    showMenu: false,
    activeColor: "rgba(255,255,255,1)",
    shortStatus: "",
    levelAttribute: {} as Attribute,
    powerAttribute: {} as Attribute,
    apiAttribute: {} as Attribute,
    rateLimit: {} as ApiRateLimit,
    attributes: []
})


function generateState() {
    if (state.levelAttribute.value) {
        state.shortStatus = (state.active) ? `ON${(state.levelAttribute) ? ', ' + state.levelAttribute.value + '%' : ''}` : 'OFF'
    } else {
        state.shortStatus = (state.active) ? `ON` : 'OFF'
    }
}

// Define the prop for this entity
let props = defineProps<{
    entity: Entity
}>()

// Inject the remote manifest
let remote = inject('remote') as Remote

// When the view loads, force the local state to update
onMounted(() => {
    state.loading = true
    updateLight(remote.attributes)
    generateState()
})

// Ripple any attribute changes down to the local reactive state
watchEffect(() => {
    updateLight(remote.attributes)
    findMode()
})


// Compare the defined order to sort the lights
function compareOrder(a: any, b: any): number {
    return a.order - b.order
}

function map_range(value: number, low1: number, high1: number, low2: number, high2: number) {
    return low2 + (high2 - low2) * (value - low1) / (high1 - low1);
}

function cctToRgb(cct: number) {
    return [map_range(cct, 2000, 8000, 255, 227),
        map_range(cct, 2000, 8000, 138, 233),
        map_range(cct, 2000, 8000, 18, 255)]
}

function findMode() {
    let hue = state.attributes.find((a: Attribute) => a.key === 'hue')
    let cct = state.attributes.find((a: Attribute) => a.key === 'cct')
    let dim = state.attributes.find((a: Attribute) => a.key === 'dim')
    if (hue && cct && dim) {
        if (moment(hue.requested).isAfter(cct.requested)) {
            state.activeColor = `hsla(${hue.value}, 100%, ${20 + parseInt(dim?.value) / 100.0 * 50}%, 0.5)`
        } else {
            state.activeColor = `rgba(${(cctToRgb(parseInt(cct.value)))[0]}, ${(cctToRgb(parseInt(cct.value)))[1]}, ${(cctToRgb(parseInt(cct.value)))[2]}, 0.5)`
        }
    } else {
        state.activeColor = 'rgba(255,255,255,0.5)'
    }
}

// Update the reactive model for the light
function updateLight(attributes: Attribute[]): Attribute[] {
    // Define the attributes for the light
    state.attributes = attributes.filter((a: Attribute) => a.entity === props.entity.id).sort(compareOrder)
    // Get the current power state of the light
    let on = state.attributes.find((a: Attribute) => a.key === 'on')
    let dim = state.attributes.find((a: Attribute) => a.key === 'dim')

    // Assign the power state to a local attribute
    if (!on) return []
    state.powerAttribute = on as Attribute
    state.active = on.value === "true" || on.request === "true"
    if (dim) {
        if (!state.active) {
            state.active = moment(dim.requested).isAfter(on.requested)
        }
        state.levelAttribute = dim as Attribute
    }
    generateState()
    state.loading = false
    let api = state.attributes.find((a: Attribute) => a.key === 'api')
    if (api) {
        state.apiAttribute = api
        state.rateLimit = JSON.parse(api.value) as ApiRateLimit
    } else {
        state.apiAttribute = {} as Attribute
    }

    return attributes
}

// Toggle the state of the context menu
function toggleMenu(): void {
    state.showMenu = !state.showMenu
    // context(state.showMenu)
}


</script>

<template>
    <ControlLight v-if="state.showMenu" :entity="props.entity" @click="state.showMenu = false"></ControlLight>
    <div class="element light" @click="state.showMenu = !state.showMenu">
        <div :style="`color: ${state.activeColor}!important;`" class="icon ">
            {{ (props.entity.icon || '􀛭') }}
        </div>
        <div class="metadata">
            <div class="name label-o5">{{ props.entity.name }}</div>
            <div class="description label-o3">{{ state.shortStatus }}</div>
        </div>
    </div>
</template>

<style lang="scss" scoped>

.element.light {

  z-index: 1 !important;

  display: flex;
  align-items: center;
  padding: 0.25rem 0.25rem;
  gap: 0.3rem;

  .icon {
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: rgba(255, 255, 255, 0.025);
    border-radius: 8px;
    aspect-ratio: 1/1 !important;
    width: 48px;
    font-size: 20px;
    line-height: 18px;
  }

  .metadata {
    display: flex;
    flex-direction: column;
    justify-content: center;

    .name {
      font-size: 18px;
      line-height: 18px;
      font-weight: 500;
      font-family: "SF Pro Rounded", sans-serif;
    }

    .description {
      font-size: 15px;
      line-height: 20px;
      font-weight: 500;
      font-family: "SF Pro Rounded", sans-serif;
    }
  }
}

</style>