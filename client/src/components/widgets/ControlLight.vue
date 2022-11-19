<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {onMounted, reactive, watchEffect} from "vue";
import type {ApiRateLimit, Attribute, Entity} from "@/types"
import moment from "moment";
import Context from "@/views/terminal/Context.vue";
import Header from "@/components/Header.vue";
import core from "@/core";
import Switch from "@/components/Switch.vue";
import Slider from "@/components/Slider.vue";
import attributeService from "@/services/attributeService";
import type {Notify} from "@/notifications";
// Establish a local reactive state
let state = reactive({
    loading: true,
    active: false,
    showMenu: false,
    activeColor: "rgba(255,255,255,1)",
    shortStatus: "",
    levelAttribute: {} as Attribute,
    powerAttribute: {} as Attribute,
    apiAttribute: {} as Attribute,
    rateLimit: {} as ApiRateLimit,
    attributes: [] as Attribute[]
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
const remote = core.remote()


onMounted(() => {
    updateLight(remote.attributes)
    generateState()
    state.loading = true
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
    if (!hue) return;
    let cct = state.attributes.find((a: Attribute) => a.key === 'cct')
    if (!cct) return;
    let dim = state.attributes.find((a: Attribute) => a.key === 'dim')
    if (!dim) return;
    if (hue && cct && dim) {
        if (moment(hue).isAfter(cct.requested)) {
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
    const attrs = attributes.filter((a: Attribute) => a.entity === props.entity.id).sort(compareOrder)
    if (!attrs) return state.attributes
    state.attributes = attrs as Attribute[]
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

const notify: Notify = core.notify()

function applyState(attribute: Attribute) {
    attributeService.request(attribute).then(res => {
        notify.success("Core", `Changed ${props.entity.name}::${attribute.key} to '${attribute.request}'`)
    }).catch(err => {
        notify.fail("Core", `Could not change '${props.entity.name}' attribute '${attribute.key}' to '${attribute.request}'`, err)
    })
}

const colorTemps = [
    {
        name: "Warm",
        value: 2400,
        color: "255, 157, 63"
    },
    {
        name: "Tungsten",
        value: 3000,
        color: "255, 180, 107"
    },
    {
        name: "Moonlight",
        value: 4100,
        color: "255, 211, 168"
    }, {
        name: "Daylight",
        value: 5000,
        color: "255, 236, 224"
    },
    {
        name: "Cloudy",
        value: 5500,
        color: "255, 236, 224"
    },
    {
        name: "White",
        value: 6500,
        color: "255, 249, 253"
    }
]

function timeSince(time: string): string {
    return moment(time).fromNow()
}

function setCCT(attribute: Attribute) {
    attributeService.request(attribute).then(res => {
        notify.success("Core", `Changed ${props.entity.name}::${attribute.key} to '${attribute.request}'`)
    }).catch(err => {
        notify.fail("Core", `Could not change '${props.entity.name}' attribute '${attribute.key}' to '${attribute.request}'`, err)
    })
}

</script>

<template>

    <Context>
        <div class="control-light" style="grid-column: 5 / span 7; grid-row: 2 / span 4;" @click.stop>
            <Header :icon="props.entity.icon" :name="props.entity.name"></Header>
            <div class="">
                <div class="d-flex gap-1 flex-column">
                    <div v-for="attribute in state.attributes" :key="`${attribute.id}`">
                        <Switch v-if="attribute.key === 'on'"
                                :active="attribute.value === 'true'"
                                :change="(val) => {attribute.request = `${val?'true':'false'}`; applyState(attribute)}"
                                name="Power"></Switch>
                        <Slider v-else-if="attribute.key === 'dim'"
                                :change="(val) => {attribute.request = `${val}`; applyState(attribute)}" :live="false"
                                :max="100"
                                :min="0"
                                :step="5"
                                :value="parseInt(attribute.request)" name="Brightness" unit="%"></Slider>
                        <div v-else-if="attribute.key === 'cct'" class="">
                            <Slider :change="(val) => {attribute.request = `${colorTemps[val].value}`; setCCT(attribute)}"
                                    :live="false"
                                    :max="colorTemps.length-1" :min="0" :step="1"
                                    :tags="colorTemps.map(c => c.name)"
                                    :value="0" name="Color Temperature" unit="K"></Slider>

                        </div>
                        <div v-else-if="attribute.key === 'api'">

                        </div>
                    </div>
                </div>
            </div>

        </div>
    </Context>

</template>

<style lang="scss" scoped>

.control-light {

  transform: translateZ(0);
  z-index: 4 !important;
}
</style>