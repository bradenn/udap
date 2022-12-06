<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Attribute, Entity as EntityType, Preferences} from "@/types";
import Slider from "@/components/Slider.vue";
import Entity from "@/components/Entity.vue";
import Switch from "@/components/Switch.vue";
import attributeService from "@/services/attributeService";
import moment from "moment";
import type {Remote} from "@/remote";
import type Core from "@/core";
import core from "@/core";
import MenuSection from "@/components/menu/MenuSection.vue";
import ColorPicker from "@/components/ColorPicker.vue";

const state = reactive({
  entities: [] as EntityType[],
  selected: [] as string[],
  types: [] as any[],

  target: {} as EntityType,
  attributes: [] as Attribute[],
})


let ctx = core.context()

const remote = inject("remote") as Remote
const prefs = inject("preferences") as Preferences

function sortByModule(a: EntityType, b: EntityType): number {
  if (a.name >= b.name) {
    return 1
  } else if (a.name < b.name) {
    return -1
  }
  return 0

}

function compareOrder(a: Attribute, b: Attribute): number {
  if (a.order >= b.order) {
    return 1
  } else if (a.order < b.order) {
    return -1
  }
  return 0
}

onMounted(() => {
  updateData()
})

watchEffect(() => {
  updateData()
  fetchAttributes()
  return remote.attributes
})

function updateData() {
  state.entities = remote.entities.filter(e => e.module === 'govee' || e.module === 'macmeta' || e.module === 'hs100').sort(sortByModule)
  sortEntities()
}

function selectEntity(id: string) {
  // state.selected = id
  if (state.selected.includes(id)) {
    state.selected = state.selected.filter(s => s != id)
  } else {
    state.selected.push(id)
  }
  fetchAttributes()
  if (state.selected.length == 0) {
    state.attributes = []
  }

}

function groupBy<T>(xs: T[], key: string): T[] {
  return xs.reduce(function (rv: any, x: any): T {
    (rv[x[key]] = x);
    return rv;
  }, {});
}

function groupByType<T>(xs: T[], key: string): T[] {
  return xs.reduce(function (rv: any, x: any): T {
    (rv[x[key]] = rv[x[key]] || []).push(x);
    return rv;
  }, {});
}

function sortEntities() {
  state.types = groupByType(state.entities, "module")
}

function fetchAttributes() {
  if (state.selected.length > 0) {
    let a = remote.attributes.filter(a => state.selected.includes(a.entity))
    if (!a) return
    let grouped = groupBy(a, "key")

    state.attributes = Object.values(grouped) as Attribute[]
    state.attributes = state.attributes.sort(compareOrder)

  } else {

    state.selected = []
  }
}


const notify = core.notify()

function attributeRequestAll(key: string, value: string) {
  let attrs = remote.attributes.filter(a => state.selected.includes(a.entity) && a.key === key) as Attribute[]
  if (!attrs) return;

  for (let attr of attrs) {
    attr.request = value
    attributeService.request(attr).then(res => {
    }).catch(err => {
      notify.fail("Core", `Could not change '${attr.key}'to '${attr.request}'`, err)
    })
  }
  // notify.success("Core", `Changed ${key} to '${value}'`)


}


function applyState(attribute: Attribute) {
  attributeRequestAll(attribute.key, attribute.value)
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
  <div class="remote-grid h-100 mt-4">
    <div style="grid-column: 2 / span 3; grid-row: 1 / span 6; ">
      <!--            <FixedScroll style="overflow-y: scroll">-->
      <div class="d-flex flex-column" style="">
        <MenuSection v-for="(v, k) in state.types" :key="k" :title="`${k}`">
          <div class="entity-grid mb-1">
            <Entity v-for="entity in v" :key="entity.id" :entity="entity"
                    :selected="state.selected.includes(entity.id)"
                    @mousedown="() => selectEntity(entity.id)">

            </Entity>
          </div>
        </MenuSection>
      </div>
      <!--            </FixedScroll>-->
    </div>
    <div style="grid-column: 5 / span 7; grid-row: 1 / span 8;">
      <MenuSection :title="'Control'">
        <div class="d-flex gap-1 flex-column">
          <div v-for="attribute in state.attributes" :key="`${attribute.id}`">
            <Switch v-if="attribute.key === 'on'"
                    :active="attribute.value === 'true'"
                    :change="(val) => attributeRequestAll(attribute.key, `${val?'true':'false'}`)"
                    name="Power"></Switch>
            <Slider v-else-if="attribute.key === 'dim'"
                    :change="(val) => {attributeRequestAll(attribute.key, `${val}`)}"
                    :live="false"
                    :max="100"
                    :min="0"
                    :step="5"
                    :value="parseInt(attribute.request)" name="Brightness"
                    unit="%"></Slider>
            <div v-else-if="attribute.key === 'cct'" class="">
              <Slider
                  :change="(val) => {attributeRequestAll(attribute.key, `${colorTemps[val].value}`)}"
                  :live="false"
                  :max="colorTemps.length-1" :min="0" :step="1"
                  :tags="colorTemps.map(c => c.name)"
                  :value="0" name="Color Temperature" unit="K"></Slider>

            </div>
            <div v-else-if="attribute.key === 'hue'">
              <ColorPicker
                  :change="(val) => {attributeRequestAll(attribute.key, `${val}`)}"
                  :selected="parseInt(attribute.value)"></ColorPicker>

            </div>
            <div v-else-if="attribute.key === 'api'">

            </div>
            <div v-else style="display: none;"></div>
          </div>
        </div>
      </MenuSection>
    </div>

  </div>
</template>

<style scoped>

.color-grid {
  align-items: start;
  display: flex;
  flex-direction: row;
  flex-wrap: nowrap;

  padding: 0.25rem;
  gap: 0.25rem;
  width: 100%;
  height: 100%;

}

.remote-controls {
  grid-column: 4 / span 6;
}

.entity-grid {
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-columns: repeat(3, 1fr);
}

.remote-grid {
  display: grid;
  grid-column-gap: 0.5rem;
  grid-row-gap: 0.25rem;
  grid-template-rows: repeat(12, 1fr);
  grid-template-columns: repeat(12, 1fr);
}
</style>