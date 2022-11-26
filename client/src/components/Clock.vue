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
    page: "",
    parent: "",
})

onMounted(() => {
    startClock()
    currentPageName()
})

const router = useRouter()
router.afterEach(() => {
    currentPageName()
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
    state.date = m.format("ddd, MMM Do, YYYY");
}


function currentPageName() {
    let last = router.currentRoute.value.matched.length
    let current = router.currentRoute.value.matched[last - 1]
    let meta = current.meta
    state.parent = ""
    if (meta) {
        if (meta.title) {
            if (!current) return
            if (!current.meta) return
            state.page = current.meta.title as string

            return
        }
    }

    state.page = current.name as string || ""

}

</script>

<template>
    <div v-if="!props.small">

        <div>

            <div class="d-flex align-items-center time-xl">
                {{ state.time }}
            </div>
            <div class="date" v-html="state.date"></div>
        </div>
    </div>
    <div v-else class="h-100">
        <div class="d-flex align-items-start flex-column justify-content-center time-xl h-100 px-1">
            {{ state.time }}
        </div>
    </div>
</template>

<style lang="scss">
.text-change {

}

@keyframes clockIn {
  0% {
    transform: scale(0.96);
    opacity: 0;
  }
  20% {
    opacity: 0.8;
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}

.vertical-colon {
  margin-bottom: 30px;
}

.time-sm {

  font-feature-settings: "ss03"; // Vertically centers the colon for macos
  font-size: 1.6rem;
  line-height: 1.6rem;
  font-style: normal;
  font-weight: 600;

  display: flex;
  align-items: center;
  text-align: center;
  font-family: "SF Compact Display", sans-serif !important;

  //text-shadow: 0 0 4px rgba(0, 0, 0, 0.2);
  color: white;

  mix-blend-mode: screen !important;
  transition: font-size 50ms ease-in;
}

.time-xl {
  //font-variant-numeric: tabular-nums;
  font-feature-settings: "ss03"; // Vertically centers the colon for macos
  font-size: 1.4rem;
  line-height: 1.4rem;
  font-weight: 600;

  //outline: 1px solid red;
  display: flex;
  align-items: center;
  text-align: center;

  justify-content: left;

  font-family: "SF Pro Display", monospace !important;

  //text-shadow: 0 0 4px rgba(0, 0, 0, 0.2);
  color: rgba(255, 255, 255, 0.6);

  mix-blend-mode: screen !important;
  transition: font-size 50ms ease-in;
}

.page-title {
  font-size: 1.6rem;
  line-height: 1.6rem;
  font-weight: 700;
  font-family: "SF Pro Display", sans-serif;
  color: rgba(220, 220, 220, 0.9);
  text-shadow: 0 0 4px rgba(0, 0, 0, 0.2);
  transition: font-size 50ms ease-in;
}

$dist: 2px;

.time {
  z-index: 22 !important;
  font-size: 2.1rem;
  line-height: 2rem;
  font-weight: 600;
  font-family: "SF Pro Display", sans-serif;
  color: rgba(255, 255, 255, 0.6);
  text-shadow: 0 0 8px rgba(0, 0, 0, 0.2);
  animation: clockIn 100ms forwards;
}

.date {
  font-size: 0.65rem;
  line-height: 0.65rem;
  font-weight: 500;


  color: rgba(255, 255, 255, 0.4);
  font-family: "SF Compact Rounded", sans-serif !important;
  text-shadow: 0 0 4px rgba(0, 0, 0, 0.1);

}

.date-sm {
  font-size: 0.75rem;
  line-height: 0.75rem;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.5);
  font-family: "SF Pro Display", sans-serif;
  text-shadow: 0 0 4px rgba(0, 0, 0, 0.2);
}
</style>
