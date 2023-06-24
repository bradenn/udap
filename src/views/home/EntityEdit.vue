<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {Attribute, Entity} from "udap-ui/types";
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
  online: false,
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

  state.online = state.attributes.map(a => (Date.now().valueOf() - new Date(a.lastUpdated).valueOf()) < (1000 * 60)).some(a => a)
}

function sendRequest(key: string, value: string) {
  let found = state.attributes.find(a => a.key === key)
  if (!found) {
    return
  }

  switch (key) {
    case "on":
      state.spectrum.on = value === "true"
      break
    case "cct":
      state.spectrum.cct = parseInt(value)
      break
    case "dim":
      state.spectrum.dim = parseInt(value)
      break
    case "hue":
      state.spectrum.hue = parseInt(value)
      break
  }
  found.request = value
  found.entity = state.entity.id
  found.key = key
  attributeService.request(found).then(e => {
    console.log(e)
  }).catch(err => {
    console.log(err)
  })
}

function setDim(value: number) {
  sendRequest("dim", `${Math.max(value, 1)}`)
}

</script>

<template>
  <div class="d-flex gap-2 flex-column">
    <div class="d-flex gap-1">
      <router-link class="sf-icon element m-0 label-w600" style="width: 3.5rem; color: rgba(255,255,255,0.2)"
                   to="/home/dashboard">􀯶
      </router-link>
      <div :class="`${state.spectrum.on?'on':'off'}`"
           class="d-flex align-items-center justify-content-start gap-4 header flex-grow-1"
           @click="() => sendRequest('on', state.spectrum.on?'false':'true')">
        <div class=" label-c2 sf-icon" style="padding-left: 1.5rem !important;">{{ state.entity.icon }}</div>
        <div class="d-flex justify-content-between w-100 align-items-center">
          <div class="d-flex align-items-center gap-2">
            <div class="label-c5 label-w700 label-o5 ">{{
                state.entity.alias ? state.entity.alias : state.entity.name
              }}
            </div>
          </div>

          <div class="px-3 label-w600 entity-status">
            {{ state.spectrum.on ? "ON" : "OFF" }}
          </div>
        </div>
      </div>
    </div>
    <div v-if="!state.online" class="element px-3 d-flex gap-2">
      <div class="sf-icon text-warning">􀇿</div>
      <div class="label-o5 label-w600">This device is offline.</div>
    </div>
    <div :class="!state.online?'block-usage':''" class="d-flex flex-column gap-1">
      <div v-if="state.attributes.map(a => a.key).includes('dim')">
        <div class="d-flex justify-content-between px-1">
          <div class="label-c4 label-w600 label-o6">Brightness</div>
          <div class="label-c4 label-w600 label-o5">{{ state.spectrum.dim }}%</div>
        </div>
        <div>
          <Slider :change="(v) => sendRequest('dim', `${v}`)" :max="100" :min="1"
                  :step="2" :value="state.spectrum.dim"
                  bg="dim"></Slider>
        </div>
      </div>
      <div v-if="state.attributes.map(a => a.key).includes('cct')">
        <div class="d-flex justify-content-between px-1">
          <div class="label-c4 label-w600 label-o6">Color Temperature</div>
          <div class="label-c4 label-w600 label-o5">{{ state.spectrum.cct }}&deg; K</div>
        </div>
        <Slider :change="(v) => sendRequest('cct', `${v}`)" :max="6000"
                :min="2000"
                :step="1" :value="state.spectrum.cct"
                bg="cct"></Slider>
      </div>
      <div v-if="state.attributes.map(a => a.key).includes('hue')">
        <div class="d-flex justify-content-between px-1">
          <div class="label-c4 label-w600 label-o6">Color Hue</div>
          <div class="label-c4 label-w600 label-o5">{{ state.spectrum.hue }}&deg;</div>
        </div>
        <Slider :change="(v) => sendRequest('hue', `${v}`)" :max="360" :min="1"
                :step="1"
                :value="state.spectrum.hue" bg="hue"></Slider>
      </div>
    </div>
  </div>
</template>

<style scoped>

.block-usage {
  filter: blur(2px) brightness(50%);
}

.header {

  padding: 0.25rem 0.25rem;

  backdrop-filter: blur(40px);
  box-shadow: inset 0 0 1px 1.5px rgba(37, 37, 37, 0.6), 0 0 3px 1px rgba(22, 22, 22, 0.6);
  /* Note: backdrop-filter has minimal browser support */
  border-radius: 11.5px;
  -webkit-backdrop-filter: blur(40px) !important;
}
</style>