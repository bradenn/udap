<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import type {Attribute, Entity} from "udap-ui/types";
import {onMounted, reactive, watchEffect} from "vue";
import core from "udap-ui/core";
import attributeService from "@/services/attributeService";
import Element from "udap-ui/components/Element.vue";
import List from "udap-ui/components/List.vue";

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
  td: 80,
  hd: 0,
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
    state.set.heat = Math.round(state.set.heat * 100) / 100
  }

  let cool = state.attributes.find(a => a.key === "cool")
  if (cool) {
    state.set.cool = Math.round(parseFloat(cool.value) * 100) / 100
    state.set.cool = Math.round(state.set.cool * 100) / 100
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
  state.td = Math.round((state.inside.temp - state.set.cool) * 100) / 100
  state.hd = Math.round((state.outside.next - state.outside.temp) * 100) / 100
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
  <Element>
    <List :row="true" style="height: 4rem">
      <button id="install" hidden>Install</button>
      <Element class="d-flex align-items-center" foreground style="height: 3rem; aspect-ratio: 1.35914091/1;">
        <div class="label-c0 sf-icon label-cool">􀇥</div>
      </Element>
      <Element :cb=" ()=> router.push('/thermostat')" :foreground="true" :mutable="true"
               class="d-flex align-items-center justify-content-center w-25">
        <div class="label-c0 lh-1">{{ Math.round(state.inside.temp) }}</div>
      </Element>
      <Element foreground style="height: 3rem; aspect-ratio: 1.35914091/1;">
        <div class="label-c0 sf-icon label-heat d-flex flex-row">
          <div>􁰹</div>
          <div>􀆇</div>
        </div>
      </Element>

    </List>
  </Element>
</template>

<style lang="scss">

.label-cool {
  color: rgba(90, 200, 245, 0.9)
}

.label-heat {
  color: hsla(0, 50%, 55%, 0.9);
}

.font-monospaced {
  //font-family: "SF Pro", sans-serf !important;
  font-weight: 500;
}

.temp {
  letter-spacing: -2px;
  font-size: 3.8rem;
  font-weight: 400;
}

</style>