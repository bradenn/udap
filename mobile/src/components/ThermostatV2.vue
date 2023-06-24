<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import type {Attribute, Entity} from "udap-ui/types";
import {onMounted, reactive, watchEffect} from "vue";
import core from "@/core";
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
    <List :row="true">


      <Element :foreground="true" :mutable="true" class="w-100  px-4">
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
                                                                                  style="font-size: 8px">o</sub> 1hr
            </div>
            <div :style="state.hd > 0?`color: hsla(0, 50%, 55%, 0.9);`:`color: hsla(200, 50%, 55%, 0.9);`"
                 class="label-c6 label-o3">
              {{ state.hd > 0 ? `+${state.hd.toFixed(2)}` : state.hd.toFixed(2) }}&deg;
            </div>
          </div>

        </div>
      </div>


    </List>
  </Element>
</template>

<style lang="scss">
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