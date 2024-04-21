<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import {onMounted, reactive, watchEffect} from "vue";
import {Attribute} from "udap-ui/types";
import core from "udap-ui/core";
import BarChart from "@/components/BarChart.vue";
import {ProcessedSensorData, Sensor} from "@/sensors/sensors";

const remote = core.remote();


const state = reactive({
  attributes: [] as Attribute[],
  sensors: [] as Sensor[],
  processedData: [] as ProcessedSensorData[],

})

onMounted(() => {
  update()
})

watchEffect(() => {
  update()
})

function prepareData(input: number[]): number[] {
  let avg = 0;
  input.forEach(e => avg += e)
  avg /= input.length
  return input.map(x => Math.abs((avg - x) / 1000))
}

function defineMagnitude(unit: Sensor): number[] {
  let out: number[] = [];
  if (!unit) return out;
  for (let i = 0; i < unit.time.length; i++) {
    out.push(Math.max(prepareData(unit.ax)[i], prepareData(unit.ay)[i], prepareData(unit.az)[i]))
  }


  return out
}

function update() {
  state.attributes = remote.attributes.filter(a => a.key.startsWith("/sensors/"))
  state.sensors = state.attributes.map(attribute => {
    return JSON.parse(attribute.value) as Sensor
  })

  state.processedData = state.sensors.map(data => {
    let processed = {
      core: 0,
      ambient: 0
    } as ProcessedSensorData

    if (data.core) {
      data.core.forEach(d => processed.core += d)
      processed.core /= data.core.length;
    }

    if (data.ambient) {
      data.ambient.forEach(d => processed.ambient += d)
      processed.ambient /= data.ambient.length;
    }

    processed.magnitude = defineMagnitude(data)

    processed.time = data.time

    return processed
  })
}

</script>

<template>
  <Element>
    <div v-for="sensor in state.processedData">
      <BarChart :data="sensor.magnitude" :scale="2" :time-data="sensor.time" unit="G"></BarChart>


    </div>
  </Element>
</template>

<style lang="scss" scoped>

</style>