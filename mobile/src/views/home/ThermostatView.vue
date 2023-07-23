<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import {onMounted, reactive, watchEffect} from "vue";
import {Attribute, Entity} from "udap-ui/types";
import core from "@/core";
import Element from "udap-ui/components/Element.vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";
import ElementLabel from "udap-ui/components/ElementLabel.vue";
import List from "udap-ui/components/List.vue";
import type {CurrentWeather} from "udap-ui/weather"
import type {Forecast} from "udap-ui/composables/weather";
import {useForecast} from "udap-ui/composables/weather";
import attributeService from "@/services/attributeService";


interface CurrentWeather {
  temperature: number
  temperatureNext: number
  humidity: number
  min: number
  max: number
}

const forecast = useForecast() as Forecast
const state = reactive({
  attributes: [] as Attribute[],
  thermostat: {} as Entity,
  expandControls: false,
  mode: "OFF",
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
  properties: {
    updateTime: "",
    generatedAt: "",
  },
  ranges: {
    temp: {
      min: 100,
      max: -100,
    },
    wind: {
      min: 100,
      max: 0,
    }, rain: {
      min: 0,
      max: 0,
    },
  },
  loaded: false,
  connected: false
})


const router = core.router();
const remote = core.remote();

onMounted(() => {
  updateState()
})

watchEffect(() => {
  updateState()
  return remote.attributes
})

function updateState() {

  let th = remote.entities.find(e => e.type === "thermostat" && e.name === "thermostat")
  if (!th) return
  state.thermostat = th
  let attrs = remote.attributes.filter(a => state.thermostat.id === a.entity)
  if (!attrs) return;
  state.attributes = attrs
  let mode = state.attributes.find(a => a.key === "mode")
  if (mode) {
    state.mode = mode.value;
  }
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
    // state.mode = val.mode
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

  let value = `${Math.round(state.set.cool + delta)}`
  if (state.mode !== "OFF") {
    sendRequest("cool", value)
  }
}

function setMode(mode: string) {
  sendRequest("mode", mode)
}

</script>

