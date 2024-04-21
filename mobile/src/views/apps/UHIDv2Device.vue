<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";

import {onBeforeUnmount, onMounted, reactive, watchEffect} from "vue";
import type {Attribute, Entity, Trigger} from "udap-ui/types";
import core from "udap-ui/core";
import List from "udap-ui/components/List.vue";
import ElementPair from "udap-ui/components/ElementPair.vue";
import ElementTextBox from "udap-ui/components/ElementTextBox.vue";


import type {SensorData} from "@/views/apps/Sensor.vue";
import SensorDOM from "@/views/apps/Sensor.vue";
import entityService from "udap-ui/services/entityService";
import attributeService from "@/services/attributeService";
import ElementHeader from "udap-ui/components/ElementHeader.vue";
import Slider from "udap-ui/components/Slider.vue";

const remote = core.remote();

const router = core.router()


interface UHIDv2 {
  name: string,
  version: string,
  rebooting: boolean,
  updating: boolean,
  progress: number,
  updateStatus: string,
  mac: string,
}


const state = reactive({
  entityId: "",
  online: false,
  rebooting: false,
  alias: "",
  live: true,
  entity: {} as Entity,
  config: {} as Attribute,
  sensorConfig: {
    nightlight: 0,
    indicator: 0,
    sensorRate: 1,
  } as UHIDv2SensorConfig,
  uncommittedConfig: {
    nightlight: 0,
    indicator: 0,
    sensorRate: 1,
  } as UHIDv2SensorConfig,
  configMutex: false,
  attributes: [] as Attribute[],
  triggers: [] as Trigger[],
  sensor: {} as SensorData,
  metadata: {} as UHIDv2,
  configString: "",
  lastUpdate: 0,
  deviceState: [] as Warning[],
  settings: false,
  timestamps: [] as number[],
  data: [] as number[],
})

interface TimeData {

}

onMounted(() => {
  update()
  // let sensor = state.attributes.find(a => a.key == "sensor")

})

onBeforeUnmount(() => {
  state.live = false
})

watchEffect(() => {
  if (state.live) {
    update()
  }

  return remote.attributes
})

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

interface UHIDv2SensorConfig {
  nightlight: number,
  indicator: number,
  sensorRate: number
}

function update() {

  state.entityId = <string>router.currentRoute.value.params["entityId"];

  let entity = remote.entities.find(e => (e.type.includes("uhidv2") || e.type.includes("vradarv4")) && e.id == state.entityId)
  if (!entity) return
  state.entity = entity as Entity;

  state.attributes = remote.attributes.filter(a => a.entity == state.entityId)
  let sensor = state.attributes.find(a => a.key == "sensor")
  if (!sensor) return


  try {
    state.sensor = JSON.parse(sensor?.value || "{}") as SensorData || {} as SensorData
  } catch {
    state.sensor = {} as SensorData
  }
  state.sensor.attributeId = sensor.id
  let metadata = state.attributes.find(a => a.key == "metadata")
  if (!metadata) return
  try {
    state.metadata = JSON.parse(metadata?.value || "{}") as UHIDv2 || {} as UHIDv2
  } catch {
    state.metadata = {} as UHIDv2
  }
  let now = new Date().valueOf()
  if (state.sensor) {
    state.lastUpdate = Math.max(new Date(sensor.updated).valueOf(), new Date(metadata.updated).valueOf())
    state.online = (now - state.lastUpdate) < 1000 * 10
  }
  let config = state.attributes.find(a => a.key == "config")
  if (!config) return
  state.config = config
  let sensorConfig = JSON.parse(state.config.value) as UHIDv2SensorConfig
  if (sensorConfig.sensorRate > 0) {
    state.sensorConfig = sensorConfig

  }

  if (!state.configMutex) {
    state.uncommittedConfig = state.sensorConfig
  }

  let localString = JSON.stringify(state.sensorConfig)
  if (state.configString == "") {
    state.configString = localString
  }
  state.triggers = remote.triggers.filter(t => t.name.includes(state.entity.name))

  if (state.metadata.updateStatus == "startup") {
    state.rebooting = false;
  }

  state.deviceState = createDeviceState()


}

