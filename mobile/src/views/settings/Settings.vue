<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import core from "@/core";
import {computed, reactive} from "vue";
import Navbar from "@/components/Navbar.vue";

const router = core.router()

const state = reactive({
  currentName: ""
})


router.afterEach((to, from, failure) => {
  if (!to.name) return
  state.currentName = <string>to.name
})

let transitionName = computed(() => {
  if (router.currentRoute.value.path === '/home/settings') {
    return 'slide-reverse'
  }
  return 'slide'
})

</script>

<template>
  <div class=" gap-1 h-100" style="height: 100vh">
    <div class="d-flex">
      <Navbar v-if="router.currentRoute.value.name"
              :back="(router.currentRoute.value.path === '/home/settings')?'/home/dashboard':'/home/settings'"
              :title="state.currentName">
      </Navbar>
    </div>
    <div class="root">
      <router-view v-slot="{ Component }">
        <transition :name="(router.currentRoute.value.path === '/home/settings')?'slide-reverse':'slide'"
                    mode="out-in">
          <component :is="Component"/>
        </transition>
      </router-view>
    </div>
  </div>
</template>

<style scoped>


.root {
  display: flex;
  flex-direction: column;
  max-height: 100%;
  height: 100%;
  overflow-y: scroll;
  /*outline: 1px solid red;*/
}

@keyframes slide-in-reverse {
  from {
    transform: translateX(-4%) translateY(-2%) scale(1.025);
    filter: blur(2px);

    opacity: 0;
  }
  to {
    transform: translateX(0%) translateY(0%) scale(1);
    filter: blur(0px);
    /**/
    opacity: 1;
  }
}

@keyframes slide-out-reverse {
  from {
    transform: translateX(0%) translateY(0%) scale(1);
    filter: blur(0px);
    opacity: 1;
  }
  to {
    transform: translateX(4%) translateY(0%) scale(0.7);
    filter: blur(8px);
    opacity: 0;
  }
}

.slide-reverse-enter-active {
  animation: slide-in-reverse 70ms linear;
}

.slide-reverse-leave-active {
  animation: slide-out-reverse 100ms linear;
}

@keyframes slide-in {
  from {
    transform: translateX(5%);
    filter: blur(4px);
    opacity: 0;
  }
  to {
    transform: translateX(0) scale(1);
    filter: blur(0px);
    opacity: 1;
  }
}

@keyframes slide-out {
  from {
    transform: translateX(0) scale(1);
    filter: blur(0px);
    opacity: 1;
  }
  to {
    transform: translateX(-5%);
    filter: blur(4px);
    opacity: 0;
  }
}

.slide-enter-active {
  animation: slide-in 100ms ease-out;
}

.slide-leave-active {
  animation: slide-out 100ms ease-out;
}

/*.slide-enter-to {*/
/*    transform: translateX(0);*/
/*    opacity: 1;*/
/*}*/

/*.slide-leave-to {*/
/*    transform: translateX(10%);*/
/*    opacity: 0;*/
/*}*/
</style>