<script lang="ts" setup>
import Encrypted from "./components/Encrypted.vue";
import Menu from "./components/Menu.vue";
import _remote, {Remote} from "@/remote";
import {onMounted, provide} from "vue";
import {Preference} from "@/preferences";
import {PreferenceTypes} from "@/types";
import core from "@/core";

/* Remote */
const remote: Remote = _remote
provide("remote", remote)
const router = core.router();

onMounted(() => {
    let storedToken = new Preference(PreferenceTypes.Token).get()
    if (storedToken === "unset") {
        router.push("/setup/enroll")
    }
    remote.connect()
    router.push("/home")
})


</script>

<template>
    <div class="d-flex flex-column gap-3 mt-2 px-2" style=" max-height: 99vh;">
      <div class="d-flex justify-content-between ">
        <div class="d-flex flex-row gap-1 justify-content-start align-items-center align-content-center px-1">
          <Encrypted></Encrypted>
          <div class="udap-logo lh-1" style=" z-index: 6 !important;">UDAP</div>
        </div>
        <Menu v-if="remote" :name="0"></Menu>
      </div>
      <div class="dock-fixed">
        <router-view></router-view>
      </div>
    </div>
</template>

<style scoped>


.dock-fixed {
  /*position: relative;*/
  height: 100%;
  width: 100%;
  overflow-y: scroll !important;
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
