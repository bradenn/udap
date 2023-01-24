<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>


import Plot from "@/components/plot/Plot.vue";
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Attribute, Entity, Macro, Zone} from "@/types";
import Light from "@/components/widgets/Light.vue";
import Select from "@/components/plot/Select.vue";
import attributeService from "@/services/attributeService";
import Button from "@/components/Button.vue";
import type {Remote} from "@/remote";

let remote = inject("remote") as Remote
let preferences = inject('preferences')

let zones = [
  {
    name: "All",
    key: "all"
  },
  {
    name: "Bedroom",
    key: "bedroom"
  },
  {
    name: "Kitchen",
    key: "kitchen"
  },
  {
    name: "Lor",
    key: "lor"
  }
]


let state = reactive({
  zone: {} as Zone,
  zones: [] as Zone[],
  lights: {} as Entity[],
  macros: {} as Macro[],
  globalColor: 0,
  globalDim: 50,
  globalCCT: 6500,
  targets: ["8c1494c3-6515-490b-8f23-1c03b87bde27", "9a3347a7-7e19-4be5-976c-22384c59142a", "c74d427b-5046-4aeb-8195-2efd05d794f8"] as string[],
  loading: true,
  colorMenu: false,
})

onMounted(() => {
  state.loading = true
  handleUpdates(remote)
  switchZone("all")
})


watchEffect(() => handleUpdates(remote))

function sortLights(a: Entity, b: Entity): number {
  if (a.type < b.type) {
    return -1
  } else {
    return 1
  }
  return 0
}

function sortZones(a: Zone, b: Zone): number {
  if (a.pinned && !b.pinned) {
    return -1;
  } else if (!a.pinned && b.pinned) {
    return 1;
  } else if (a.pinned && b.pinned && a.name > b.name) {
    return 2;
  } else if (a.pinned && b.pinned && a.name < b.name) {
    return -2;
  }
  return 0
}

function handleUpdates(remote: Remote) {
  state.lights = remote.entities.filter((entity: Entity) => state.targets.includes(entity.id)).sort(sortLights)
  state.zones = remote.zones.filter((zone: Zone) => !zone.deleted).sort(sortZones)
  state.macros = remote.macros
  state.loading = false

  if (!state.zone.name) {
    switchZone("all")
  }

  return remote.entities
}

function switchZone(name: string) {
  const zone = remote.zones.find(z => z.name === name)
  if (!zone) return
  state.zone = zone
  state.targets = zone.entities.map(e => e.id)
  handleUpdates(remote)
}

function setAttributes(key: string, value: string) {
  remote.attributes.filter((a: Attribute) => a.key == key && state.targets.includes(a.entity)).forEach(v => {
    v.request = value
    attributeService.request(v)
  })
}

function closeMenu() {
  state.colorMenu = false
}

function changeGlobalColor() {
  setAttributes("hue", `${state.globalColor}`)
}

function changeGlobalDim() {
  setAttributes("dim", `${state.globalDim}`)
}

function changeGlobalCCT() {
  setAttributes("cct", `${state.globalCCT}`)
}

</script>

