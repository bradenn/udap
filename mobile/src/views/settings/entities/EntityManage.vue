<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import core from "@/core";
import Element from "udap-ui/components/Element.vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";
import List from "udap-ui/components/List.vue";
import ElementNumber from "udap-ui/components/ElementNumber.vue";
import type {Attribute, Entity} from "udap-ui/types"
import {onMounted, reactive} from "vue";

const router = core.router()

const state = reactive({
  entityId: "",
  entity: {} as Entity,
  attributes: [] as Attribute[]
})

const remote = core.remote()

const supportedAttributes = [
  {
    key: "hue",
    input: "scalar",
    scalar: {
      min: 0,
      max: 360
    }
  },
  {
    key: "on",
    input: "enum",
    labels: ["ON", "OFF"],
    options: ["true", "false"]
  },
  {
    key: "dim",
    input: "scalar",
    scalar: {
      min: 0,
      max: 360
    }
  },
  {
    key: "cct",
    input: "scalar",
    scalar: {
      min: 0,
      max: 360
    }
  },
]

onMounted(() => {
  let eid = router.currentRoute.value.params["entityId"]
  if (!eid) return
  state.entityId = eid as string
  let e = remote.entities.find(e => e.id == state.entityId)
  if (!e) return;
  state.entity = e
  state.attributes = remote.attributes.filter(a => a.entity == state.entityId)

})

function getTitle(key: string): string {
  switch (key) {
    case "dim":
      return 'Luminance'
    case "hue":
      return 'Hue'
    case "cct":
      return 'Temperature'
    case "on":
      return 'Power'
  }
  return ""
}

function getIcon(key: string): string {
  switch (key) {
    case "dim":
      return '􀇯'
    case "hue":
      return '􀟗'
    case "cct":
      return '􀆮'
    case "on":
      return '􀆨'
  }
  return ""
}

function getValue(key: string, value: string): string {
  switch (key) {
    case "dim":
      return `%`
    case "hue":
      return `&deg;`
    case "cct":
      return `&deg; K`
    case "on":
      return value === "true" ? "ON" : "OFF"
  }
  return ""
}

</script>

<template>
  <List>

    <Element class="p-3" foreground>
      <div class="label-c3 lh-1 label-o5 d-flex align-items-start">
        <div class="sf-icon label-c5" style="width: 15px; margin-right: 6px">
          {{ state.entity.icon }}
        </div>
        <div>
          {{ state.entity.alias ? `${state.entity.alias}ª` : state.entity.name }}
          <div class="label-c5 lh-1 label-o3">{{ state.entity.type }}</div>
        </div>
      </div>
    </Element>

    <ElementHeader class="px-3 pb-0 m-0" title="Attributes"></ElementHeader>
    <List>
      <ElementNumber v-for="attr in state.attributes.sort((a, b) => a.order - b.order)" :icon="getIcon(attr.key)"
                     :immutable="attr.key === 'on'" :title="getTitle(attr.key)"
                     :unit="getValue(attr.key, attr.value)"
                     :value="(parseInt(attr.value) || 0)">

      </ElementNumber>
    </List>
  </List>
</template>

<style lang="scss" scoped>
.entity {

}
</style>