<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";

import {onBeforeUnmount, onMounted, onUpdated, reactive, watchEffect} from "vue";
import type {Attribute, Entity, Trigger} from "udap-ui/types";
import core from "udap-ui/core";
import List from "udap-ui/components/List.vue";


import SensorDOM, {SensorData} from "@/views/apps/Sensor.vue";
import triggerService from "udap-ui/services/triggerService";
import attributeService from "@/services/attributeService";

const remote = core.remote();

const router = core.router()

const props = defineProps<{
  entityId: string,
  alignment?: SensorData
}>()

interface UHIDv2 {
  name: string,
  version: string,
  rebooting: boolean,
  updating: boolean,
  progress: number,
  updateStatus: string,
}


const state = reactive({
  entityId: "",
  live: true,
  entity: {} as Entity,
  alignmentAttribute: {} as Attribute,
  attributes: [] as Attribute[],
  triggers: [] as Trigger[],
  sensor: {} as SensorData,
  aligned: {} as SensorData,
  alignedTo: {} as SensorData,
  alignment: {} as SensorData,
  metadata: {} as UHIDv2,
  lastUpdate: 0,
  deviceState: [] as Warning[]

})

onMounted(() => {
  update()
})

onBeforeUnmount(() => {
  state.live = false
})

watchEffect(() => {
  if (state.live) {
    update()
    calculateDeltas();
  }
  return remote
})

onUpdated(() => {
  calculateDeltas();
})

watchEffect(() => {
  calculateDeltas()
  return props.alignment
})

function calculateDeltas() {
  if (props.alignment) {

    state.alignedTo = deltaData(state.sensor, props.alignment)

  }
  applyAlignment()
}

interface Warning {
  message: string,
  criteria: string,
  level: number
}

function createDeviceState(): Warning[] {

  let errors = [] as Warning[]

  if (state.sensor.coreTemp >= 80) {
    errors.push({
      message: "Microcontroller Overheating",
      criteria: "Core temp exceeds 50&deg; C",
      level: 6,
    })
  } else if (state.sensor.coreTemp >= 60) {
    errors.push({
      message: "Microcontroller Critical Overheating",
      criteria: "Core temp exceeds 60&deg; C",
      level: 10,
    })
  } else if (state.sensor.coreTemp >= 50) {
    errors.push({
      message: "Microcontroller Overheating",
      criteria: "Core temp exceeds 50&deg; C",
      level: 6,
    })
  }


  if (state.sensor.temp >= 60) {
    errors.push({
      message: "Critical Ambient Heat",
      criteria: "External ambient temperature exceeds 60&deg; C",
      level: 10,
    })
  } else if (state.sensor.temp >= 40) {
    errors.push({
      message: "Excessive Ambient Heat",
      criteria: "External ambient temperature exceeds 40&deg; C",
      level: 6,
    })
  }

  if (state.sensor.humidity >= 90) {
    errors.push({
      message: "Critical Ambient Humidity",
      criteria: "External ambient humidity exceeds 90%",
      level: 10,
    })
  } else if (state.sensor.humidity >= 60) {
    errors.push({
      message: "Excessive Ambient Humidity",
      criteria: "External ambient humidity exceeds 60%",
      level: 6,
    })
  }

  if (state.sensor.rssi <= -90) {
    errors.push({
      message: "Critical Signal Strength",
      criteria: "Wi-Fi RSSI falls below -90dBm",
      level: 4,
    })
  } else if (state.sensor.rssi <= -75) {
    errors.push({
      message: "Poor Signal Strength",
      criteria: "Wi-Fi RSSI falls below -75dBm",
      level: 2,
    })
  }

  if (state.sensor.lux >= 600) {
    errors.push({
      message: "Light Sensor Overload",
      criteria: "Lux measurement exceeds 600 Lux",
      level: 1,
    })
  }

  if (state.sensor.ratio > 0.9 || state.sensor.ratio < 0.2) {
    errors.push({
      message: "Light Sensor Out of Range",
      criteria: "Light sensor ration not between 0.2 and 0.6",
      level: 1,
    })
  }

  return errors
}

