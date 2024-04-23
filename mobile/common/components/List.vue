<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {defineProps, onUpdated, reactive} from "vue"
import {v4 as uuidv4} from "uuid";

const props = defineProps<{
  row?: boolean
  fixedWidth?: boolean
  scrollX?: boolean
  scrollY?: boolean
  scrollLock?: boolean
}>()

const state = reactive({
  uuid: uuidv4()
})

onUpdated(() => {
  updatePosition()
})

onUpdated(() => {

})

function updatePosition() {
  let element = document.getElementById(state.uuid) as HTMLElement

  if (element && props.scrollY) {
    let sum = 0;
    // console.log(`List: ${state.uuid} => ${element} ${sum}`)

    if (!element.parentElement) {
      return
    }

    if (!element.parentElement.children) {
      return
    }

    //@ts-ignore
    for (let child of element.parentElement.children) {

      if (child.id.toString() !== element.id.toString()) {
        sum += child.clientHeight

      }

    }

    let parent = element.parentElement.clientHeight;
    let children = parent - sum;
    // element.style.height = `${children}px`;
    // element.style.maxHeight = `${children}px`;
    element.style.overflowY = `scroll`;


  } else {

  }
}

</script>

<template>
  <div :id="`${state.uuid}`"
       :class="`${props.scrollX?'element-list-scroll-x':''} ${props.scrollY?'element-list-scroll-y':''}`"
  >
    <div :class="`element-list ${props.row?'flex-row':'flex-column'} ${props.fixedWidth?'fixed-width':''} `">
      <slot></slot>
    </div>
  </div>
</template>

<style lang="scss">
.element-list-scroll-x {
  overflow-x: scroll;

  .subplot .element {
    flex-shrink: 1;
  }
}

.element-list-scroll-lock.element-list-scroll-y > .element-header {
  margin-bottom: 5vh;
}

.element-list-scroll {


  overflow: hidden;
}

.element-list-scroll-lock.element-list-scroll-y > .subplot {
  scroll-snap-align: start;
}

.element-list-scroll-lock {
  scroll-snap-type: y mandatory;


  .element-list {
    scroll-snap-align: start;
  }
}


.element-list-scroll-y > .element-list {
  padding-bottom: 15vh;
}

.element-list-scroll-y {
  overflow-y: scroll;
  overscroll-behavior: contain;
  -webkit-overflow-scrolling: touch !important;
  min-height: 0 !important;
  height: 100%;
  //height: 100% !important;

}

.element-list .fixed-width {
  div {
    flex-basis: 0 !important;
  }

}

.element-list {
  display: flex;
  //flex: 1 1 auto;
  gap: 0.25rem;


  min-height: 0 !important;
  //padding-inline: 0.125rem;
  margin-inline: 0 !important;

}


</style>