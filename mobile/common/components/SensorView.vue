<!-- Copyright (c) 2024 Braden Nicholson -->

<script lang="ts" setup>

import {onMounted, reactive} from "vue";
import {v4 as uuidv4} from "uuid";
import LineChart from "./LineChart.vue";
import ElementPair from "./ElementPair.vue";
import attributeService2 from "udap-ui/services/attributeService";
import {SensorEntry} from "@/views/apps/Sensor.vue";
import List from "udap-ui/components/List.vue";
import {SelectOption} from "udap-ui/components/Select.vue";
import ElementSelect, {Item} from "udap-ui/components/ElementSelect.vue";

const props = defineProps<{
  value: number,
  unit: string,
  title: string,
  alt: string,
  icon: string,
  sensorKey: string,
  attributeId: string,
}>()

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

const state = reactive({
  canvas: {} as HTMLCanvasElement,
  ctx: {} as CanvasRenderingContext2D,
  animation: 0,
  toggled: false,
  uuid: uuidv4(),
  ready: false,
  history: {} as SensorEntry,
  metadata: {
    mode: "",
    window: 0,
    from: 0,
    to: 0,
    points: 0,
  }
})

onMounted(() => {
  state.metadata = {
    mode: "AVG",
    window: 60,
    to: Date.now(),
    from: Date.now() - (1000 * 60 * 60 * 4),
    points: 0
  }

})

function pullData() {
  fetchKey(props.sensorKey, state.metadata.from, state.metadata.to, state.metadata.window, state.metadata.mode)
}

function fetchKey(key: string, from: number, to: number, window: number, mode: string) {

  queryHistory(key, from, to, window, mode).then(res => {
    state.history = res as SensorEntry
    state.ready = true
    state.metadata = {
      mode: mode,
      window: window,
      to: to,
      from: from,
      points: state.history.values.length
    }
    console.log(res)
  }).catch(err => {

  })

}

function queryHistory(key: string, from: number, to: number, window: number, mode: string): Promise<SensorEntry> {

  return attributeService2.summary(props.attributeId, key, from, to, 1000 * window, mode)
      .then(r => {
        let entry: SensorEntry = {
          timestamps: [],
          values: []
        } as SensorEntry
        let prs = r.data
        let kys = Object.keys(prs)

        for (let i = 0; i < kys.length; i++) {

          entry.timestamps.push(Number(kys[i]))
          entry.values.push(Number(prs[kys[i]]))
        }

        return entry

      }).catch(r => {
        return {} as SensorEntry;
      })
}

function toggle() {
  if (!state.ready && !state.toggled) {
    pullData()
  }
  state.toggled = !state.toggled
}

function toggleSample() {
  if (state.metadata.window == 60) {
    state.metadata.window = 120;
  } else if (state.metadata.window == 120) {
    state.metadata.window = 240
  } else if (state.metadata.window == 240) {
    state.metadata.window = 480
  } else {
    state.metadata.window = 60;
  }
  pullData()
}

const periods = [15 * 60 * 1000, 30 * 60 * 1000, 60 * 60 * 1000, 2 * 60 * 60 * 1000, 4 * 60 * 60 * 1000, 8 * 60 * 60 * 1000, 12 * 60 * 60 * 1000]

function togglePeriod() {
  if (state.metadata.window == 60) {
    state.metadata.window = 120;
  } else if (state.metadata.window == 120) {
    state.metadata.window = 240
  } else if (state.metadata.window == 240) {
    state.metadata.window = 480
  } else {
    state.metadata.window = 60;
  }
  pullData()
}

function updatePeriod(value: string) {
  state.metadata.from = Date.now().valueOf() - 1000 * parseInt(value)
  state.metadata.to = Date.now().valueOf()
  pullData()
}

function updateSample(value: string) {
  state.metadata.window = Math.min(Math.max(parseInt(value), 1), 60 * 60)
  pullData()
}

</script>

<template>
  <ElementPair :alt="props.alt"
               :cb="toggle"
               :icon="props.icon"
               :title="props.title">
    {{ props.value }} {{ props.unit }}
    <template v-slot:body @touchstart.stop>
      <List v-if="state.toggled">
        <LineChart v-if=state.ready :data="state.history.values" :time="state.history.timestamps"
                   class="pt-2"></LineChart>
        <List>
          <div class="d-flex gap-3">
            <div class="d-flex gap-2">
              <div>
                <div class="label-c5 label-o2">Mode</div>
                <div class="label-c5 label-o2">Sample</div>
              </div>
              <div>
                <div class="label-c5 label-o3">{{ state.metadata.mode }}</div>
                <div class="label-c5 label-o3">{{ state.metadata.window }}s</div>
              </div>
            </div>

            <div class="d-flex gap-2">
              <div>
                <div class="label-c5 label-o2">Period</div>
                <div class="label-c5 label-o2">Points</div>
              </div>
              <div>
                <div class="label-c5 label-o3">{{ (state.metadata.to - state.metadata.from) / 1000 / 60 / 60 }}h</div>
                <div class="label-c5 label-o3">{{ state.metadata.points }}</div>
              </div>
            </div>
          </div>

          <div class="label-c4 label-o4 label-w500">Period</div>
          <ElementSelect
              :change="(v) => updatePeriod(v)"
              :options="[{key: '1 h', value: '3600'}, {key: '4 h', value: '14400'},{key: '8 h', value: '28800'}, {key: '12 h', value: '43200'}, {key: '24 h', value: '86400'}]"
              :value="`${(state.metadata.to-state.metadata.from)/1000}`"></ElementSelect>
          <div class="label-c4 label-o4 label-w500">Sample</div>
          <ElementSelect
              :change="(v) => updateSample(v)"
              :options="[{key: '1s', value: '1'}, {key: '10s', value: '10'},{key: '1m', value: '60'}, {key: '5m', value: '300'}]"
              :value="`${state.metadata.window}`"></ElementSelect>
        </List>
      </List>
    </template>

  </ElementPair>

</template>

<style lang="scss" scoped>

</style>