function goBack() {
  router.go(-1);
}

function callTrigger(trigger: Trigger) {
  // triggerService.invoke(trigger).then(() => {
  //
  // }).catch(() => {
  //
  // })
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

function toggleSettings() {
  state.settings = !state.settings
}

function setAlias() {
  entityService.setAlias(state.entityId, state.alias).then(() => {
    state.settings = false
  }).catch(() => {
  })
}

function beginUpdate() {
  attributeService.request({
    id: state.config.id,
    key: "config",
    entity: state.entityId,
    request: "{\"update\": 1}"
  } as Attribute).then(() => {

  }).catch()
}

function pushConfig() {
  attributeService.request({
    id: state.config.id,
    key: "config",
    entity: state.entityId,
    request: JSON.stringify(state.uncommittedConfig)
  } as Attribute).then(() => {
    state.configMutex = false
  }).catch()
}

function updateConfig(nightlight: number, indicator: number) {
  state.configMutex = state.sensorConfig.nightlight != nightlight || state.sensorConfig.indicator != indicator
  state.uncommittedConfig = {
    indicator: indicator,
    nightlight: nightlight,
    sensorRate: 1
  } as UHIDv2SensorConfig

}

function deleteDevice() {
  Promise.all(state.attributes.filter(a => a.entity == state.entityId).map(e => attributeService.delete(e.id))).then(res => {
    entityService.delete(state.entityId).then(() => {
      router.go(-1)
    }).catch((err) => {
      alert("Could not delete entity, but the attributes are gone.. " + err)
    })
  }).catch(() => {
    alert("Could not delete entity attributes")
  })

}

function restartNow() {
  state.rebooting = true
  attributeService.request({
    id: state.config.id,
    key: "config",
    entity: state.entityId,
    request: "{\"restart\": 1}"
  } as Attribute).then(() => {

  }).catch()
}

</script>

<template>


  <div class="scrollable-fixed d-flex flex-column gap-1">
    <Element class="flex-shrink-0">
      <List row>
        <Element :cb="() => goBack()" class="sf d-flex align-items-center justify-content-center" foreground
                 mutable
                 style="width: 3.5rem !important; flex-shrink: 0; flex-grow: 0">
          <div class="d-flex align-items-center justify-content-center" style="width: 3rem">􀆉</div>
        </Element>
        <Element class="d-flex flex-row align-items-start" foreground mutable>

          <div class="sf">{{ state.online ? '􁸑' : '􀙥' }}</div>
          <div class="d-flex flex-column">
            <div class="label-c3 label-o6 label-w600 px-2 mono lh-1" style="">

              {{ ((state.entity.alias) ? state.entity.alias : state.entity.name) }}
            </div>
            <div class="label-c5 label-o4 label-w600 mono px-2 lh-1">
              v{{ state.metadata.version }}
            </div>
          </div>
        </Element>
        <Element :cb="() => toggleSettings()" class="sf d-flex align-items-center justify-content-center" foreground
                 mutable
                 style="width: 3.5rem !important; flex-shrink: 0; flex-grow: 0">
          <div class="d-flex align-items-center justify-content-center" style="width: 3rem">􀍟</div>
        </Element>
      </List>
    </Element>
    <List class="flex-grow-1" scroll-y>
      <Element>
        <List row>


          <Element :cb="() => beginUpdate()" class="d-flex align-items-center justify-content-center gap-1 ox-2"
                   foreground
                   mutable>
            <div class="sf">􁐂</div>
            Update
          </Element>
          <Element :cb="() => restartNow()" class="d-flex align-items-center justify-content-center gap-1 ox-2"
                   foreground
                   mutable>
            <div class="sf">􀊯</div>
            Restart
          </Element>
        </List>
      </Element>

      <Element v-if="state.metadata?.progress != 0">
        <ElementPair alt="This device is currently updating it's firmware..." icon='􀇿'
                     title="Updating...">
          {{ state.metadata.progress }} %
        </ElementPair>
      </Element>
      <Element v-if="state.deviceState.length > 0">
        <List>
          <!--        <ElementPair icon='􀯛' title="Last Update" alt="Elapsed time since last packet">-->
          <!--          <Time :since="state.lastUpdate" precise live></Time>-->
          <!--        </ElementPair>-->
          <ElementPair v-for="warning in state.deviceState" :key="warning.message" :alt="warning.criteria"
                       :title="warning.message"
                       icon='􀇿'>
          </ElementPair>
        </List>

      </Element>

      <Element>
        <List>

          <SensorDOM :sensor-data="state.sensor"></SensorDOM>
        </List>
      </Element>
      <Element>
        <ElementHeader title="Buttons"></ElementHeader>
        <List row>
          <Element v-for="trigger in state.triggers" :key="trigger.id" :cb="() => callTrigger(trigger)"
                   class="px-4 label-c5 label-w600 d-flex justify-content-center align-items-center text-uppercase"
                   foreground micro="􀋦"
                   mutable
                   style="height: 3rem">{{ parseName(trigger.name) }}
          </Element>
        </List>
      </Element>

      <Element>
        <List>
          <div class="d-flex justify-content-between align-items-baseline">
            <ElementHeader title="Indicator"></ElementHeader>
            <div class="label-o4 label-c5 label-w600 px-2">{{ state.uncommittedConfig.indicator }} / 2048</div>
          </div>

          <Slider :change="(v) => updateConfig(state.uncommittedConfig.nightlight, v)" :max="2048" :min="1" :step="1"
                  :value="state.uncommittedConfig.indicator || 0"
                  bg="dim"></Slider>
          <div class="d-flex justify-content-between align-items-baseline">
            <ElementHeader title="Buttons"></ElementHeader>
            <div class="label-o4 label-c5 label-w600 px-2">{{ state.uncommittedConfig.nightlight }} / 2048</div>
          </div>

          <Slider :change="(v) => updateConfig(v, state.uncommittedConfig.indicator)" :max="2048" :min="1" :step="1"
                  :value="state.uncommittedConfig.nightlight || 0"
                  bg="dim"></Slider>
          <ElementHeader title="Commit Changes"></ElementHeader>
          <Element :accent="state.configMutex" :cb="() => pushConfig()"
                   class="d-flex align-items-center justify-content-center gap-1 ox-2"
                   foreground
                   mutable>
            Save Changes

          </Element>
        </List>
      </Element>

    </List>
  </div>

  <div v-if="state.settings" class="modal" @click="toggleSettings">
    <List>
      <Element @click.stop>
        <List>
          <Element class="d-flex flex-column" foreground mutable>
            <div class="label-c3 label-o6 label-w600 px-2 mono lh-1" style="">
              {{ (state.alias != "" ? state.alias : state.entity.name) }}
            </div>
            <div class="label-c5 label-o4 label-w600 mono px-2 ">
              {{ state.alias == "" ? 'No alias' : state.entity.name }}
            </div>
          </Element>

          <List class="" row>
            <ElementTextBox :change="(v) => state.alias = v" class=""></ElementTextBox>
            <Element :cb="() => setAlias()" class="sf d-flex align-items-center justify-content-center " foreground
                     mutable style="width: 4rem !important;">
              <div class="d-flex align-items-center justify-content-center">Save</div>
            </Element>

          </List>

        </List>

      </Element>
      <Element>
        <List>
          <Element :cb="() => deleteDevice()" class="sf d-flex align-items-center justify-content-center " foreground
                   mutable>
            <div class="d-flex align-items-center justify-content-center label-w600">Delete</div>
          </Element>
        </List>
      </Element>

    </List>
  </div>

</template>

<style lang="scss" scoped>
@keyframes modal-open {
  0% {
    transform: scale(0.92);
    opacity: 0;
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}

.modal {
  animation: modal-open 100ms forwards ease-out;
  position: absolute;
  top: -0.2px;
  left: -0.2px;
  width: 100vw;
  height: 100vh;
  padding: 1rem 1rem;
  background-color: rgba(0, 0, 0, 0.05);
  //box-shadow: inset 0 0 16px 8px rgba(0, 0, 0, 0.2);
  backdrop-filter: blur(22px);
  -webkit-backdrop-filter: blur(22px);
  z-index: 1 !important;

}
</style>