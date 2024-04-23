<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import ElementLabel from "udap-ui/components/ElementLabel.vue";
import List from "udap-ui/components/List.vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";
import ElementLink from "udap-ui/components/ElementLink.vue";

import core from "@/core";
import {PreferencesCtrl, PreferencesRemote} from "udap-ui/persistent";
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Endpoint} from "udap-ui/types";


const remote = core.remote()
const preferences = core.preferences() as PreferencesRemote

interface Key {
  key: string
  value: any
  children: Key[]
}

const state = reactive({
  endpoint: {} as Endpoint,
  items: [] as Key[],
  ctx: {} as CanvasRenderingContext2D,
  loaded: false
})


let pctl = inject("preferencesCtrl") as PreferencesCtrl

watchEffect(() => {

  return preferences
})

function localReset() {
  pctl.reset()
  state.items = recursiveKeyMap(preferences, 2)
}

onMounted(() => {

  state.items = recursiveKeyMap(preferences, 2)
  state.loaded = true

})

function recursiveKeyMap(obj: any, depth: number): Key[] {
  let keys = Object.keys(obj)
  let out: Key[] = []
  keys.forEach(key => {
    let k = {} as Key
    k.key = key

    let target = obj[key]
    if (target) {
      if (target instanceof Object) {
        if (depth > 0) {
          k.children = recursiveKeyMap(target, depth - 1)
        }
        k.value = "---"
      } else {
        let value = `${target}`
        k.value = value.slice(0, Math.min(value.length, 20))
        k.children = []
      }
    }
    out.push(k)
  })

  out = out.sort((a, b) => a.key.localeCompare(b.key))
  return out

}


</script>

<template>
  <List>

    <div v-if="state.loaded">
      <ElementHeader title="Local Preferences"></ElementHeader>
      <List>
        <div v-for="item in state.items" :key="item.key">
          <ElementLabel :title="item.key " icon="">
            <div class="label-monospace">{{ item.value }}</div>
          </ElementLabel>
          <List v-if="item.children?item.children.length > 0:false" class="mt-1" style="padding-left: 1rem">
            <ElementLabel v-for="child in item.children" :key="child.key" :title="child.key" icon="">
              <div class="label-monospace">{{ child.value }}</div>
            </ElementLabel>
          </List>
        </div>
      </List>

    </div>
    <ElementHeader title="Options"></ElementHeader>
    <ElementLink :cb="() => localReset()" button class="flex-shrink-0" icon="􀲯"
                 title="Reset"></ElementLink>

  </List>
</template>

<style scoped>
.label-monospace {
  font-family: "JetBrains Mono", monospace;
  font-weight: 400;
  font-size: 0.9rem;
}
</style>