function update() {
  state.entityId = props.entityId
  state.entity = remote.entities.find(e => e.type.includes("uhidv2") && e.id == state.entityId)
  if (!state.entity) return
  state.attributes = remote.attributes.filter(a => a.entity == state.entityId)
  let sensor = state.attributes.find(a => a.key == "sensor")
  if (!sensor) return

  state.sensor = JSON.parse(sensor?.value || "{}") as SensorData || {} as SensorData

  state.attributes = remote.attributes.filter(a => a.entity == state.entityId)
  let alignment = state.attributes.find(a => a.key == "alignment")
  if (!alignment) {
    alignment = sensor
  }
  let protoAlignment = JSON.parse(alignment?.value || "{}") as SensorData

  state.alignment = protoAlignment
  if (Object.keys(protoAlignment).length != 0) {
    applyAlignment()

  } else {
    reset()
    applyAlignment()

  }

  state.alignedTo = protoAlignment
  state.alignmentAttribute = alignment
  let metadata = state.attributes.find(a => a.key == "metadata")
  if (!metadata) return
  state.metadata = JSON.parse(metadata?.value || "{}") as UHIDv2 || {} as UHIDv2
  let now = new Date().valueOf()
  if (state.sensor) {
    state.lastUpdate = Math.max(new Date(sensor?.updated | 0).valueOf(), new Date(metadata?.updated | 0).valueOf())
  }
  state.triggers = remote.triggers.filter(t => t.name.includes(state.entity.name))

  state.deviceState = createDeviceState()

}

function reset() {
  let alignment = {} as SensorData
  let keys = Object.keys(state.sensor);
  for (let i = 0; i < keys.length; i++) {
    let name = keys[i]
    alignment[name] = 0
  }
  state.alignment = alignment
  updateAlignment()

}

function applyAlignment() {

  state.aligned = alignData(state.sensor, state.alignment)

}

function alignData(data: SensorData, alignment: SensorData): SensorData {
  let aligned: SensorData = {} as SensorData
  Object.keys(data).forEach((key: string) => {

    let offset = 0
    if (alignment) {
      aligned[key] = data[key] + alignment[key]
    } else {
      aligned[key] = data[key]
    }

  })

  return aligned
}

function deltaData(data: SensorData, alignment: SensorData): SensorData {
  let aligned: SensorData = {} as SensorData
  Object.keys(data).forEach((key: string) => {
    if (key == "lux" || key == "coreTemp") {
      if (key == "coreTemp") {
        aligned[key] = -(5 / 9) * 32
      } else {
        aligned[key] = 0
      }

      return
    }

    aligned[key] = (data[key] - alignment[key])
    if (key == "temp") {
      // aligned[key] -= 32
      // aligned[key] *= 5 / 9
    }
  })

  return aligned
}

function goBack() {
  router.go(-1);
}

function callTrigger(trigger: Trigger) {
  triggerService.invoke(trigger).then(() => {

  }).catch(() => {

  })
}

function parseName(input: string): string {
  let chunks = input.split("-")
  switch (chunks[chunks.length - 1]) {
    case "dial":
      return "fn1"
    case "on":
      return "On"
    case "off":
      return "Off"
  }

  return "unknown"

}

function updateAlignment() {
  let alignment = state.alignment
  alignment.coreTemp = 0
  alignment.temp = -state.alignedTo.temp
  alignment.humidity = -state.alignedTo.humidity
  state.alignmentAttribute.request = JSON.stringify(alignment)
  // state.alignment = alignment
  attributeService.request(state.alignmentAttribute).then(() => {
  }).catch(() => {
  })
}


</script>

<template>

  <Element class=" d-flex flex-column gap-1">
    <List class="flex-grow-0 scrollable-fixed">
      <div class="d-flex justify-content-between align-items-center">
        <div class="label-c4 label-o5 label-w600 px-1 mono" style="">
          {{ ((state.entity.alias) ? state.entity.alias : state.entity.name) }}
        </div>
        <div class="label-c6 label-o3 label-w600 mono px-1 lh-1">
          v{{ state.metadata.version }}
        </div>
      </div>
      <List class="scrollable-fixed">
        <SensorDOM :sensor-data="state.aligned" refined></SensorDOM>
        <!--        <SensorDOM raw refined :sensor-data="state.alignment"></SensorDOM>-->
        <!--        <SensorDOM raw refined :sensor-data="state.alignedTo"></SensorDOM>-->
        <!--        -->
      </List>
    </List>
  </Element>

</template>

<style lang="scss" scoped>

</style>