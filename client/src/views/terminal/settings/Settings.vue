<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>

import {reactive} from "vue";
import router from "@/router/index";

let state = reactive({
  transitionName: 'slide-right',
  to: 0,
  from: 0,
  isDragging: false,
  verified: false,
  dragA: {x: 0, y: 0},
  scrollX: 0,
  scrollXBack: 0
})

// same as beforeRouteUpdate option with no access to `this`
router.beforeResolve((to, from) => {
  // only fetch the user if the id changed as maybe only the query or the hash changed
  let toDepth = to.meta.order;
  state.to = toDepth as number
  let fromDepth = from.meta.order || 0;
  state.from = fromDepth as number

  state.transitionName = state.to < state.from ? "slide-right" : "slide-left";
})

// When called, if the user is still dragging, evoke the action confirming drag intent
function timeout() {
  // Check if user is till dragging
  if (state.isDragging) {
    // Verify the drag intention
    state.verified = true
  }
}


// When the user starts dragging, initialize drag intent
function dragStart(e: MouseEvent) {
  dragStop(e)
  // Record the current user position
  let a = {x: e.screenX, y: e.screenY}
  if (!e.view) return
  // If the drag has started near the bottom of the screen
  // Set the dragging status for later verification
  state.isDragging = true;
  // Record the drag position
  state.dragA = a
  // Verify drag intent if the user is still dragging after 100ms
  setTimeout(timeout, 100)
  // Otherwise, we consider the swipes
}

// While the user is still dragging
function dragContinue(e: MouseEvent) {


  let current = router.currentRoute.value.matched[1].children.map(c => c.path).indexOf(router.currentRoute.value.matched[2].path);

  state.isDragging = true;
  if (state.verified) {
    // Record the current position
    let dragB = {x: e.screenX, y: e.screenY}

    if (!e.view) return

    let height = e.view.screen.availHeight;
    let width = e.view.screen.availWidth;
    let thresholdOffset = 60;

    state.scrollX = (state.dragA.x - dragB.x) / 4

    if (state.scrollX < -100 && current >= 1) {
      router.replace(router.currentRoute.value.matched[1].children[current - 1])
      dragStop(e);
    } else if (state.scrollX > 100 && current < router.currentRoute.value.matched[1].children.length) {
      router.replace(router.currentRoute.value.matched[1].children[current + 1])
      dragStop(e);
    }


  }

}

function dragStop(e: MouseEvent) {
  // Discard the drag intent
  state.isDragging = false;
  // Reset the distance
  // state.scrollX = 0;
  // Reset verified drag intent
  state.verified = false;
  state.dragA = {x: 0, y: 0}

  if (Math.abs(state.scrollX) != 0 && state.scrollXBack == 0) {
    if (state.scrollXBack == 0) {
      state.scrollXBack = setInterval(() => {
        if (state.isDragging) return
        state.scrollX -= state.scrollX * 0.05
        if (Math.abs(state.scrollX) < 0.01) {
          state.scrollX = 0
          state.scrollXBack = 0
          clearInterval(state.scrollXBack)
        }
      }, 5)
    } else {
      clearInterval(state.scrollXBack)
    }
  }
}
</script>

<template>
  <div class="d-flex flex-row align-content-start align-items-start gap h-100"
       style="height: calc(100% - 4rem) !important;">
    <div class="flex-grow-1 h-100">
      <router-view v-slot="{ Component }" class="h-100">
        <component :is="Component"/>

        <!--        <transition :name="state.transitionName"-->
        <!--                    mode="out-in">-->
        <!--        </transition>-->
      </router-view>
    </div>
  </div>
</template>

<style lang="scss" scoped>

.bg {
  width: 100%;
  height: 100%;
  position: absolute;
  top: 0;
  left: 0;
  z-index: -1;
  background-color: rgba(12, 12, 13, 1) !important;
}

.slide-right-enter-active {
  position: relative;
  animation: animateSlideRightIn 75ms ease;
}

.slide-right-leave-active {
  position: relative;
  animation: animateSlideRightOut 75ms ease;
}

.slide-left-enter-active {
  position: relative;
  animation: animateSlideLeftIn 75ms ease;
}

.slide-left-leave-active {
  position: relative;
  animation: animateSlideLeftOut 75ms ease;
}

$displacement: 3rem;
$displacementStart: 0rem;

@keyframes animateSlideRightIn {
  0% {
    right: $displacement;
  }
  100% {
    right: 0;
  }
}

@keyframes animateSlideRightOut {
  0% {
    right: -$displacementStart;
  }
  100% {
    right: -$displacement;
  }
}

@keyframes animateSlideLeftIn {
  0% {
    left: $displacement;
  }
  100% {
    left: 0;
  }
}

@keyframes animateSlideLeftOut {
  0% {
    left: -$displacementStart;
  }
  100% {
    left: -$displacement;
  }
}


/*
  Enter and leave animations can use different
  durations and timing functions.
*/
.slide-fade-enter-active {
  transition: all 600ms ease-out;
}

.slide-fade-leave-active {
  transition: all 300ms cubic-bezier(1, 0.5, 0.8, 1);
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  position: absolute;
  transform: scale(0.98);
}

.sidebar-container {
  width: 14rem;
}
</style>