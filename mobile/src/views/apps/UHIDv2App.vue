<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import ElementLink from "udap-ui/components/ElementLink.vue";
import {onMounted, reactive, watchEffect} from "vue";
import {Attribute, Entity} from "udap-ui/types";
import core from "udap-ui/core";
import List from "udap-ui/components/List.vue";

import type {SensorData} from "@/views/apps/Sensor.vue";
import SensorDOM from "@/views/apps/Sensor.vue";
import UHIDv2Preview from "@/views/apps/UHIDv2Preview.vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";

const remote = core.remote();

interface UHIDv2 {
  name: string,
  version: string,
  rebooting: boolean,
  updating: boolean,
  progress: number,
  updateStatus: string,
}


const state = reactive({
  baseline: {} as SensorData,
  thermostat: {} as Entity,
  entities: [] as Entity[],
  attributes: [] as Attribute[],
  sensors: [] as Attribute[],
  alignments: [] as Attribute[],
  devices: [] as Entity[],
  hubs: [] as Entity[],


})

onMounted(() => {
  update()
})

watchEffect(() => {
  update()
  return remote.entities
})


function alignSensors() {
  for (let alignment in state.sensors) {

  }
}

function update() {
  state.entities = remote.entities.filter(e => (e.type.includes("uhidv2") || e.type.includes("vradarv4")))
  let entityIds: string[] = state.entities.map(e => e.id)
  state.attributes = remote.attributes.filter(a => entityIds.includes(a.entity))
  state.sensors = state.attributes.filter(a => a.key == "sensor")
  state.alignments = state.attributes.filter(a => a.key == "alignment")

  state.devices = []
  state.hubs = []
  for (let entity of state.entities) {
    if (entity.type == "uhidv2") {
      state.devices.push(entity)
    } else if (entity.type == "vradarv4") {
      state.hubs.push(entity)
    }
  }

  let thermostat = remote.entities.find(e => e.name == "thermostat")
  if (!thermostat) return;
  state.thermostat = thermostat
  let temp = remote.attributes.find(a => a.key === "thermostat")
  if (temp) {
    let val = JSON.parse(temp.value) as {
      connected: boolean
      temperature: number
      humidity: number
      ecoCool: number
      ecoHeat: number
      mode: string
    }
    state.baseline.temp = (val.temperature - 32) * (5 / 9)
    state.baseline.humidity = val.humidity
    state.baseline.lux = 0
    state.baseline.coreTemp = -32 * (5 / 9)
  }
}


</script>

<template>
  <List scroll-y>
    <Element>
      <List>

        <ElementHeader title="Sensors"></ElementHeader>
        <ElementLink v-for="sensor in state.devices"
                     :alt="sensor.type"
                     :title="sensor.alias?sensor.alias:sensor.name"
                     :to="`/apps/uhidv2/${sensor.id}`" icon="􁔊"></ElementLink>
        <ElementHeader title="Hubs"></ElementHeader>
        <ElementLink v-for="sensor in state.hubs"
                     :title="sensor.alias?sensor.alias:sensor.name"
                     :to="`/apps/uhidv2/${sensor.id}`"
                     :alt="sensor.type" icon="􁅏">
        </ElementLink>
      </List>
    </Element>

    <List class="scrollable-fixed">
      <Element class="d-flex flex-column gap-1">
        <div class="d-flex justify-content-between align-items-center">
          <div class="label-c4 label-o5 label-w600 px-1 mono" style="">
            Thermostat
          </div>
          <div class="label-c6 label-o3 label-w600 mono px-1 lh-1">

          </div>
        </div>
        <SensorDOM :sensor-data="state.baseline" refined></SensorDOM>
      </Element>
      <List>
        <div v-for="item in state.sensors">
          <UHIDv2Preview :alignment="state.baseline" :entity-id="item.entity"></UHIDv2Preview>
        </div>
      </List>
    </List>
    <!--      <PanTiltVisualizer v-for="sensor in state.sensors" :time="sensor.time" :x="sensor.ax" :y="sensor.ax"-->
    <!--                         :z="sensor.az"-->
    <!--                         style="width: 100%"></PanTiltVisualizer>-->

  </List>
</template>

<style lang="scss" scoped>

</style>