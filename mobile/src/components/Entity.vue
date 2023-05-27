<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import type {Attribute, Entity} from "@/types";
import {onMounted, reactive, watchEffect} from "vue";
import core from "@/core";
import attributeService from "@/services/attributeService";

const props = defineProps<{
    aspect?: number
    secondary?: boolean
    entity: Entity
}>()

const state = reactive({
    attributes: [] as Attribute[],
    spectrum: {
        on: false,
        dim: 0,
        cct: 0,
        hue: 0,
    },
    pressed: false,
    holding: false,
    timeout: 0,
    down: 0,
})

const router = core.router();
const remote = core.remote();

onMounted(() => {
    updateAttribute()
})

watchEffect(() => {
    updateAttribute()
    return remote.attributes
})

function updateAttribute() {
    state.attributes = remote.attributes.filter(a => a.entity === props.entity.id)
    let found = state.attributes.find(a => a.key === 'on')
    if (!found) return;
    state.spectrum.on = found.value === "true"
}

function sendRequest(key: string, value: string) {
    let found = state.attributes.find(a => a.key === key)
    if (!found) return
    found.request = value
    attributeService.request(found).then(_ => {

    }).catch(err => {

    })
}

function powerOn(id: string) {
    sendRequest("on", "true")
}

function powerOff(id: string) {
    sendRequest("on", "false")
}


function togglePower() {
    if (state.spectrum.on) {
        sendRequest("on", "false")
    } else {
        sendRequest("on", "true")
    }
}

function goEdit() {
    router.push(`/home/entity/${props.entity.id}`)
}

function mouseHold(e: MouseEvent) {
    console.log("Hold")
}

function mouseDown(e: TouchEvent) {
    e.preventDefault();
    state.down = Date.now();
    state.holding = false;
    state.pressed = true
    state.timeout = window.setTimeout(() => {
        state.holding = true;
        goEdit();
    }, 250); // Adjust the timeout as per your requirement
}

function mouseUp(e: TouchEvent) {
    e.preventDefault();
    clearTimeout(state.timeout);
    state.pressed = false
    const elapsedTime = Date.now() - state.down;

    if (elapsedTime < 250 && !state.holding) {
        togglePower();
    }
}

</script>

<template>
    <div :class="`${state.pressed?'pressed':''} ${state.spectrum.on?'on':'off'}`" class="entity px-3 py-2"
         @touchend="mouseUp"
         @touchstart="mouseDown">
        <div class="d-flex gap-2 justify-content-between align-items-center py-1 gap-3 w-100">
            <div class="d-flex justify-content-center align-items-center gap-3">
                <div :class="`${state.spectrum.on?'icon-on':'icon-off'}`" class="sf-icon">{{
                    props.entity.icon
                    }}
                </div>
                <div>
                    <div class="label-c5 label-w700  name">{{
                        props.entity.alias ? props.entity.alias : props.entity.name
                        }}
                    </div>

                </div>
            </div>
            <div v-if="false" class="d-flex">
                <div v-for="attr in state.attributes">
                    <div v-if="attr.key === 'on'" class="d-flex justify-content-center" style="width: 4rem">
                        <div :class="``"
                             class="power-toggle label-w500 label-c2"
                             @click="attr.value === 'false'?powerOn(attr.id):powerOff(attr.id)">
                            {{ attr.value === 'false' ? "OFF" : "ON" }}
                        </div>
                    </div>
                </div>

                <router-link :to="`/home/entity/${props.entity.id}`" class="sf-icon label-o3">􀆊</router-link>
            </div>
        </div>
    </div>

</template>

<style lang="scss">


.entity.pressed {
  transform: scale(0.99); /* Scale down the button when pressed */
  box-shadow: inset 0 0 5px rgba(0, 0, 0, 0.1);
}

.on {
  background-color: rgba(46, 46, 48, 0.4);

  .sf-icon {
    color: rgba(255, 255, 255, 0.9) !important;
    filter: drop-shadow(0 0 4px rgba(255, 255, 255, 0.6));
  }

  .entity-status {
    color: rgba(255, 255, 255, 0.8) !important;
    filter: drop-shadow(0 0 2px rgba(255, 255, 255, 0.4));
  }
}

.off {
  background-color: rgba(22, 22, 22, 0.2);

  .sf-icon {
    color: rgba(255, 255, 255, 0.1) !important;
    filter: drop-shadow(0 0 4px rgba(255, 255, 255, 0.6));
  }

  .entity-status {
    color: rgba(255, 255, 255, 0.2) !important;
    filter: drop-shadow(0 0 4px rgba(255, 255, 255, 0.6));
  }
}

.entity {


  backdrop-filter: blur(40px);
  box-shadow: inset 0 0 1px 1.5px rgba(37, 37, 37, 0.6), 0 0 3px 1px rgba(22, 22, 22, 0.6);
  /* Note: backdrop-filter has minimal browser support */
  aspect-ratio: 2.71828183/0.8;
  border-radius: 11.5px;
  -webkit-backdrop-filter: blur(40px) !important;
  display: flex;
  justify-content: center;

  .sf-icon {
    font-size: 1rem;
    /* Label Color/Light/Primary */
    color: #FFF;

    mix-blend-mode: overlay;
  }

  .name {

    /* Label Color/Light/Primary */
    color: rgba(255, 255, 255, 0.8);

    mix-blend-mode: overlay;
  }

}

.surface {
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 0.5rem;
  padding: 1rem 0.25rem;
}
</style>