<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import type {Attribute, Entity} from "@/types";
import {onMounted, reactive, watchEffect} from "vue";
import core from "@/core";
import attributeService from "@/services/attributeService";
import TempSlider from "./TempSlider.vue";

const props = defineProps<{
  entity: Entity
}>()

const state = reactive({
  attributes: [] as Attribute[],
  mode: "",
  outside: {
    humidity: 1,
    temp: 68,
    next: 70,
  },
  inside: {
    humidity: 1,
    temp: 50,
  },
  min: 0,
  max: 100,
  emin: 40,
  emax: 80,
  ecomin: 20,
  ecomax: 80,
  set: {
    cool: 90,
    heat: 32,
  },
  loaded: false,
  connected: false
})

const router = core.router();
const remote = core.remote();

function toF(val: number): number {
  return val * (9 / 5) + 32
}

onMounted(() => {
  updateAttribute()
})

watchEffect(() => {
  updateAttribute()
  return remote.attributes
})

interface CurrentWeather {
  temperature: number
  temperatureNext: number
  humidity: number
  min: number
  max: number
}

function updateAttribute() {

  state.attributes = remote.attributes.filter(a => a.entity === props.entity.id)

  let temp = state.attributes.find(a => a.key === "thermostat")
  if (temp) {
    let val = JSON.parse(temp.value) as {
      connected: boolean
      temperature: number
      humidity: number
      ecoCool: number
      ecoHeat: number
      mode: string
    }
    state.inside.temp = Math.round((val.temperature) * 100) / 100
    state.inside.humidity = Math.round((val.humidity) * 100) / 100
    state.mode = val.mode
    state.ecomax = val.ecoCool
    state.ecomin = val.ecoHeat
  }

  let heat = state.attributes.find(a => a.key === "heat")
  if (heat) {
    state.set.heat = Math.round(parseFloat(heat.value) * 100) / 100
    state.set.heat = Math.max(state.set.heat, 55)
    state.set.heat = Math.round(state.set.heat / 0.5) * 0.5
  }

  let cool = state.attributes.find(a => a.key === "cool")
  if (cool) {
    state.set.cool = Math.round(parseFloat(cool.value) * 100) / 100
    state.set.cool = Math.round(state.set.cool / 0.5) * 0.5
  }

  let weather = remote.attributes.find(a => a.key === "weather")
  if (weather) {
    let cw = JSON.parse(weather.value) as CurrentWeather

    state.outside.temp = cw.temperature
    state.outside.next = cw.temperatureNext
    state.outside.humidity = cw.humidity
    state.emin = cw.min
    state.min = Math.min(state.inside.temp, state.ecomin, state.ecomax, state.set.cool, state.set.heat, state.outside.temp, state.outside.next, cw.min) - 5
    state.emax = cw.max
    state.max = Math.max(state.inside.temp, state.ecomin, state.ecomax, state.set.cool, state.set.heat, state.outside.temp, state.outside.next, cw.max) + 5
  }

  state.loaded = true

}

function sendRequest(key: string, value: string) {
  let found = state.attributes.find(a => a.key === key)
  if (!found) return
  found.request = value
  attributeService.request(found).then(suc => {
    console.log(suc)
  }).catch(err => {
    console.log(err)
  })
}

function setCool(delta: number) {

  let value = `${state.set.cool + delta}`

  sendRequest("cool", value)
}

function mouseHold(e: MouseEvent) {
  console.log("Hold")
}

</script>

<template>
  <div class="d-flex flex-column element p-1">

    <div v-if="true" class="py-2" style="width: 100%">
      <TempSlider :ecoMax="state.ecomax"
                  :ecoMin="state.ecomin"
                  :inside="state.inside.temp" :max="state.max"
                  :min="state.min" :outside="state.outside.temp"
                  :outside-next="state.outside.next" :set="state.set"
                  :weather-max="state.emax" :weather-min="state.emin" class="w-100"></TempSlider>
    </div>
    <div class="d-flex justify-content-between align-items-center gap-1 p-1">
      <div class="subplot d-flex justify-content-center align-items-center sf-icon"
           style="width: 4rem; height: 4rem"
           @click="(e) => setCool(-0.5)">􀆈
      </div>
      <div
          class="subplot flex-grow-1 d-flex flex-column align-items-center justify-content-center gap-0 sf-icon"
          style="height: 4rem">
        <div class="label-c0 label-w800 lh-1">{{ state.set.cool }}&deg; F</div>
        <div v-if="state.mode === 'COOL'" class="label-c6 label-w700 label-o5"
             style="color: rgba(90,200,245,0.9)">
          {{ state.mode }}
        </div>
        <div v-if="state.mode === 'HEAT'" class="label-c7 label-w700 label-o5"
             style="color: rgba(255,69,55,1)">
          {{ state.mode }}
        </div>
        <div v-if="state.mode === 'AUTO'" class="label-c7 label-w700 label-o5"
             style="color: rgba(90,90,90,0.9)">
          {{ state.mode }}
        </div>
        <div v-if="state.mode === 'OFF'" class="label-c7 label-w700 label-o5"
             style="color: rgba(90,90,90,0.9)">
          {{ state.mode }}
        </div>
      </div>
      <div class="subplot d-flex justify-content-center align-items-center sf-icon"
           style="width: 4rem; height: 4rem"
           @click="(e) => setCool(0.5)">
        􀆇
      </div>
    </div>
  </div>
</template>

<style lang="scss">

.subplot {
  border-radius: 0.3rem !important;
  //box-shadow: 0 0 8px 2px rgba(0, 0, 0, 0.05);
  //border: 1px solid rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.55);
  font-weight: 500;
  font-size: 1rem;
  background-color: rgba(255, 255, 255, 0.025);
}

.subplot-hidden {
  background-color: transparent !important;

}

.subplot-hidden:active {
  background-color: inherit !important;
}

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