<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {onMounted, reactive} from "vue";
import moment from "moment/moment";
import Clock from "@/components/Clock.vue";


const state = reactive({
  time: "",
  date: "",
  timer: 0
})


onMounted(() => {
  startClock()
})

function startClock() {
  updateTime()
  setTimeout(() => {
    state.timer = setInterval(updateTime, 1000)
  }, 500 - new Date().getMilliseconds())
}

function updateTime() {
  state.time = moment().format("h:mm:ss");
  let m = moment();
  m.year(m.year() + 10000)
  state.date = m.format("dddd, MMMM Do, YYYY");
}

</script>

<template>
  <div class="w-100 d-flex align-items-center justify-content-center flex-column">

    <div class="lock-header">
      <Clock :large="true"></Clock>
    </div>
    <div class="mt-2 label-o4">
      􀎡
    </div>
  </div>
</template>


<style lang="scss">


.lock-header {
  display: flex;
  justify-content: center;
}

.time-lock {
  display: block;
  position: absolute;
  font-family: 'SF Pro Rounded', serif;
  font-style: normal;
  font-weight: 500;
  font-size: 4rem;
  color: #FFF2;
  background: linear-gradient(0deg, rgba(255, 255, 255, 0.35) 20.17%, rgba(255, 255, 255, 0.25) 67.23%), linear-gradient(0deg, rgba(255, 255, 255, 0.35), rgba(255, 255, 255, 0.35));
  -webkit-background-clip: text;

  //background-color: rgba(255, 255, 255, 0.2);


}

.time-lock::before {

}

</style>