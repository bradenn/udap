<script lang="ts" setup>
import Encrypted from "./components/Encrypted.vue";
import Menu from "./components/Menu.vue";
import type {Remote} from "udap-ui/remote";
import _remote from "udap-ui/remote";
import {onBeforeUnmount, onMounted, provide, reactive, watchEffect} from "vue";
import core from "@/core";
import Status from "@/components/Status.vue";
import Dock from "@/components/Dock.vue";
import {PreferencesRemote, usePersistent} from "udap-ui/persistent";

/* Remote */
const remote: Remote = _remote
provide("remote", remote)
const router = core.router();

const preferences: PreferencesRemote = usePersistent();

watchEffect(() => {
  updateBackground()
  return preferences
})

function updateBackground() {

  // document.body.style.backgroundColor = preferences.background
  // document.body.style.backgroundImage = img
  // document.head.style.backgroundSize = 'scale';

  // document.body.style.backgroundSize = "20px";
  document.body.style.backgroundColor = preferences.background
  document.body.style.backgroundImage = preferences.pattern.svg;
  document.body.style.backgroundRepeat = 'repeat';


  document.body.style.backgroundSize = `${preferences.pattern.scale}%`;
}

let state = reactive({
  ready: false
})

onMounted(() => {
  updateBackground()
  if (!remote.client.connect()) {
    router.push("/setup/enroll")
    state.ready = false
    return
  } else {
    state.ready = true
    router.push("/home")
  }

})


onBeforeUnmount(() => {
  remote.client.disconnect()
})


</script>

<template>

  <div class="d-flex flex-column gap-3 mt-2 px-2" style=" max-height: 95vh; height: 95vh; background-color: transparent"
  >

    <div class="d-flex justify-content-between ">
      <div class="d-flex flex-row gap-1 justify-content-start align-items-center align-content-center px-1">
        <Encrypted></Encrypted>
        <div class="udap-logo lh-1" style=" z-index: 6 !important;">UDAP</div>

      </div>
      <Menu v-if="remote.client.connected" :name="0"></Menu>
    </div>

    <div v-if="state.ready && !remote.client.connected"
         class="d-flex flex-column align-items-center justify-content-center"
         style="height: 10rem">
      <div class="label-c1 label-w600">Connection Lost</div>
      <div class="label-c5 label-w600 label-o3">You are currently out of range of all UDAP nodes.</div>
      <Status :remote="remote"></Status>
    </div>
    <div v-else class="dock-fixed">
      <router-view></router-view>
    </div>
    <Dock class="float-end"></Dock>
  </div>
</template>

<style>

.element {

}


.dock-fixed {
  /*position: relative;*/
  height: 100%;
  width: 100%;
  overflow-y: scroll !important;
//padding-right: ;
  /*filter: blur(20px);*/
}

.label-muted {
  font-family: "IBM Plex Sans Medium", sans-serif;
  font-weight: 800;
  font-size: 0.7rem;
  color: rgba(255, 255, 255, 0.25);
}


.udap-logo {
  font-family: "IBM Plex Sans Medium", sans-serf;
  font-size: 2rem;
  font-weight: 700;
}

</style>
