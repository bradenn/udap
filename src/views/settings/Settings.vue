<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import core from "udap-ui/core";
import {PreferencesRemote} from "udap-ui/persistent";
import {reactive} from "vue";
import Navbar from "@/components/Navbar.vue";
import SettingsMenu from "@/views/settings/SettingsMenu.vue";

const router = core.router()
const preferences = core.preferences() as PreferencesRemote

const state = reactive({
  currentName: ""
})


router.afterEach((to, from, failure) => {
  if (!to.name) return
  state.currentName = <string>to.name || "Settings"
})
//
// let transitionName = computed(() => {
//   // if (router.currentRoute.value.path === '/home/settings') {
//   //   return 'slide-up'
//   // }
//   // return 'slide-down'
// })

</script>

<template>
  <div class=" gap-0 d-flex flex-column w-100">
    <div class="d-flex">
      <Navbar v-if="router.currentRoute.value.name"
              :back="(router.currentRoute.value.path === '/home/settings')?'/home/dashboard':'/home/settings'"
              :title="state.currentName">
      </Navbar>
    </div>

    <div v-if="preferences.landscape">
      <div class="d-flex flex-row gap-1">
        <div class="col-4">
          <SettingsMenu></SettingsMenu>
        </div>
        <div class="flex-fill w-25">
          <router-view v-if="router.currentRoute.value.path !== '/home/settings'" v-slot="{ Component }">
            <transition :name="(router.currentRoute.value.path === '/home/settings')?'slide-reverse':'slide'"
                        mode="out-in">
              <component :is="Component"/>
            </transition>
          </router-view>
        </div>
      </div>
    </div>

    <div v-else class="root flex-fill mx-1">

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

//overflow-y: auto; //max-height: 50vh; //overflow-y: scroll !important; //height: max-content; //margin: auto; //outline: 1px solid white; //display: block; //display: flex; //flex-direction: column; //max-height: 100%; //height: 100%; //overflow-y: scroll;
}

@keyframes slide-in-reverse {
  from {
    transform: translateY(-10px);
    filter: blur(4px);
    opacity: 0.5;
  }
  to {
    transform: translateY(0);
    filter: blur(0px);
    opacity: 1;
  }
}

@keyframes slide-out-reverse {
  from {
    transform: translateY(-10px);
    filter: blur(4px);
    opacity: 0.5;
  }
  to {
    transform: translateY(0);
    filter: blur(0px);
    opacity: 1;
  }
}

.slide-reverse-enter-active {
//animation: slide-in 70ms linear;
}

.slide-reverse-leave-active {
//animation: slide-out-reverse 100ms linear;
}

@keyframes slide-in {
  from {
    transform: scale(0.98);
    filter: blur(4px);
    opacity: 0.5;
  }
  to {
    transform: scale(1);
    filter: blur(0px);
    opacity: 1;
  }
}

@keyframes slide-out {
  from {
  //transform: translateY(0) scale(0.98); //filter: blur(0px); opacity: 1;
  }
  to {
  //transform: translateY(10px) scale(0); //filter: blur(1px); opacity: 0;
  }
}

.slide-enter-active {
  animation: slide-in 100ms ease-out;
}

.slide-leave-active {
  animation: slide-out 75ms ease-out;
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