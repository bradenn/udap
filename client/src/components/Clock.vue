<script lang="ts" setup>
import moment from 'moment'
import {onMounted, reactive} from "vue";


let props = defineProps<{
  small?: boolean,
}>()

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
    <div :class="`time${props.small?'-sm':''}`" class="top" @click="$router.push('/terminal/home')"
         v-html="state.time"></div>
    <div v-if="props.small">
      <div class="page-title"> {{ $router.currentRoute.value.matched[1].name }}</div>
    </div>
    <div v-if="!props.small" class="date" v-html="state.date"></div>
  </div>
</template>

<style lang="scss" scoped>
.page-title {
  font-size: 2rem;
  line-height: 2rem;
  font-weight: 600;
  font-family: "SF Pro Rounded", sans-serif;
  color: rgba(255, 255, 255, 0.6);
  text-shadow: 0 0 4px rgba(0, 0, 0, 0.2);
  transition: font-size 50ms ease-in;
}

.time {
  z-index: 22 !important;
  font-size: 2rem;
  line-height: 2rem;
  font-weight: 600;
  font-family: "SF Pro Rounded", sans-serif;
  color: rgba(255, 255, 255, 0.6);
  text-shadow: 0 0 4px rgba(0, 0, 0, 0.2);
  transition: font-size 50ms ease-in;
}

.date {
  font-size: 0.75rem;
  line-height: 0.75rem;
  font-weight: 500;

  color: rgba(255, 255, 255, 0.5);
  font-family: "SF Pro Rounded", sans-serif;
  text-shadow: 0 0 4px rgba(0, 0, 0, 0.2);
}

.time-sm {
  font-size: 0.75rem;
  line-height: 0.75rem;
  font-weight: 500;

  color: rgba(255, 255, 255, 0.5);
  font-family: "SF Pro Rounded", sans-serif;
  text-shadow: 0 0 4px rgba(0, 0, 0, 0.2);
}

.date-sm {
  font-size: 0.75rem;
  line-height: 0.75rem;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.5);
  font-family: "SF Pro Rounded", sans-serif;
  text-shadow: 0 0 4px rgba(0, 0, 0, 0.2);
}
</style>
