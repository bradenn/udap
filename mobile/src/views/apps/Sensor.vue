<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import {onMounted, reactive, watchEffect} from "vue";
import core from "udap-ui/core";
import List from "udap-ui/components/List.vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";
import SensorView from "udap-ui/components/SensorView.vue";


const props = defineProps<{
  sensorData?: SensorData,
  refined?: boolean
  raw?: boolean,
  attribute?: string
}>()

const remote = core.remote();


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
  data: {} as SensorData,
  aligned: {} as SensorData,
  sensorName: "" as string,
  ready: false as boolean,
  history: {
    temp: null,
  } as SensorHistory
})

const router = core.router()
onMounted(() => {
  update()

})

watchEffect(() => {
  update()
  return props.sensorData
})
//
// function prepareData(input: number[]): number[] {
//   let avg = 0;
//   input.forEach(e => avg += e)
//   avg /= input.length
//   return input.map(x => (x - avg) / 100) as number[]
// }
//
// function defineMagnitude(unit: Sensor): number[] {
//   let out: number[] = [];
//   let input: number[] = []
//   for (let i = 0; i < unit.ax.length; i++) {
//     input.push(Math.sqrt(Math.pow(unit.ax[i], 2) + Math.pow(unit.ay[i], 2) + Math.pow(unit.az[i], 2)))
//   }
//   if (!unit.time) return out;
//   for (let i = 0; i < unit.time.length; i++) {
//     out.push(prepareData(input)[i])
//     // out.push(prepareData(unit.ay)[i])
//     // out.push(prepareData(unit.az)[i])
//   }
//
//
//   return out
// }

function cToF(input: number): number {


  return (input * (9 / 5) + 32)
}

function roundTo(input: number, factor: number): number {

  let a = Math.pow(10, factor)

  return Math.round(input * a) / a
}

//
function update() {

  // let id
  // if (props.sensorId) {
  //   id = props.sensorId;
  // } else {
  //   if (!router) return
  //   id = router?.currentRoute.value.params["sensorId"] as string;
  // }
  // if (!id) return

  // state.sensorName = id
  // state.attribute = remote.attributes.find(a => a.key.includes(id)) || {} as Attribute
  // if (!state.attribute) return;
  state.data = props.sensorData || {} as SensorData


  state.aligned.rssi = state.data.rssi
  state.aligned.lux = roundTo(state.data.lux, 2)
  state.aligned.rssi = roundTo(state.data.rssi, 0)
  if (props?.raw) {
    state.aligned.temp = roundTo(state.data.temp, 2) //  - 14,
    state.aligned.coreTemp = roundTo(state.data.coreTemp, 2)
  } else {
    state.aligned.temp = roundTo(cToF(state.data.temp), 2) //  - 14,
    state.aligned.coreTemp = roundTo(cToF(state.data.coreTemp), 2)
  }


  state.aligned.humidity = roundTo(state.data.humidity, 2) //  + 25
  state.aligned.ratio = roundTo(state.data.ratio * 100, 2)
  // state.data = JSON.parse(state.attribute.value) as SensorData
  //
  state.ready = true

}

</script>

