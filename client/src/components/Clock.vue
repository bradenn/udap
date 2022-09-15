<script lang="ts" setup>
import moment from 'moment'
import {onMounted, reactive} from "vue";
import {useRouter} from "vue-router";


let props = defineProps<{
  small?: boolean,
  large?: boolean
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
  state.time = moment().format("hh:mm:ss");
  let m = moment();
  m.year(m.year() + 10000)
  state.date = m.format("dddd, MMMM Do");
}

const router = useRouter()

function currentPageName() {
  let last = router.currentRoute.value.matched.length
  return router.currentRoute.value.matched[last - 1].name
}

</script>

<template>
  <div v-if="props.large">

    <div class="time-xl mt-3">
      <div v-for="c in state.time">
        <div v-if="c === ':'" class="vertical-colon">{{ c }}</div>
        <div v-else>{{ c }}</div>
      </div>
    </div>
    <div class="date" v-html="state.date"></div>
  </div>
  <div v-else>
    <div :class="`time${props.small?'-sm':''}`" class="top" @click="$router.push('/terminal/home')"
         v-html="state.time"></div>
    <div v-if="props.small">
      <div class="page-title"> {{ currentPageName() }}</div>
    </div>
    <div v-if="!props.small" class="date" v-html="state.date"></div>
  </div>
</template>

<style lang="scss" scoped>
.vertical-colon {
  margin-bottom: 30px;
}

.time-xl {
  font-size: 4.5rem;
  line-height: 3.5rem;
  font-style: normal;
  font-weight: 500;

  display: flex;
  align-items: center;
  text-align: center;
  letter-spacing: 0;
  font-family: "SF Pro Rounded", sans-serif;
  background: linear-gradient(0deg, rgba(255, 255, 255, 0.25) 20.17%, rgba(255, 255, 255, 0.15) 67.23%), linear-gradient(0deg, rgba(255, 255, 255, 0.25), rgba(255, 255, 255, 0.25)), rgba(255, 255, 255, 0.4);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;

  //text-shadow: 0 0 4px rgba(0, 0, 0, 0.2);

  mix-blend-mode: screen !important;
  transition: font-size 50ms ease-in;
}

.page-title {
  font-size: 1.9rem;
  line-height: 2rem;
  font-weight: 600;
  font-family: "SF Pro Display", sans-serif;
  color: rgba(255, 255, 255, 0.6);
  text-shadow: 0 0 4px rgba(0, 0, 0, 0.2);
  transition: font-size 50ms ease-in;
}

$dist: 2px;

.time {
  z-index: 22 !important;
  font-size: 2rem;
  line-height: 2rem;
  font-weight: 600;
  font-family: "SF Pro Rounded", sans-serif;
  color: rgba(232, 231, 231, 0.9);
  text-shadow: 0 0 8px rgba(0, 0, 0, 0.9);
  transition: font-size 50ms ease-in;
}

.date {
  font-size: 0.75rem;
  line-height: 0.75rem;
  font-weight: 500;

  color: rgba(255, 255, 255, 0.4);
  font-family: "SF Pro Rounded", sans-serif;
  text-shadow: 0 0 4px rgba(0, 0, 0, 0.2);
  mix-blend-mode: screen;
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
