<script lang="ts" setup>
import moment from 'moment'
import {defineProps, onMounted, reactive} from "vue";

let props = defineProps({
  size: String,
  inner: Boolean,
  large: Boolean,
})

let state = reactive({
  time: "",
  date: "",
  timer: 0,
  day: "",
})

onMounted(() => {
  startClock()
})

function formatDay(wd: number) {
  let m = moment()
  let weekday = m.weekday()
  return {
    numeric: m.weekday(weekday + wd).format("DD"),
    long: m.weekday(weekday + wd).format("dddd"),
  }
}

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
  <div>
    <div class="clock-time-inner top" v-html="state.time"></div>
    <div class="clock-date-inner" v-html="state.date"></div>
  </div>
</template>

<style lang="scss" scoped>
.clock-container {
  /*height: 48px !important;*/
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-items: flex-start;
}

.weekday {

}

.clock-large {
  font-size: 3rem !important;
}

.clock-sm {

}


.clock {
  font-family: "SF Compact Display", sans-serif;
  font-size: 0.8rem;
  font-weight: 600;
  line-height: 2rem !important;


  color: rgba(255, 255, 255, 0.8);
  text-shadow: 1px 1px 12px rgba(0, 0, 0, 0.25);
  transition: font-size 0.25s ease-in-out;
}
</style>
