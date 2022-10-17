<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import PaneMenu from "@/components/pane/PaneMenu.vue";
import type {Attribute, Entity, Remote} from "@/types";
import Slider from "@/components/Slider.vue";
import Switch from "@/components/Switch.vue";
import PaneMenuToggle from "@/components/pane/PaneMenuToggle.vue";
import attributeService from "@/services/attributeService";
import moment from "moment";

const state = reactive({
  entities: [] as Entity[],
  selected: "",
  target: {} as Entity,
  attributes: [] as Attribute[],
})

const remote = inject("remote") as Remote

function sortByModule(a: Entity, b: Entity): number {
  if (a.name >= b.name) {
    return 1
  } else if (a.name < b.name) {
    return -1
  }
  return 0

}

onMounted(() => {
  updateData()
})

watchEffect(() => {
  updateData()
  return remote.entities
})

function updateData() {
  state.entities = remote.entities.filter(e => e.type === 'spectrum' || e.type === 'dimmer').sort(sortByModule)
}

function selectEntity(id: string) {
  state.selected = id
  state.target = state.entities.find(e => e.id === id)
  state.attributes = remote.attributes.filter(a => a.entity === id).sort((a, b) => a.order - b.order)
}


function applyState(attribute: Attribute) {
  attributeService.request(attribute)
}

const colorTemps = [
  {
    name: "Warm",
    value: 2400,
    color: "255, 157, 63"
  },
  {
    name: "Tungsten",
    value: 3000,
    color: "255, 180, 107"
  },
  {
    name: "Moonlight",
    value: 4100,
    color: "255, 211, 168"
  }, {
    name: "Daylight",
    value: 5000,
    color: "255, 236, 224"
  },
  {
    name: "Cloudy",
    value: 5500,
    color: "255, 236, 224"
  },
  {
    name: "White",
    value: 6500,
    color: "255, 249, 253"
  }
]

function timeSince(time: string): string {
  return moment(time).fromNow()
}

</script>

<template>
  <div class="remote-grid">
    <PaneMenu title="Devices">
      <PaneMenuToggle v-for="entity in state.entities" :active="state.selected === entity.id"
                      :fn="() => selectEntity(entity.id)"
                      :subtext="entity.module"
                      :title="entity.name"
      ></PaneMenuToggle>
    </PaneMenu>

    <div class="d-flex flex-column remote-controls ">
      <div class="d-flex gap-1 flex-column">
        <div v-for="attribute in state.attributes" :key="`${attribute.id}`">
          <Switch v-if="attribute.key === 'on'" :active="false" name="Power"></Switch>
          <Slider v-else-if="attribute.key === 'dim'"
                  :change="(val) => {attribute.request = `${val}`; applyState(attribute)}" :live="false" :max="100"
                  :min="0"
                  :step="5"
                  :value="parseInt(attribute.request)" name="Brightness" unit="%"></Slider>
          <div v-else-if="attribute.key === 'cct'" class="">
            <Slider :change="(val) => {attribute.request = `${val}`; applyState(attribute)}" :max="colorTemps.length-1"
                    :min="0" :step="1" :tags="colorTemps.map(c => c.name)"
                    :value="0"
                    confirm
                    name="Color Temperature" unit="K"></Slider>

          </div>
          <div v-else-if="attribute.key === 'api'">

          </div>
        </div>
      </div>

    </div>
    <div class="d-flex flex-column gap-1">
      <div class="element p-2">
        <div class="d-flex label-o4 label-c2 justify-content-between">
          <div class="label-w600 label-r label-c2">
            Last Updated
          </div>
          <div class="label-w500 label-o3 label-c2">
            {{ timeSince(state.target.updated) }}
          </div>
        </div>
        <div class="d-flex label-o4 label-c2 justify-content-between">
          <div class="label-w600 label-r label-c2">
            Created
          </div>
          <div class="label-w500 label-o3 label-c2">
            {{ timeSince(state.target.created) }}
          </div>
        </div>
        <div class="d-flex label-o4 label-c2 justify-content-between">
          <div class="label-w600 label-r label-c2">
            Neural
          </div>
          <div class="label-w500 label-o3 label-c2">
            {{ state.target.neural }}
          </div>
        </div>
        <div class="d-flex label-o4 label-c2 justify-content-between">
          <div class="label-w600 label-r label-c2">
            Alias
          </div>
          <div class="label-w500 label-o3 label-c2">
            {{ state.target.alias }}
          </div>
        </div>
      </div>
      <div class="element d-flex d-flex justify-content-between">
        <div class="label-w500 label-c1 label-o4 px-1">Neural</div>
        <div class="label-w500 label-c1 label-o3 px-1">{{ state.target.neural }}</div>
      </div>
    </div>

  </div>
</template>

<style scoped>

.remote-controls {
  grid-column: span 2 / 4;
}

.remote-grid {
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-rows: repeat(1, 1fr);
  grid-template-columns: repeat(4, 1fr);
}
</style>