<template>
  <div v-if="!state.loading" class="d-flex flex-column gap-1">
    <div>
      <Select
          :selected="`${state.zone.name.charAt(0).toUpperCase()}${state.zone.name.substring(1)}`">
        <div v-for="zone in state.zones"
             :class="state.zone.name !== zone.name?'subplot-inline':''"
             class="subplot"
             @click="() => switchZone(zone.name)">
          <div
              class="d-flex align-content-center align-items-center justify-content-center">
            <span class="label-c4 label-o2 lh-1"
                  style="width: 0.75rem; margin-top: 4px;"><span
                v-if="zone.pinned">􀎧</span></span>
            <div class="text-capitalize lh-1">
              {{ zone.name }}
            </div>
          </div>
          <div class="label-c3 label-o3">
            {{ zone.entities.length }} 􀛮
          </div>

        </div>
      </Select>
    </div>
    <div class="light-grid">
      <Light v-for="light in state.lights.slice(0, 5)"
             :key="light.id"
             :entity="light"></Light>
    </div>

    <Plot :cols="5" :rows="1" style="width: 100%;">
      <Button :active="false" text="OFF"
              @click="() => setAttributes('on', 'false')"></Button>
      <Button :active="false" text="ON"
              @click="() => setAttributes('on', 'true')"></Button>
      <Button :active="false" text="􀆫"
              @click="() => setAttributes('dim', '20')"></Button>
      <Button :active="false" text="􀆮"
              @click="() => setAttributes('dim', '60')"></Button>
      <Button :active="!state.colorMenu" :to="`/terminal/remote`"
              text="􀎘"></Button>
    </Plot>

  </div>
  <div v-else>

  </div>
  <div v-if="state.colorMenu"
       class="context-container d-flex flex-column justify-content-start gap-1">
    <div class="element" style="width: 40%;  margin-top: 10rem">
      <div class="d-flex justify-content-between align-items-end">
        <div class="label-xs label-w500 label-o5 px-1">Hue</div>
        <div class="label-c1 label-w500 label-o3 px-1">{{
            state.globalColor
          }}°
        </div>
      </div>
      <input
          v-model="state.globalColor"
          :class="`slider-hue`"
          :max="360"
          :min="0"
          :step="1"
          class="range-slider slider subplot pb-0" style="opacity: 0.9"
          type="range"
          v-on:mouseup="(e) => changeGlobalColor()">
    </div>
    <div class="element" style="width: 40%;">
      <div class="d-flex justify-content-between align-items-end">
        <div class="label-xs label-w500 label-o5 px-1">Color Temperature</div>
        <div class="label-c1 label-w500 label-o3 px-1">{{ state.globalCCT }}°
          K
        </div>
      </div>
      <input
          v-model="state.globalCCT"
          :class="`slider-cct`"
          :max="8000"
          :min="2000"
          :step="1"
          class="range-slider slider subplot pb-0" style="opacity: 0.9"
          type="range"
          v-on:mouseup="(e) => changeGlobalCCT()">
    </div>

    <!--    <Slider name="Brightness" :min="0" :max="100" :step="5" unit="%" style="z-index: 4000 !important;"></Slider>-->
    <!--    <Slider name="Hue" :min="0" :max="360" :step="36" unit="%" style="z-index: 4000 !important;"></Slider>-->
    <div class="element" style="width: 40%;">
      <div class="d-flex justify-content-between align-items-end">
        <div class="label-xs label-w500 label-o5 px-1">Brightness</div>
        <div class="label-c1 label-w500 label-o3 px-1">{{
            state.globalDim
          }}%
        </div>
      </div>

      <input
          v-model="state.globalDim"
          :class="`slider-dim`"
          :max="100"
          :min="0"
          :step="1"
          class="range-slider slider subplot pb-0" style="opacity: 0.9"
          type="range"
          v-on:mouseup="(e) => changeGlobalDim()">
    </div>
  </div>
</template>


<style lang="scss" scoped>

.macro-cell {
  aspect-ratio: 4/2;
  //outline: 1px solid white;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.light-grid {
  width: 100%;
  height: 100%;
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  //grid-template-rows: repeat(5, 1fr);
  grid-template-columns: repeat(2, 1fr);
  grid-template-rows: repeat(5, 1fr);
}


.macro-grid {
  width: 100%;
  height: 100%;
  display: grid;
  grid-column-gap: 0.125rem;
  grid-row-gap: 0.125rem;
  grid-auto-flow: column;
  //grid-template-rows: repeat(5, 1fr);

  grid-template-rows: repeat(1, minmax(1rem, 2.5rem));
}

.context-container {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;

}

.context-container > .element {
  z-index: 8 !important;
}


.color-menu {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  opacity: 1;

}
</style>