<template>

  <div v-if="state.ready" class="d-flex gap-1 flex-column">
    <List v-if="!props.refined">
      <ElementHeader title="Lighting"></ElementHeader>
      <SensorView :attribute-id="props.sensorData?.attributeId || ''" :value="state.aligned.lux"
                  alt="Luminous flux intensity per square meter" icon="􀇯"
                  sensor-key="lux" title="Light Intensity"
                  unit="lux"></SensorView>


      <ElementHeader title="Environment"></ElementHeader>
      <SensorView :attribute-id="props.sensorData?.attributeId || ''" :value="state.aligned.temp"
                  alt="External ambient temperature" icon="􀇬"
                  sensor-key="temp" title="Temperature"
                  unit="&deg; F"></SensorView>
      <SensorView :attribute-id="props.sensorData?.attributeId || ''" :value="state.aligned.humidity"
                  alt="External ambient relative humidity" icon="􁃛"
                  sensor-key="humidity" title="Humidity"
                  unit="%"></SensorView>

      <ElementHeader title="Diagnostic"></ElementHeader>
      <SensorView :attribute-id="props.sensorData?.attributeId || ''" :value="state.aligned.coreTemp"
                  alt="Microcontroller's internal temperature" icon="􁊁"
                  sensor-key="coreTemp" title="Core Temperature"
                  unit="&deg; F"></SensorView>
      <SensorView :attribute-id="props.sensorData?.attributeId || ''" :value="state.aligned.rssi"
                  alt="Device's Wi-Fi signal reception quality" icon="􀖀"
                  sensor-key="rssi" title="Signal Strength"
                  unit="dBm"></SensorView>
      <SensorView :attribute-id="props.sensorData?.attributeId || ''" :value="state.aligned.ratio"
                  alt="Ratio of infrared and visible photodiodes" icon="􀇯"
                  sensor-key="ratio" title="Light Ratio"
                  unit="%"></SensorView>

    </List>
    <List v-else fixed-width row>
      <Element class="align-items-center d-flex justify-content-center gap-1" foreground style="height: 2rem">
        <span class="mono">{{ state.aligned.lux.toFixed(2) }}</span>
        <div class="label-o3 label-c5">lux</div>
      </Element>
      <Element class="align-items-center d-flex justify-content-center gap-1" foreground style="height: 2rem">
        <span class="mono">{{ state.aligned.coreTemp.toFixed(2) }}</span>
        <div class="label-o3 label-c5">&deg;F</div>
      </Element>
      <Element class="align-items-center d-flex justify-content-center gap-1" foreground style="height: 2rem">
        <span class="mono">{{ state.aligned.temp.toFixed(2) }}</span>
        <div class="label-o3 label-c5">&deg;F</div>
      </Element>
      <Element class="align-items-center d-flex justify-content-center gap-1" foreground style="height: 2rem">
        <div class="mono">{{ state.aligned.humidity.toFixed(2) }}</div>
        <div class="label-o3 label-c5">%</div>
      </Element>
    </List>
    <Element v-if="false" class="w-100" foreground>
      <!--      <List>-->

      <!--        <div style="width: 100%">-->
      <!--          <div class="d-flex gap-1 align-items-baseline">-->
      <!--            <div class="mono">-->
      <!--              {{ (Math.round(state.data.lux * 10) / 10).toFixed(1) }}-->
      <!--            </div>-->
      <!--            <div class="label-c5 label-o3 label-w500">lux</div>-->
      <!--          </div>-->
      <!--          <div class="d-flex  gap-1 align-items-baseline">-->
      <!--            <div class="mono">-->
      <!--              {{ (Math.round((state.data.temp * (9 / 5) + 18) * 10) / 10).toFixed(1) }}-->
      <!--            </div>-->
      <!--            <div class="label-c5 label-o3 label-w500">&deg; F</div>-->
      <!--          </div>-->
      <!--          <div class="d-flex gap-1 align-items-baseline">-->
      <!--            <div class="mono">-->
      <!--              {{ (Math.round((state.data.humidity + 23) * 10) / 10).toFixed(1) }}-->
      <!--            </div>-->
      <!--            <div class="label-c5 label-o3 label-w500">% RH</div>-->
      <!--          </div>-->
      <!--          <div class="d-flex gap-1 align-items-baseline">-->
      <!--            <div class="mono">-->
      <!--              {{ (Math.round((state.data.ratio) * 10000) / 100).toFixed(2) }}-->
      <!--            </div>-->
      <!--            <div class="label-c5 label-o3 label-w500">&deg; K</div>-->
      <!--          </div>-->
      <!--          <div class="d-flex gap-1 align-items-baseline">-->
      <!--            <div class="mono">-->
      <!--              {{ (Math.round((state.data.rssi) * 100) / 100).toFixed(0) }}-->
      <!--            </div>-->
      <!--            <div class="label-c5 label-o3 label-w500">dBm</div>-->
      <!--          </div>-->
      <!--          <div class="d-flex gap-1 align-items-baseline">-->
      <!--            <div class="mono">-->
      <!--              {{ (Math.round(((state.data.temp * 9 / 5) + 32) * 100) / 100).toFixed(2) }}-->
      <!--            </div>-->
      <!--            <div class="label-c5 label-o3 label-w500">&deg; F</div>-->
      <!--          </div>-->

      <!--        </div>-->

      <!--        &lt;!&ndash;        <ElementPair title="Illumination" icon="􀇯"&ndash;&gt;-->
      <!--        &lt;!&ndash;                     :value="(Math.round(state.data.lux * 100)/100).toFixed(2)+' lux'"></ElementPair>&ndash;&gt;-->
      <!--        &lt;!&ndash;        <ElementPair title="Temperature" icon="􁷉"&ndash;&gt;-->
      <!--        &lt;!&ndash;                     :value="(Math.round((state.data.temp * (9/5) + 18) * 100)/100).toFixed(2)+'&deg; F'"></ElementPair>&ndash;&gt;-->
      <!--        &lt;!&ndash;        <ElementPair title="Relative Humidity" icon="􁘰"&ndash;&gt;-->
      <!--        &lt;!&ndash;                     :value="(Math.round((state.data.humidity +23 )* 100)/100).toFixed(2) + '% RH'"></ElementPair>&ndash;&gt;-->
      <!--      </List>-->


    </Element>

  </div>
</template>

<style lang="scss" scoped>

</style>