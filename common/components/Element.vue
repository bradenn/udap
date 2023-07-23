<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {defineProps, inject, onMounted, onUpdated, reactive} from "vue"
import {useRouter} from "vue-router";
import {PreferencesRemote} from "../persistent";

const props = defineProps<{
  foreground?: boolean
  mutable?: boolean
  surface?: boolean
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
}

onUpdated(() => {
  update()
})

router.afterEach(update)

function pointerDown(event: TouchEvent | MouseEvent) {
  state.position.w = 0.99
  state.down = Date.now().valueOf()
  // if (event.cancelable) event.preventDefault()
  if (props.longCb) {
    //@ts-ignore
    state.timeout = setTimeout(revisitLongCb, 250) as number
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

  <div
      :class="`${props.foreground?(props.link?(state.active?'subplot px-2 py-2':'subplot  px-2 py-2'):'subplot px-2 py-2'):`element back`} ${props.surface?'subplot-surface':''} `"
      :style="`box-shadow: inset 0 0 1px 1px ${state.active||props.accent?preferences.accent+'':'transparent'}; ${props.mutable?'transform: scale(' + state.position.w+');':''}  ${!props.foreground?`backdrop-filter: blur(${preferences.blur}px); -webkit-backdrop-filter: blur(${preferences.blur}px);`:''}`"

      @touchend.passive="pointerUp"
      @touchleave.passive="pointerUp"
      @touchmove.passive="pointerDrag"
      @touchstart.passive="pointerDown">
    <slot></slot>
  </div>
</template>

<style lang="scss" scoped>
.subplot-inactive {
  background-color: transparent !important;
}

.subplot.disabled {
  background-color: rgba(255, 0, 0, 1);
}

.subplot {
  cursor: none !important;
  user-select: none !important;
}

.element {
  cursor: none !important;
  user-select: none !important;
  padding: 0.375rem;
  border-radius: 1rem !important;
  box-shadow: inset 0 0 1px 0.5px rgba(255, 255, 255, 0.05) !important;


  > .subplot {
    backdrop-filter: none !important;
    -webkit-backdrop-filter: none !important;
  }
}

.subplot {
  cursor: none !important;
  user-select: none !important;
  -webkit-user-select: none !important;
  //z-index: 99 !important;
  border-radius: calc(1rem - 0.375rem) !important;
  border: 1px solid rgba(255, 255, 255, 0.025);
  //scroll-padding: 1rem;
  //filter: brightness(150%);
  transition: transform 80ms ease-out;
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