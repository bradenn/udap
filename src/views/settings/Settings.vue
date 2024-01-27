<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import core from "udap-ui/core";
import {PreferencesRemote} from "udap-ui/persistent";
import {reactive} from "vue";
import Navbar from "@/components/Navbar.vue";
import SettingsMenu from "@/views/settings/SettingsMenu.vue";
import Element from "udap-ui/components/Element.vue";
import List from "udap-ui/components/List.vue";

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
  <div class=" nest-highlight">

    <div v-if="preferences.landscape">
      <div class="d-flex flex-row gap-1">
        <Element class="w-25">
          <SettingsMenu></SettingsMenu>
        </Element>
        <div class="flex-fill w-25">
          <Element class="">
            <List v-if="router.currentRoute.value.path !== '/settings'" scroll-y
                  style="max-height: 80vh; height: 80vh">
              <router-view v-slot="{ Component }">
                <transition :name="(router.currentRoute.value.path === '/settings')?'slide-reverse':'slide'"
                            mode="out-in">
                  <component :is="Component"/>
                </transition>
              </router-view>
            </List>
          </Element>

        </div>
      </div>
    </div>
    <div v-else style="height: 100%">
      <List>
        <Element style="height: 100%">
          <Navbar
              :back="(router.currentRoute.value.path === '/settings')?'/':'/settings'"
              :title="state.currentName">
          </Navbar>
        </Element>
        <Element class="">
          <List scroll-y style="max-height: 77vh; height: 77vh">
            <router-view v-slot="{ Component }">
              <transition :name="(router.currentRoute.value.path === '/settings')?'slide-reverse':'slide'"
                          mode="out-in">
                <component :is="Component"/>
              </transition>
            </router-view>
          </List>
        </Element>
      </List>
    </div>
  </div>
</template>

<style lang="scss" scoped>
//.nest-highlight {
//  outline: 1px solid blue;
//  display: block;
//  max-height: 100%;
//  height: 83vh;
//  overflow: hidden;
//}

.nest-highlight > * {
  //outline: 1px solid red;
  border-radius: 0.25rem;
}

.root {

  //overflow-y: auto; //max-height: 50vh; //overflow-y: scroll !important; //height: max-content; //margin: auto; //outline: 1px solid white; //display: block; //display: flex; //flex-direction: column; //max-height: 100%; //height: 100%; //overflow-y: scroll;
}

@keyframes slide-in-reverse {
  from {
    transform: translateX(-25px) scale(0.98);

    opacity: 0.9;
  }
  to {
    transform: translateX(0px) scale(1);
    filter: blur(0px);
    opacity: 1;
  }
}

@keyframes slide-out-reverse {
  from {
    transform: translateX(0px) scale(1);
    filter: blur(0px);
    opacity: 1;
  }
  to {
    transform: translateX(25px) scale(0.98);
    filter: blur(8px);
    opacity: 0;
  }
}

.slide-reverse-enter-active {
  animation: slide-in-reverse 100ms ease-out;
}

.slide-reverse-leave-active {
  animation: slide-out-reverse 100ms ease-out;
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