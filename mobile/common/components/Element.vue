<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {defineProps, inject, onMounted, onUpdated, reactive} from "vue"
import {useRouter} from "vue-router";
import {PreferencesRemote} from "../persistent";
import {v4 as uuidv4} from "uuid";

const props = defineProps<{
  foreground?: boolean
  mutable?: boolean
  layout?: boolean
  micro?: string
  surface?: boolean
  fill?: boolean
  to?: string
  link?: boolean
  cb?: () => void
  longCb?: () => void
  accent?: boolean
}>()

const state = reactive({
  position: {
    w: 1,
    x: 0,
    y: 0,
    a: 0,
    b: 0,
    c: 0,
  },
  dx: 0,
  down: 0,
  timeout: 0,
  active: false,
  uuid: uuidv4()
})

const router = useRouter();
const preferences = inject("preferences") as PreferencesRemote

onMounted(() => {
  if (props.accent || props.link || props.cb) {
    update()
  }
})

function update() {
  if (props.to) {
    state.active = router.currentRoute.value.fullPath == (props.to)
  }
  updateSize()
}

function updateSize() {
  let element = document.getElementById(state.uuid) as HTMLElement

  if (element && props.fill) {
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

    let parent = element.parentElement?.clientHeight || 100;
    let children = parent - sum;
    element.style.height = `${children}px`;
    element.style.maxHeight = `${children}px`;


  } else {

  }
}

onUpdated(() => {
  update()

})

router.afterEach(update)

function pointerDown(event: TouchEvent | MouseEvent) {
  state.position.w = 0.995
  state.down = Date.now().valueOf()
  // if (event.cancelable) event.preventDefault()
  if (props.longCb) {
    //@ts-ignore
    state.timeout = setTimeout(revisitLongCb, 400) as number
  }
}

function revisitLongCb() {
  if (props.longCb) {
    props.longCb()
  }
}

function pointerDrag(event: TouchEvent | MouseEvent) {

  let posX = 0, posY = 0;
  // if (event.cancelable) event.preventDefault()
  let evt = event as TouchEvent
  let touching = event.currentTarget as HTMLDivElement


  let bounds = touching.getBoundingClientRect()
  if (event instanceof TouchEvent) {
    let itm = evt.touches.item(0)
    if (!itm) return
    posX = itm.clientX - bounds.x
    posY = itm.clientY - bounds.y
  } else if (event instanceof MouseEvent) {
    let itm = event as MouseEvent

    if (!itm) return
    posX = itm.clientX - bounds.x
    posY = itm.clientY - bounds.y
  }


  state.dx++

  posX = Math.max(0, Math.min(posX, bounds.width))
  posY = Math.max(0, Math.min(posY, bounds.height))

  state.position.a = posX / bounds.width
  state.position.b = posY / bounds.height
  state.position.c = posY

  let tx = 5, ty = 20;
  let angleX = Math.atan((posX - (bounds.width / 2)) / bounds.width)
  let angleY = Math.atan((posY - (bounds.height / 2)) / bounds.height)

  state.position.b = (angleX * (180.0 / Math.PI)) * 0.7
  state.position.a = -(angleY * (180.0 / Math.PI)) * 0.7

  state.position.x = posX
  // state.position.c = angle
  state.position.y = posY


}


function pointerUp(event: TouchEvent | MouseEvent) {
  state.position.a = 0
  state.position.b = 0
  // event.preventDefault()
  state.position.w = 1
  if (props.longCb) {
    clearTimeout(state.timeout);
  }
  if ((props.to || props.cb) && (Date.now().valueOf() - state.down) < 200 && state.dx < 2) {
    if (props.to) {
      router.push(props.to)
    } else if (props.cb) {
      props.cb()
    }
  }
  state.dx = 0
}


</script>

<template>

  <div :id="`${state.uuid}`"
       :class="`${props.layout?'element-layout':''} ${props.foreground?(props.link?(state.active?'subplot px-2 py-2':'subplot  px-2 py-2'):'subplot px-2 py-2'):`element back`} ${props.surface?'subplot-surface':''} `"
       :style="`box-shadow: inset 0 0 0.5px 1px ${state.active||props.accent?preferences.accent+'':'transparent'}; ${props.accent?'background-color: rgba(255,255,255,0.05);':''} ${props.mutable?'transform: scale(' + state.position.w + `);`:''}  ${props.foreground?`${state.position.w!=1?'background-color: rgba(255,255,255,0.05);':''} `:`backdrop-filter: blur(${preferences.blur}px); -webkit-backdrop-filter: blur(${preferences.blur}px);`}`"

       @touchend.passive="pointerUp"
       @touchleave.passive="pointerUp"
       @touchmove.passive="pointerDrag"
       @touchstart.passive="pointerDown">
    <slot></slot>
  </div>
