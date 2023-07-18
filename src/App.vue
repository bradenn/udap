<script lang="ts" setup>
import Encrypted from "./components/Encrypted.vue";
import type {Remote} from "udap-ui/remote";
import _remote from "udap-ui/remote";
import {onBeforeUnmount, onMounted, provide, reactive, watchEffect} from "vue";
import core from "udap-ui/core";
import Status from "@/components/Status.vue";
import Dock from "@/components/Dock.vue";
import {PreferencesRemote, usePersistent} from "udap-ui/persistent";
import Element from "udap-ui/components/Element.vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";
import List from "udap-ui/components/List.vue";
import type {RemoteTimings} from "udap-ui/timings";
import useTimings from "udap-ui/timings";

/* Remote */
const remote: Remote = _remote
provide("remote", remote)
const router = core.router();

const preferences: PreferencesRemote = usePersistent();

const timings: RemoteTimings = useTimings(remote);
provide("timings", timings)
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
  // document.body.style.backgroundRepeat = 'round';


  document.body.style.backgroundSize = `auto`;
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
  <div class="app-frame">
    <div class="app-container">
      <div
          class="d-flex flex-row gap-1 justify-content-start align-items-center align-content-center px-1 pb-1 flex-shrink-0"
          style="height: 1.5rem">
        <Encrypted :compact="!(router.currentRoute.value.fullPath === '/home/dashboard')"></Encrypted>
        <div
            :class="`udap-logo-container ${router.currentRoute.value.fullPath === '/home/dashboard'?'udap-logo':'udap-logo-sm'}`"
            class="lh-1"
            style=" z-index: 6 !important;">UDAP
        </div>

      </div>

      <div v-if="state.ready && !remote.client.connected"
           class="dock-fixed flex-fill">
        <Element>
          <ElementHeader title="Connection Lost"></ElementHeader>
          <List>
            <Element class="d-flex" foreground>
              <Status :remote="remote"></Status>
              <div class=" py-2 px-2 label-c5 label-w600 label-o3">You are currently out of range of all UDAP nodes.
              </div>

            </Element>

          </List>
        </Element>
      </div>
      <div v-else class="dock-fixed flex-grow-1 flex-shrink-1 overflow-hidden">
        <router-view></router-view>
      </div>

      <div class="d-flex align-items-center w-100 flex-shrink-0 mb-3" style="height: 60px">
        <Element class="w-100">
          <Dock class=""></Dock>
        </Element>
      </div>
    </div>
  </div>

</template>

<style lang="scss">
//.app-container > * {
//  box-shadow: inset 0 0 2px 2px rgba(128, 128, 255, 1);
//}

.app-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  justify-content: space-between;
  grid-template-rows: repeat(10, minmax(1px, 1fr));
  grid-row-gap: 0.25rem;
  min-width: 0; //  FF flexbox overflow

  //box-shadow: inset 0 0 2px 2px rgba(128, 255, 128, 1);
}

.app-frame {
  height: 100vh;
  width: 100vw;
  //box-shadow: inset 0 0 2px 2px rgba(255, 128, 255, 1);
  padding: 0.5rem;
  border-radius: 0.25rem;
}

.app-frame > * {
  //outline: 1px solid rgba(255, 255, 255, 0.5);
  z-index: 1;
  border-radius: 0.25rem;
}

.udap-logo-container {
  transition: font-size 100ms ease;
}

//.select-outline-all {
//  //* {
//  //  box-shadow: inset 0 0 1px 1px rgba(255, 255, 255, 0.1) !important;
//  //  border: 0.5px solid rgba(255, 255, 255, 0.1);
//  //  border-radius: 4px;
//  //}
//}


.dock-fixed {
  /*position: relative;*/

  //overflow-y: clip !important;
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
  font-family: "IBM Plex Sans Medium", sans-serif;
  font-size: 2rem;
  font-weight: 700;
}

.udap-logo-sm {
  font-family: "IBM Plex Sans Medium", sans-serif;
  font-size: 1.45rem;
  font-weight: 700;
}
</style>
