<script lang="ts" setup>
import type {Remote} from "udap-ui/remote";
import _remote from "udap-ui/remote";
import {onBeforeMount, onBeforeUnmount, onMounted, provide, reactive, watchEffect} from "vue";
import core from "udap-ui/core";
import Status from "@/components/Status.vue";
import Dock from "@/components/Dock.vue";
import {PreferencesRemote, usePersistent} from "udap-ui/persistent";
import Element from "udap-ui/components/Element.vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";
import List from "udap-ui/components/List.vue";
import type {RemoteTimings} from "udap-ui/timings";
import useTimings from "udap-ui/timings";
import useDevice, {Device} from "udap-ui/device";
import StatusBar from "@/components/StatusBar.vue";

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
  // document.body.style.background = preferences.pattern.svg;

  document.body.style.backgroundRepeat = 'repeat';
  // document.body.style.backgroundRepeat = 'round';

  // document.body.style.accentColor = preferences.accent

  document.body.style.backgroundSize = `auto`;
}

let state = reactive({
  ready: false,
  connected: false,
})

function connected(success: boolean): void {
  state.connected = true
}

onBeforeMount(() => {
  updateBackground()
  if (!remote.client.connect(connected)) {
    router.push("/setup/enroll")
    state.ready = true
    return
  } else {

    // state.ready = true
    // router.push("/home/dashboard")
  }
  state.ready = true
})

onMounted(() => {

  if (navigator) {
    //@ts-ignore
    navigator?.clearAppBadge()
  }

  const device: Device = useDevice(remote) as Device;
  provide("device", device)
  if (typeof window !== 'undefined')
      //@ts-ignore
    import('./pwa')
})


onBeforeUnmount(() => {
  remote.client.disconnect()
})


</script>

<template>


  <div class="app-frame">
    <StatusBar class="udap-header"></StatusBar>


    <div v-if="!state.ready && !state.connected"
         class=" flex-fill">
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
    <div v-else class="udap-content scrollable-fixed">
      <router-view v-slot="{ Component }">
        <component :is="Component" class="scrollable-fixed"/>
      </router-view>
    </div>


    <div class="d-flex align-items-center udap-dock pt-1">
      <Element class="w-100">
        <Dock class=""></Dock>
      </Element>
    </div>
    <div class="udap-spacer"></div>

  </div>
</template>

<style lang="scss">
.full-blur {
  z-index: -0;
  width: 100vw;
  height: 100vh;
  backdrop-filter: blur(40px);
}


.udap-header {
  flex: 0 0 32px;
}

.udap-content {
  //flex: 1 1 auto;
  min-height: 0;
  overflow: hidden;
  height: 100%;

}

.udap-dock {
  flex: 0 0 60px;
}

.udap-spacer {
  flex: 0 0 10px;
}

.app-frame {

  height: calc(100vh - 0.25rem);
  //width: 100vw;
  display: flex;
  flex-direction: column;

  //gap: 0.375rem;
  //box-shadow: inset 0 0 2px 2px rgba(255, 128, 255, 1);
  padding: 0.5rem;


  //border-radius: 0.25rem;
}

.app-frame > div {
  //outline: 1px solid rgba(255, 255, 255, 0.5);
  z-index: 1;

  border-radius: 0.25rem;
}

.app-frame > div {
  //box-shadow: 0 0 1px 1px red;
}


//.blur-frame {
//  z-index: -1 !important;
//  position: absolute;
//  top: 0;
//  left: 0;
//  width: 100vw;
//  height: 100vh;
//  background-color: hsla(0, 0, 10, 0.1);
//
//}

//.app-container > * {
//  box-shadow: inset 0 0 2px 2px rgba(128, 128, 255, 1);
//}

.app-container {
  //display: flex;
  //flex-direction: column;
  //height: 100%;
  //justify-content: space-between;
  //grid-template-rows: repeat(10, minmax(1px, 1fr));
  //grid-row-gap: 0.25rem;
  //min-width: 0; //  FF flexbox overflow

  //box-shadow: inset 0 0 2px 2px rgba(128, 255, 128, 1);
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
