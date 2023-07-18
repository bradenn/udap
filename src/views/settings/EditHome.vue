<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import useListModules, {ModuleController} from "@/controller/modulesController";
import Element from "udap-ui/components/Element.vue";
import {onMounted, reactive} from "vue";
import {v4 as uuid} from "uuid"
import List from "udap-ui/components/List.vue";
import core from "udap-ui/core";
import {PreferencesRemote} from "udap-ui/persistent";

let state = useListModules() as ModuleController;

interface Slot {
  id: string
  value: number
}

let obj = reactive({
  slots: new Map<string, Slot>(),
  items: ["Wobble", "Bobble", "Britz", "Foobs", "Croobs"],
  assignments: [],
  last: "",
  x: 0,
  y: 0,
  cx: {} as any,
  cy: 0,
  active: {} as HTMLElement
})

onMounted(() => {

  for (let i = 0; i < 5; i++) {
    let id = uuid()
    obj.slots.set(id, {
      id: id,
      value: Math.floor(Math.random() * 100)
    } as Slot)
  }
})

function swap(key: string, next: string) {
  // let src = obj.slots.find(i => i.id == key) as Slot
  // let srcIdx = obj.slots.indexOf(src)
  //
  // let dst = obj.slots.find(i => i.id == next) as Slot
  // let dstIdx = obj.slots.indexOf(dst)

  if (obj.slots.has(key) && obj.slots.has(next)) {
    let dst = obj.slots.get(next)
    let src = obj.slots.get(key)
    if (!dst || !src) return
    obj.slots.set(next, src)
    obj.slots.set(key, dst)
  }

  //
  // obj.slots[srcIdx].value = dst.value
  // obj.slots[dstIdx].value = src.value

  // console.log(srcIdx + " -> " + dstIdx)
  //
  // obj.slots = obj.slots.map((a) => {
  //   if (a.id == key) {
  //     a.value = dst.value
  //   } else if (a.id == next) {
  //     a.value = src.value
  //   }
  //
  //   return a
  // })


}

const preferences = core.preferences() as PreferencesRemote

function drag(e: TouchEvent) {
  let touch = e.touches.item(0);
  if (!touch) return;
  let t = touch.target as HTMLDivElement
  e.preventDefault()
  obj.cx = t.parentElement?.id
  // obj.cx = touch.pageX - obj.x
  obj.cy = touch.pageY - obj.y
  // obj.active.style.transitionDuration = `1s`;
  obj.active.style.boxShadow = `inset 0 0 2px 2px ${preferences.accent}`;
  obj.active.style.transform = `translate(${touch.clientX - obj.x}px, ${touch.clientY - obj.y}px)`;
  let stack = document.elementsFromPoint(touch.clientX, touch.clientY).find(e => e.id.includes("-")) as HTMLElement
  if (!stack) return
  if (stack.id != obj.last && obj.last != "") {
    document.getElementById(obj.last)?.style.removeProperty('box-shadow')
  }
  if (stack) {
    obj.last = stack.id
    stack.style.boxShadow = `inset 0 0 2px 2px ${preferences.accent}`
  }

}

function dragStart(e: TouchEvent) {
  let touch = e.touches.item(0)
  if (!touch) return
  obj.active = e.currentTarget as HTMLElement
  obj.x = touch.pageX
  obj.y = touch.pageY


}

function dragStop(e: TouchEvent) {

  obj.active.style.removeProperty('transition-duration');
  obj.active.style.removeProperty('transform');
  obj.active.style.removeProperty('box-shadow');
  if (obj.last != "") {
    let ls = document.getElementById(obj.last)
    if (ls) ls.style.removeProperty('box-shadow')
  }
  // let touch = e.touches.item(0);
  let t = obj.active as HTMLDivElement

  if (obj.cy > t.getBoundingClientRect().height || obj.cy < -t.getBoundingClientRect().height) {
    if (!t.parentElement) {
      obj.x = 0;
      obj.y = 0;

      obj.active = {} as HTMLElement;
      return
    }
    swap(t.parentElement.id, obj.last)
  }

  obj.x = 0;
  obj.y = 0;

  obj.active = {} as HTMLElement;

}

</script>

<template>


  <List>

    <Element v-for="a in obj.slots" :id="a[0]" :key="a[0]" foreground
             style="height: 4rem; padding: 0.25rem !important;">
      <div class="d-flex align-items-center justify-content-center h-100 drag-me "
           style=" border-radius: 0.375rem"
           @touchcancel="dragStop"
           @touchend="dragStop"
           @touchmove="drag" @touchstart="dragStart">
        {{ a[1].value }}

      </div>
    </Element>


  </List>


</template>

<style lang="scss" scoped>

.drag-me {
  transition: 0.5s cubic-bezier(0.2, 1, 0.1, 1);
  user-select: none;
  cursor: move;
  animation: animateIn 500ms ease forwards;
}

@keyframes animateIn {
  0% {
    transform: scale(0.2);
  }
  100% {
    transform: scale(1);
  }
}

.home-grid {
  display: grid;
  grid-template-columns: repeat(1, minmax(1rem, 1fr));
  grid-template-rows: repeat(5, minmax(4rem, 1fr));
  grid-gap: 0.25rem;

  .subplot {


  }
}
</style>