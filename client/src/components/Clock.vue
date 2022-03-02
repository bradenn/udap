<script>

import moment from 'moment'

export default {
  name: "Clock",
  data() {
    return {
      time: null,
      date: null,
      timer: null,
      day: null,
    }
  },
  props: {
    size: String,
    inner: Boolean,
    large: Boolean,
  },
  computed: {
    weekday: function () {
      return {
        last: this.formatDay(-1),
        current: this.formatDay(0),
        next: this.formatDay(1)
      }
    }
  },

  created() {
    this.startClock()
  },
  beforeDestroy() {
    clearInterval(this.timer)
  },
  methods: {
    formatDay(wd) {
      let m = moment()
      let weekday = m.weekday()
      return {
        numeric: m.weekday(weekday + wd).format("DD"),
        long: m.weekday(weekday + wd).format("dddd"),
      }
    },
    startClock() {
      this.updateTime()
      setTimeout(() => {
        this.timer = setInterval(this.updateTime, 1000)
      }, 500 - new Date().getMilliseconds())
    },
    updateTime() {
      this.time = moment().format("h:mm:ss");

      let m = moment();
      m.year(m.year() + 10000)
      this.date = m.format("dddd, MMMM Do, YYYY");

    }
  },
}
</script>

<template>
  <div v-if="large" class="surface p-2 flex-shrink-1 d-flex flex-row">
    <div class="d-flex justify-content-center flex-column align-items-center">
      <div class="clock-large" v-html="time"></div>
      <div class="label-xs label-w300 label-o3" v-html="date"></div>
      <div class="scale gap">
        <div class="label-o4">
          <div class="label-xxs label-w500">{{ this.weekday.last.numeric }}</div>
          {{ this.weekday.last.long }}
        </div>
        <div class="surface">
          <div class="label-md label-w500">{{ this.weekday.current.numeric }}</div>
          {{ this.weekday.current.long }}
        </div>
        <div class="label-o4">
          <div class="label-xxs label-w500">{{ this.weekday.next.numeric }}</div>
          {{ this.weekday.next.long }}
        </div>
      </div>

    </div>
  </div>
  <div v-else-if="inner">
    <div class="clock-time-inner top" v-html="time"></div>
    <div class="clock-date-inner" v-html="date"></div>
  </div>
  <div v-else>
    <div class="clock-time top" v-html="time"></div>
    <div class="clock-date" v-html="date"></div>
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