<template>
  <div>
    <Element>
      <List>
        <List row>

          <Element :foreground="true" :long-cb="() => router.push('/thermostat')" :mutable="true"
                   class="w-100  px-4">
            <div class="d-flex flex-column align-items-center gap-0 justify-content-center w-100">
              <div class="d-flex align-items-center gap-0">
                <div class="d-flex">
                  <div class="temp lh-1">{{ state.inside.temp.toString().split(".")[0] }}</div>
                </div>
                <div class="d-flex flex-column justify-content-between" style="height: 3rem;">
                  <div class="label-o3 label-c5 ">&nbsp;&deg;F</div>
                  <div class="label-c5 label-o2" style="margin-bottom: -2px;">.{{
                      state.inside.temp.toString().split(".")[1]
                    }}
                  </div>
                </div>

              </div>

            </div>
          </Element>

          <div class="w-100  px-4 d-flex gap-2 py-2">
            <!--        <div style="width: 2px; border-radius: 1px; height: 3.7rem; background-color: rgba(255,255,255,0.04)"></div>-->
            <div class="d-flex flex-column gap-2 align-items-start justify-content-between font-monospaced"
                 style="width: 4.6rem">
              <div class=" lh-1 d-flex justify-content-between gap-1">
                <div class="label-c6 label-o3" style="width: 2rem">cool</div>
                <div class="label-c6 label-o3"
                     style="color: hsla(200, 50%, 55%, 0.9);">{{ state.set.cool.toFixed(2) }}&deg;
                </div>
              </div>
              <div v-if="false" class="d-flex justify-content-between lh-1 gap-1">
                <div class="label-c6 label-o3" style="width: 2rem">heat</div>
                <div class="label-c6 label-o3" style="color: hsla(0, 50%, 55%, 0.9);">{{ state.set.heat.toFixed(2) }}&deg;
                </div>
              </div>

              <div class=" lh-1 d-flex justify-content-between gap-1">
                <div class="label-c6 label-o3" style="width: 2rem; margin-bottom: -4px">H2O<sub
                    style="font-size: 8px">i</sub>
                </div>
                <div class="label-c6 label-o3"
                     style="color: hsla(200, 20%, 50%, 0.9);">{{ state.inside.humidity.toFixed(0) }}%
                </div>
              </div>
              <div class="d-flex justify-content-between lh-1 gap-1">
                <div class="label-c6 label-o3" style="width: 2rem">&Delta;&tau;<sub class=""
                                                                                    style="font-size: 8px">i</sub>
                </div>
                <div class="label-c6 label-o3" style="color: hsla(0, 0%, 50%, 1);">{{
                    state.td > 0 ? `+${state.td.toFixed(2)}` : state.td.toFixed(2)
                  }}&deg;
                </div>
              </div>
            </div>
            <div class=""
                 style="width: 2px; border-radius: 1px; height: 3.7rem; background-color: rgba(255,255,255,0.04)"></div>
            <div class="d-flex flex-column gap-2 align-items-start justify-content-between font-monospaced"
                 style="width: 5rem">
              <div class=" lh-1 d-flex justify-content-between gap-1">
                <div class="label-c6 label-o3" style="width: 2.8rem">outside
                </div>
                <div class="label-c6 label-o4"
                >{{ state.outside.temp.toFixed(2) }}&deg;
                </div>
              </div>
              <div class=" lh-1 d-flex justify-content-between gap-1">
                <div class="label-c6 label-o3" style="width: 2.8rem;  margin-bottom: -4px">H2O<sub
                    style="font-size: 8px">o</sub>
                </div>
                <div class="label-c6 label-o3"
                     style="color: hsla(200, 20%, 50%, 0.9);">{{ state.outside.humidity.toFixed(0) }}%
                </div>
              </div>
              <!--          <div class=" lh-1 d-flex justify-content-between gap-1">-->
              <!--            <div class="label-c6 label-o3" style="width: 2.8rem">high-->
              <!--            </div>-->
              <!--            <div class="label-c6 label-o4"-->
              <!--            >{{ state.ecomax.toFixed(2) }}&deg;-->
              <!--            </div>-->
              <!--          </div>-->
              <!--          <div class=" lh-1 d-flex justify-content-between gap-1">-->
              <!--            <div class="label-c6 label-o3" style="width: 2.8rem">low-->
              <!--            </div>-->
              <!--            <div class="label-c6 label-o4"-->
              <!--            >{{ state.ecomin.toFixed(2) }}&deg;-->
              <!--            </div>-->
              <!--          </div>-->
              <div class="d-flex justify-content-between lh-1 gap-1">
                <div class="label-c6 label-o3" style="width: 2.8rem">&Delta;&tau;<sub class=""
                                                                                      style="font-size: 8px">o</sub>
                  1hr
                </div>
                <div :style="state.hd > 0?`color: hsla(0, 50%, 55%, 0.9);`:`color: hsla(200, 50%, 55%, 0.9);`"
                     class="label-c6 label-o3">
                  {{ state.hd > 0 ? `+${state.hd.toFixed(2)}` : state.hd.toFixed(2) }}&deg;
                </div>
              </div>

            </div>
          </div>


        </List>
        <List scroll-y style="max-height: 74vh">
          <ElementHeader title="Controls"></ElementHeader>
          <List row style="height: 3.25rem; z-index: 1">
            <Element :cb="() => setCool(-1)" class="align-items-center d-flex justify-content-center  flex-row"
                     foreground
                     mutable>
              <div class="label-c2 lh-1 label-w600 label-o3 sf-icon">􀆈
              </div>

            </Element>
            <Element :cb="() => state.expandControls = !state.expandControls"
                     class="align-items-center d-flex justify-content-center flex-column" foreground
                     style="width: 25%">
              <div
                  v-if="state.mode ==='COOL'" class="d-flex flex-column align-items-center justify-content-center">
                <div class="label-c2 lh-1 label-w500 label-o4">{{ state.set.cool }}&deg;</div>
                <div class="label-c6 lh-1" style="color: hsla(200, 50%, 55%, 0.9);">cool to</div>
              </div>
              <div v-else-if="state.mode === 'OFF'"
                   class="align-items-center d-flex justify-content-center flex-column">
                <div class="label-c2 lh-1 label-w500 label-o4">OFF</div>
                <div class="label-c6 lh-1 sf" style="color: hsla(200, 50%, 55%, 0.9);">􀍠</div>

              </div>
            </Element>
            <Element :cb="() => setCool(1)" class="align-items-center d-flex justify-content-center sf-icon" foreground
                     mutable>
              <div class="label-c2 lh-1 label-w600 label-o3 sf-icon">􀆇</div>
            </Element>
          </List>
          <List v-if="state.expandControls" row style="height: 3.25rem; z-index: 1">
            <Element :accent="state.mode === 'OFF'" :cb="() => setMode('OFF')"
                     class="align-items-center d-flex justify-content-center  flex-row"
                     foreground mutable
            >
              <div class="label-c5 label-w400 label-o3">OFF</div>
            </Element>
            <Element :accent="state.mode === 'COOL'" :cb="() => setMode('COOL')"
                     class="align-items-center d-flex justify-content-center  flex-row"
                     foreground mutable
            >
              <div class="label-c5 label-w400 label-o3">COOL</div>
            </Element>
            <Element :accent="state.mode === 'HEAT'" :cb="() => setMode('HEAT')"
                     class="align-items-center d-flex justify-content-center  flex-row"
                     foreground mutable
            >
              <div class="label-c5 label-w400 label-o3">HEAT</div>
            </Element>
            <Element :accent="state.mode === 'HEATCOOL'" :cb="() => setMode('HEATCOOL')"
                     class="align-items-center d-flex justify-content-center  flex-row"
                     foreground mutable
            >
              <div class="label-c5 label-w400 label-o3">HEAT COOL</div>
            </Element>
          </List>
          <List>
            <ElementHeader title="Inside"></ElementHeader>
            <List class="">
              <ElementLabel icon="􁏄" title="Current">{{ state.inside.temp }}&deg;</ElementLabel>
              <ElementLabel icon="􁃛" title="Humidity">{{ state.inside.humidity }}%</ElementLabel>

            </List>

            <ElementHeader title="Outside"></ElementHeader>
            <List class="">
              <ElementLabel icon="􁏄" title="Current">{{ state.outside.temp }}&deg;</ElementLabel>
              <ElementLabel icon="􁃛" title="Humidity">{{ state.outside.humidity }}%</ElementLabel>
              <List row>
                <ElementLabel icon="􀇫" style="border-top-left-radius: 0.375rem !important;" title="Low">{{
                    state.emin
                  }}&deg;
                </ElementLabel>
                <ElementLabel icon="􀦜" style="border-top-right-radius: 0.375rem !important;" title="High">{{
                    state.emax
                  }}&deg;
                </ElementLabel>
              </List>

            </List>

            <ElementHeader title="Forecast"></ElementHeader>
            <List v-if="!forecast.loading" class="" row>

              <Element v-for="item in forecast.hourly"
                       class="sf-icon d-flex flex-column align-items-center justify-content-center"
                       foreground>
                <div class="label-c6 label-o2">{{ item.time }}</div>
                <div class="label-c5 label-o4 label-w500 py-1">{{ item.icon }}</div>
                <div class="label-c5 label-o4 label-w500">{{ item.temp }}&deg;</div>
              </Element>
            </List>


          </List>
        </List>
      </List>
    </Element>
  </div>

</template>

<style lang="scss" scoped>

</style>