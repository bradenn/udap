<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "../Element.vue";

import {defineProps, onMounted, reactive} from "vue"
import ElementHeader from "../ElementHeader.vue";
import List from "../List.vue";


const props = defineProps<{
  buckets: string[]
  bins: string[]
  values: string[]
}>()

interface BinItem {
  name: string,
  values: string[]
}

const state = reactive({
  position: {
    x: 0,
    y: 0,
    w: 0,
    id: ""
  },
  map: [] as BinItem[],
  loaded: false,
})

onMounted(() => {
  for (let i = 0; i < props.buckets.length; i++) {
    state.map.push({
      name: props.buckets[i],
      values: [] as string[]
    } as BinItem)
  }

  for (let i = 0; i < props.bins.length; i++) {
    let list = props.buckets.indexOf(props.values[i])
    state.map[list].values.push(props.bins[i])
  }
  state.loaded = true

})

function pointerDown(event: TouchEvent | MouseEvent) {

  state.position.w = 0.99
}

function moveItem<T>(array: T[], currentIndex: number, targetIndex: number): T[] {
  const newArray = [...array];
  const item = newArray.splice(currentIndex, 1)[0];

  if (currentIndex < targetIndex) {
    newArray.splice(targetIndex - 1, 0, item);
  } else {
    newArray.splice(targetIndex, 0, item);
  }

  return newArray;
}

function pointerDrag(event: TouchEvent) {

  let posX = 0, posY = 0;

  let evt = event as TouchEvent

  let touching = evt.currentTarget as EventTarget
  let obj = evt.target as HTMLDivElement
  state.position.id = obj.id

  let bounds = obj.getBoundingClientRect()
  posX = evt.touches.item(0).clientX - bounds.x
  posY = evt.touches.item(0).clientY - bounds.y

  if (posY > bounds.height * 1.5) {
    let id = 0;

    let item = state.map.find((m: BinItem) => obj.id.startsWith(m.name))

    let index = item.values.indexOf(obj.id.split('.')[1])
    moveItem(item.values[item], index, index + 1)
    // state.map = state.map.map(m => m)
  }

  state.position.x = posX
  // state.position.c = angle
  state.position.y = posY

}


function pointerUp(event: TouchEvent | MouseEvent) {

  state.position.w = 1
}


</script>

<template>
  <Element v-if="state.loaded" :foreground="false">
    <ElementHeader></ElementHeader>
    {{ state.position.id }}
    <div @touchdown="pointerDown" @touchmove="pointerDrag" @touchup="pointerUp">
      <div v-for="(bucket) in state.map">
        <ElementHeader :title="bucket.name"></ElementHeader>
        <List :id="`list.${bucket.name}`">
          <Element v-for="bin in bucket.values" :id="`${bucket.name}.${bin}`" :key=bin :foreground="true"
                   :mutable="true">
            {{ bin }}
          </Element>
        </List>
      </div>
    </div>
  </Element>
</template>

<style lang="scss" scoped>

</style>