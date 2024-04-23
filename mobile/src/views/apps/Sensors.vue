<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import ElementLink from "udap-ui/components/ElementLink.vue";
import {onMounted, reactive, watchEffect} from "vue";
import {Attribute} from "udap-ui/types";
import core from "udap-ui/core";
import List from "udap-ui/components/List.vue";


import type {SensorData} from "@/views/apps/Sensor.vue";
import SensorDOM from "@/views/apps/Sensor.vue";

const remote = core.remote();

interface Sensor {
  lux: number,
  temp: number,
  humidity: number,
}

const state = reactive({

  attributes: [] as Attribute[],
  sensors: [] as SensorData[],


})

onMounted(() => {
  update()
})

watchEffect(() => {
  update()
  return remote.attributes
})

const nameMap = [{key: "df08", name: "Living Room TV"}, {key: "df16", name: "Kitchen"}]

function update() {
  state.attributes = remote.attributes.filter(a => a.type == "uhidv2-sensor")
  state.sensors = state.attributes.map(attribute => {
    return JSON.parse(attribute?.value || "{}") as SensorData || {} as SensorData
  })


}

</script>

<template>
  <List>
    <Element>
      <List>
        <ElementLink v-for="sensor in state.attributes"
                     :title="nameMap.find(n => sensor.key.includes(n.key))?.name || sensor.key"
                     :to="`/apps${sensor.key}`"
                     icon="􁔊"></ElementLink>
      </List>
    </Element>
    <Element>
      <List>

        <SensorDOM v-for="sensor in state.attributes" :sensor-id="sensor.key"></SensorDOM>


      </List>
      <!--      <PanTiltVisualizer v-for="sensor in state.sensors" :time="sensor.time" :x="sensor.ax" :y="sensor.ax"-->
      <!--                         :z="sensor.az"-->
      <!--                         style="width: 100%"></PanTiltVisualizer>-->
    </Element>
  </List>
</template>

<style lang="scss" scoped>

</style>