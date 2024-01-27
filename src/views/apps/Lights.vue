<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import type {Lights} from "udap-ui/composables/lights"
import {useLights} from "udap-ui/composables/lights"
import Element from "udap-ui/components/Element.vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";
import List from "udap-ui/components/List.vue";
import {onMounted, onUnmounted, reactive, watchEffect} from "vue";
import core from "@/core";
import {PreferencesRemote} from "udap-ui/persistent";
import Slider from "@/components/Slider.vue";
import attributeService from "@/services/attributeService";
import {Attribute} from "@/types";


const lights = useLights() as Lights

const state = reactive({
  selected: [] as string[],
  spectrum: {
    on: false,
    dim: 0,
    cct: 0,
    hue: 0,
  },
  discoTimer: 0,
  discoLength: 20,
  discoDuration: 0,
  discoRemains: 0,
  disco: false,
  online: false,
})

const preferences = core.preferences() as PreferencesRemote

onMounted(() => {
  updateAttribute()
})

watchEffect(() => {
  // updateAttribute()
  return lights.attributes
})

function select(id: string) {
  if (state.selected.includes(id)) {
    state.selected = state.selected.filter(s => s != id)
  } else {
    state.selected.push(id)
  }
  updateAttribute()
}

function updateAttribute() {

  let attributes = lights.attributes.filter(a => state.selected.includes(a.entity))
  let found = attributes.find(a => a.key === 'on')
  state.spectrum.on = found?.value === "true"

  found = attributes.find(a => a.key === 'dim')
  state.spectrum.dim = parseInt(found?.value || "0")

  found = attributes.find(a => a.key === 'cct')
  state.spectrum.cct = parseInt(found?.value || "0")

  found = attributes.find(a => a.key === 'hue')
  state.spectrum.hue = parseInt(found?.value || "0")

  state.online = attributes.map(a => (Date.now().valueOf() - new Date(a.lastUpdated).valueOf()) < (1000 * 60)).some(a => a)
}

function sendRandomHue() {
  let found: Attribute = lights.attributes.find(a => a.key === "hue") || {} as Attribute
  if (!found) {
    return
  }
  found.key = "hue"
  state.selected.forEach(s => {
    found.entity = s
    found.request = `${Math.round(Math.random() * 360)}`
    console.log(found)
    attributeService.request(found).then(e => {
      console.log(e)
    }).catch(err => {
      console.log(err)
    })
  })
}

function sendRequest(key: string, value: string) {
  let found = lights.attributes.find(a => a.key === key)
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
  found.key = key
  state.selected.forEach(s => {
    if (!found) return;
    found.entity = s
    if (!found) return
    attributeService.request(found).then(e => {
      console.log(e)
    }).catch(err => {
      console.log(err)
    })
  })


}

onUnmounted(() => {
  haltDisco()
})

function haltDisco() {
  clearInterval(state.discoTimer)
  state.discoDuration = 0
  state.discoTimer = 0;
  state.discoRemains = 0;
  state.disco = false
}

function toggleDisco() {
  if (!state.disco) {
    sendRequest("on", "true")
    //@ts-ignore
    state.discoTimer = setInterval(() => {
      if (state.discoDuration >= state.discoLength) {
        haltDisco()

        sendRequest("cct", "5000")
        sendRequest("dim", "25")
        return
      }
      if (state.discoDuration == 0) {
        sendRequest("dim", "100")
      }
      sendRandomHue()
      state.discoDuration++;
      state.discoRemains = (state.discoDuration / state.discoLength)
    }, 500)
    state.disco = true
  } else {
    haltDisco()
  }

}
</script>

<template>

  <List>

    <Element>
      <List scroll-y style="max-height: 30vh">
        <div v-for="zone in lights.zones" v-if="lights.loaded">
          <ElementHeader :title="zone.name"></ElementHeader>
          <List>
            <Element v-for="entity in zone.entities" :accent="state.selected.includes(entity.id)"
                     :cb="() => select(entity.id)"
                     class="d-flex align-items-center justify-content-between" foreground
                     mutable
                     style="height: 3.25rem">
              <div class="d-flex gap-2">
                <div class="sf-icon">{{ entity.icon }}</div>
                <div>{{ entity.alias ? entity.alias : entity.name }}</div>
              </div>
              <div>
                <div v-if="!state.selected.includes(entity.id)" class="sf-icon px-3 label-o1">􀀀</div>
                <div v-else :style="`color: ${preferences.accent};`" class="sf-icon px-3">􀷙</div>
              </div>
            </Element>
          </List>
        </div>
      </List>
    </Element>
    <Element :class="!state.online?'block-usage':''">
      <List scroll-y style="max-height: 40vh">
        <Element v-if="lights.attributes.map(a => a.key).includes('on')"
                 :cb="() => sendRequest('on', state.spectrum.on?'false':'true')"
                 class="d-flex align-items-center justify-content-between px-3 py-3"
                 foreground style="height: 3.25rem">
          <div class="d-flex gap-2">

            <div class="label-o4 label-w600">Power</div>
          </div>
          <div class="label-o4 label-w600 d-flex gap-2">
            <div :class="`${state.spectrum.on?'label-o4':'label-o2'}`">ON</div>
            <div :class="`${state.spectrum.on?'label-o2':'label-o4'}`">OFF</div>
          </div>
        </Element>
        <Element v-if="lights.attributes.map(a => a.key).includes('dim')" foreground>
          <ElementHeader :alt="`${state.spectrum.dim}%`" title="Brightness"></ElementHeader>
          <Slider :change="(v) => sendRequest('dim', `${v}`)" :max="100" :min="1"
                  :step="2" :value="state.spectrum.dim"
                  bg="dim"></Slider>
        </Element>
        <Element v-if="lights.attributes.map(a => a.key).includes('cct')" foreground>
          <ElementHeader :alt="`${state.spectrum.cct}&deg; K`" title="Color Temperature"></ElementHeader>
          <Slider :change="(v) => sendRequest('cct', `${v}`)" :max="6000"
                  :min="2000"
                  :step="1" :value="state.spectrum.cct"
                  bg="cct"></Slider>
        </Element>
        <Element v-if="lights.attributes.map(a => a.key).includes('hue')" foreground>
          <ElementHeader :alt="`${state.spectrum.hue}&deg;`" title="Color Hue"></ElementHeader>
          <Slider :change="(v) => sendRequest('hue', `${v}`)" :max="360" :min="1"
                  :step="1"
                  :value="state.spectrum.hue" bg="hue"></Slider>
        </Element>
      </List>
    </Element>
  </List>
</template>

<style>

</style>