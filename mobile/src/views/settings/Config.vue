<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import ElementLabel from "udap-ui/components/ElementLabel.vue";
import List from "udap-ui/components/List.vue";
import Title from "udap-ui/components/Title.vue";
import ElementLink from "udap-ui/components/ElementLink.vue";

import core from "@/core";
import {PreferencesCtrl, PreferencesRemote} from "udap-ui/persistent";
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Endpoint} from "udap-ui/types";


const remote = core.remote()

interface Key {
  key: string
  value: any
  children: Key[]
}

const state = reactive({
  endpoint: {} as Endpoint,
  items: [] as Key[],
  ctx: {} as CanvasRenderingContext2D
})

let preferences = core.preferences() as PreferencesRemote

let pctl = inject("preferencesCtrl") as PreferencesCtrl

watchEffect(() => {

  return preferences
})

function localReset() {
  pctl.reset()
  state.items = recursiveKeyMap(preferences)
}

onMounted(() => {

  state.items = recursiveKeyMap(preferences)

})

function recursiveKeyMap(obj: any): Key[] {
  let keys = Object.keys(obj)
  let out: Key[] = []
  keys.forEach(key => {
    let k = {} as Key
    k.key = key

    let target = obj[key]
    if (target) {
      if (target instanceof Object) {
        k.children = recursiveKeyMap(target)
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

onMounted(() => {
  // setup()
  // animate()
})


</script>

<template>
  <div class="d-flex flex-column gap-3">
    <Element>
      <Title title="Local Preferences"></Title>
      <List>
        <div v-for="item in state.items" :key="item.key">
          <ElementLabel :title="item.key " icon="">
            <div class="label-monospace">{{ item.value }}</div>
          </ElementLabel>
          <List v-if="item.children.length > 0" class="mt-1" style="padding-left: 1rem">
            <ElementLabel v-for="child in item.children" :key="child.key" :title="child.key" icon="">
              <div class="label-monospace">{{ child.value }}</div>
            </ElementLabel>
          </List>
        </div>
      </List>

    </Element>
    <Element>
      <Title title="Options"></Title>
      <ElementLink :cb="() => localReset()" button icon="􀲯" title="Reset"></ElementLink>
    </Element>

  </div>
</template>

<style scoped>
.label-monospace {
  font-family: "JetBrains Mono", monospace;
  font-weight: 400;
  font-size: 0.9rem;
}
</style>