</template>

<style lang="scss" scoped>
//backdrop-filter: blur(${preferences.blur}px); -webkit-backdrop-filter: blur(${preferences.blur}px);
.subplot-inactive {
  background-color: transparent !important;
}

.micro-corner {
  //background-color: rgba(255, 255, 255, 0.2);
  position: absolute;
  border-radius: 50%;
  width: 16px;
  aspect-ratio: 1/1;
  left: 1.85rem;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.3);
}

.subplot.disabled {
  background-color: rgba(255, 0, 0, 1);
}


.subplot {
  cursor: none !important;
  user-select: none !important;
}

$border-element: hsla(0deg, 0%, 20%, 0.4);
.element-layout {
  min-height: 0 !important;
  overflow: hidden !important;
  display: block !important;
}

.element {


  overflow: clip;
  cursor: none !important;
  user-select: none !important;
  padding: 0.375rem 0.375rem;
  //z-index: 1000 !important;
  //transform: translate3d(0, 0, 0);
  //background-color: rgba(22, 22, 24, 0.05);
  //backdrop-filter: brightness(120%);
  //-webkit-backdrop-filter: brightness(120%);
  border-radius: 1rem !important;
  box-shadow: inset 0 0 2px 0.5px $border-element !important;
  //border: 1px solid $border-element;
  //border: 1px solid hsla(0deg, 0, 17, 0.2);
  //filter: drop-shadow(2px 0 0 red);
  //background-blend-mode: overlay !important;

  > .subplot {
    //backdrop-filter: brightness(190%) !important;
    //
    //-webkit-backdrop-filter: brightness(190%) !important;
  }
}

.subplot {
  cursor: none !important;
  user-select: none !important;
  -webkit-user-select: none !important;
  transform: translate3d(0, 0, 0);

  z-index: 0 !important;
  //margin-inline: 0;
  //padding: 0 !important;
  border-radius: calc(1rem - 0.375rem) !important;
  border: 1px solid hsla(0deg, 0, 22, 0.3);
  //scroll-padding: 1rem;
  //background-blend-mode: overlay !important;
  //box-shadow: inset 0 0 1px 0.5px hsla(0deg, 0, 14, 0.8), 0 0 8px 2px hsla(0deg, 0, 2, 0.1), 0 0 1px 0.5px hsla(0deg, 0, 12, 0.2) !important;
  //background: rgba(255, 255, 255, 0.03);
  //box-shadow: 0 8px 32px 0 rgba(31, 38, 135, 0.37);
  //backdrop-filter: blur(8px) brightness(120%);
  //-webkit-backdrop-filter: blur(8px) contrast(98%);

  //filter: brightness(150%);
  transition: transform 40ms ease-out;

  //
  //
  //-webkit-backdrop-filter: brightness(190%) !important;
  //backdrop-filter: brightness(190%) !important;
}

.subplot {

}

.flex-column > .subplot:last-of-type:not(:first-of-type) {
  flex-grow: 1;
  border-radius: 0.375rem 0.375rem calc(1rem - 0.375rem) calc(1rem - 0.375rem) !important;

  .subplot {

    border-radius: 0.25rem 0.125rem 0.3125rem 0.25rem !important;

  }
}

.flex-column > .subplot:not(:first-of-type, :last-of-type) {
  border-radius: 0.375rem 0.375rem 0.375rem 0.375rem !important;
  flex-grow: 1;

  .subplot {
    border-radius: 0.25rem 0.125rem 0.125rem 0.25rem !important;

  }
}

.flex-column > .subplot:first-of-type:not(:last-of-type) {
  flex-grow: 1;
  border-radius: calc(1rem - 0.375rem) calc(1rem - 0.375rem) 0.375rem 0.375rem !important;

  .subplot {
    border-radius: 0.25rem 0.25rem 0.125rem 0.25rem !important;

  }
}

.flex-row > .subplot:last-of-type:not(:first-of-type) {
  flex-grow: 1;
  border-radius: 0.375rem calc(1rem - 0.375rem) calc(1rem - 0.375rem) 0.375rem !important;


}

.flex-row > .subplot:not(:first-of-type, :last-of-type) {
  border-radius: 0.375rem 0.375rem 0.375rem 0.375rem !important;
  flex-grow: 1;

  .subplot {
    border-radius: calc(1rem - 0.125rem) 0.125rem 0.125rem calc(1rem - 0.125rem) !important;
  }
}

.flex-row > .subplot:first-of-type:not(:last-of-type) {
  flex-grow: 1;
  border-radius: calc(1rem - 0.375rem) 0.375rem 0.375rem calc(1rem - 0.375rem) !important;

  .subplot {
    border-radius: calc(1rem - 0.125rem) 0.125rem 0.125rem calc(1rem - 0.125rem) !important;
  }
}


</style>