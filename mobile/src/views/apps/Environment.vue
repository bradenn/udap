<!-- Copyright (c) 2024 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import core from "udap-ui/core";
import {onMounted, reactive} from "vue";
import type {Attribute, Entity} from "udap-ui/types.ts";
import {Remote} from "udap-ui/remote";
import List from "udap-ui/components/List.vue";
import type {Trace, TraceRequest} from "udap-ui/services/traceService";
import traceService from "udap-ui/services/traceService";

import MultiChart from "udap-ui/components/trace/MultiChart.vue";
import ElementSelect from "udap-ui/components/ElementSelect.vue";
import RadialChart from "udap-ui/components/trace/RadialChart.vue";
import Floorplan from "udap-ui/components/floorplan/Floorplan.vue";
import Timeline from "udap-ui/components/floorplan/Timeline.vue";


export interface SensorData {
  attributeId: string,
  lux: number,
  temp: number,
  humidity: number,
  ratio: number,
  rssi: number,
  coreTemp: number,
}

export interface SensorEntry {
  values: number[],
  timestamps: number[],
}

export interface SensorHistory {
  lux: SensorEntry,
  temp: SensorEntry,
  humidity: SensorEntry,
  ratio: SensorEntry,
  rssi: SensorEntry,
  coreTemp: SensorEntry,
}

const remote = core.remote() as Remote

const state = reactive({
  traces: [] as Trace[],
  sources: [] as ED[],
  duration: 60 * 60 * 4,
  window: 60,
  historical: [] as Trace[]
})

onMounted(() => {
  pullData()
})

function pullData() {
  queryHistory("", "", Date.now() - state.duration * 1000, Date.now(), 1000 * state.window, "AVG")
  locateSensorData()
  queryHistorical(Date.now() - 86400000 * 14, Date.now(), 0, "AVG")
}

let sensors: string[] = ["temp", "coreTemp", "lux", "humidity", "ratio"]

export interface ED {
  entity: Entity,
  attributes: Attribute[],

}

function locateSensorData() {
  let entities = remote.entities.filter(e => e.type == "uhidv2") as Entity[]
  if (!entities) return;

  state.sources = []
  for (let entity: Entity of entities) {

    state.sources.push({
      entity: entity,
      attributes: remote.attributes.filter(a => a.entity == entity.id)
    })
  }

  for (let source of state.sources) {

  }


}

function queryHistory(id: string, key: string, from: number, to: number, window: number, mode: string) {

  let req: TraceRequest = {
    to: to,
    from: from,
    window: window,
    mode: mode,
    labels: ["class=temp", "model=uhidv2-sensor"]
  } as TraceRequest

  traceService.trace(req).then(res => {
    let trc = res.traces
    state.traces = []
    for (let trace of trc) {
      if (trace.labels.hasOwnProperty("serial") && trace.labels["serial"].includes(" ")) {
        state.traces.push(trace)
      }

    }

  }).catch(err => {
    console.error(err)
  })

}

function queryHistorical(from: number, to: number, window: number, mode: string) {


  let req: TraceRequest = {
    to: to,
    from: from,
    window: 86400000 / 8,
    mode: mode,
    labels: ["class=temp", "model=uhidv2-sensor", "serial!=Prototype"]
  } as TraceRequest

  traceService.trace(req).then(res => {
    let trc = res.traces
    state.historical = []
    for (let trace of trc) {
      if (trace.labels.hasOwnProperty("serial") && trace.labels["serial"].includes(" ")) {
        state.historical.push(trace)
      }

    }

  }).catch(err => {
    console.error(err)
  })

}

function setDuration(dur: string) {
  state.duration = parseInt(dur)
  pullData()
}

function updateWindow(window: string) {
  state.window = parseInt(window)
  pullData()
}


function cToF(degC: number): number {
  return (degC * (9 / 5)) + 32
}
</script>

<template>
  <div class="scrollable-fixed">
    <List scroll-y>
      <Timeline :traces="state.historical"></Timeline>
      <Element class="d-flex flex-column gap-1">
        <MultiChart :traces="state.traces"
                    :transforms="[(c: number) => (c-6.5),(c: number) => ((c * (9.0 / 5.0)) + 32)]"
                    style="">
        </MultiChart>
        <Floorplan :traces="state.traces"
                   :transforms="[(c: number) => (c-6.5),(c: number) => ((c * (9.0 / 5.0)) + 32)]"></Floorplan>
        <RadialChart v-if=false :traces="state.traces"
                     :transforms="[(c: number) => (c-6.5),(c: number) => ((c * (9.0 / 5.0)) + 32)]"
                     style="height: 40rem">
        </RadialChart>

        <div class="label-c4 label-o4 label-w500">Period</div>
        <ElementSelect
            :change="(v) => setDuration(v)"
            :options="[{key: '1 h', value: '3600'}, {key: '4 h', value: '14400'},{key: '8 h', value: '28800'}, {key: '12 h', value: '43200'}, {key: '24 h', value: '86400'}, {key: '2 d', value: '172800'}, {key: '4 d', value: '345600'}, {key: '1 w', value: '604800'}]"
            :value="`${state.duration}`"></ElementSelect>
        <div class="label-c4 label-o4 label-w500">Sample</div>
        <ElementSelect
            :change="(v) => updateWindow(v)"
            :options="[{key: '1s', value: '1'}, {key: '1m', value: '60'}, {key: '5m', value: '300'}, {key: '15m', value: '900'}, {key: '30m', value: '1800'}, {key: '1h', value: '3600'}]"
            :value="`${state.window}`"></ElementSelect>
      </Element>
    </List>
  </div>
</template>

<style lang="scss" scoped>
.dashboard-grid {
  display: grid;
  height: 100%;
  grid-gap: 0.375rem;
  grid-template-columns: repeat(12, minmax(1rem, 1fr));

}
</style>