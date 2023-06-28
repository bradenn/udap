<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import {onMounted, reactive} from "vue";
import {Attribute} from "@/types";
import core from "@/core";
import Element from "@/App.vue";

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

onMounted(() => {
  updateState()
})

function updateState() {
  state.attributes = remote.attributes.filter(a => a.key === "thermostat")

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


  state.min = Math.min(state.inside.temp, state.ecomin, state.ecomax, state.set.cool, state.set.heat, state.outside.temp, state.outside.next) - 5

  state.max = Math.max(state.inside.temp, state.ecomin, state.ecomax, state.set.cool, state.set.heat, state.outside.temp, state.outside.next) + 5

  state.td = Math.round((state.inside.temp - state.set.cool) * 100) / 100
  state.hd = Math.round((state.outside.next - state.outside.temp) * 100) / 100
  state.loaded = true;

}

</script>

<template>
  <Element>
    {{ state.outside.temp }}
  </Element>
</template>

<style lang="scss" scoped>

</style>