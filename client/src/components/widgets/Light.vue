<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import AttributeComponent from "@/components/entity/Attribute.vue"
import type {ApiRateLimit, Attribute, Entity, Remote} from "@/types"
import moment from "moment";
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
    <div v-if="state.showMenu" class="context context-light" @click="toggleMenu"></div>
    <div v-if="state.loading" class="w-100 h-100">
        <div class="entity-small element">
            <div class="entity-header mb-2 ">
                <div class="label-o5">
                    {{ props.entity.icon }}
                </div>
                <div class="label-c1 label-w400 label-o4 px-2">
                    {{ props.entity.name }}
                </div>
                <div class="fill"></div>

            </div>
        </div>
    </div>
    <div v-else class="w-100 h-100">
        <div v-if="!state.showMenu" class="element d-flex justify-content-start align-items-center "
             style="height: 2.125em !important; margin-bottom: 0rem; width: 100%"
             @click="toggleMenu">
            <div class="d-flex justify-content-start align-items-center align-content-center ">
                <div
                        :style="`width: 0.8rem; margin-left:0.25rem; color:${state.activeColor}; height: 100%;border-radius: 0.25rem; `"

                        class="label-c2 light-on d-flex align-items-center label-w700 justify-content-center">
                    {{ (props.entity.icon || '􀛭') }}
                </div>
                <div class="d-flex flex-column">
                    <div class="label-c2 label-o3 label-w500 px-1 lh-1">{{ props.entity.name }}</div>
                    <div class="label-c3 label-o2  label-w500 px-1 lh-1">{{ state.shortStatus }}</div>
                </div>
            </div>
        </div>
        <div v-else class="entity-small element quickedit" style="z-index: 99999!important;">
            <div class="entity-header mb-1 d-flex justify-content-between align-items-center w-100">
                <div class="d-flex">
                    <div class="label-o5">
                        {{ props.entity.icon || '􀛮' }}
                    </div>
                    <div class="label-c1 label-w400 label-o4 px-2">
                        {{ props.entity.name }}
                    </div>
                </div>
                <AttributeComponent :attribute="state.powerAttribute" :entity-id="props.entity.id"
                                    small></AttributeComponent>
                <!--        <div class="" @click="toggleMenu">-->
                <!--          <i class="fa-solid fa-circle-xmark label-o3 label-c1 label-w400 px-1"></i>-->
                <!--        </div>-->
            </div>
            <div class="w-100 d-flex flex-column gap">
                <div
                        v-for="attribute in state.attributes.filter((attr: Attribute) => attr.key === 'dim' || attr.key === 'cct' || attr.key === 'hue')"
                        :key="attribute.id">
                    <AttributeComponent :attribute="attribute" :entity-id="props.entity.id" primitive
                                        small></AttributeComponent>
                </div>
            </div>
            <div v-if="false" class="entity-small element ">
                <div class="element surface d-flex flex-column gap py-4 px-3 pt-2">
                    <div class="d-flex justify-content-start align-items-end align-content-end" v-on:click.stop>
                        <div class="mt-1">
              <span :style="`text-shadow: 0 0 8px ${state.activeColor};`"
                    class="label-md label-w600 label-o3">{{ props.entity.icon }}</span>
                            <span class="label-md label-w600 label-o6 px-2">{{ props.entity.name }}</span>
                        </div>
                        <div class="fill "></div>

                        <div class="h-bar">
                            <AttributeComponent :attribute="state.powerAttribute" :entity-id="props.entity.id"
                                                small></AttributeComponent>
                        </div>
                    </div>
                    <div class="h-sep"></div>
                    <div
                            v-for="attribute in state.attributes.filter((a: Attribute) => a.key !== 'on')">
                        <AttributeComponent :key="attribute.id" :attribute="attribute" :entity-id="props.entity.id"
                                            primitive
                                            small></AttributeComponent>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style lang="scss" scoped>

.macro-grid {

}


.light-on {
  //color: rgba(178, 178, 178, 0.76);
  background-blend-mode: luminosity;
  //text-shadow: 0 0 8px rgba(255, 255, 255, 0.1);
}

.light-off {
  color: rgba(255, 255, 255, 0.3);
}

.entity-small:not(.quickedit):active {
  animation: click 100ms ease forwards;
}

.entity-small {
  animation: click 100ms ease forwards;
}


@keyframes click {
  0% {
    transform: scale(1.0);
  }
  25% {
    transform: scale(0.98);
  }
  30% {
    transform: scale(0.97);
  }
  100% {
    transform: scale(1);
  }
}

.entity-context {
  position: absolute;
  backdrop-filter: blur(22px);
  top: 0;
  width: 100%;
  padding: 1rem;
  height: calc(100% - 4.5rem);
}


</style>