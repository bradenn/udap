<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>


import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Attribute, Calendar, Remote} from "@/types";
import moment from "moment";

let remote = inject("remote") as Remote
let preferences = inject('preferences')

let zones = [
  {
    name: "All",
    key: "all"
  },
  {
    name: "Bedroom",
    key: "bedroom"
  },
  {
    name: "Kitchen",
    key: "kitchen"
  },
  {
    name: "Lor",
    key: "lor"
  }
]


let state = reactive({
  calendarAttribute: {} as Attribute,
  calendar: [] as Calendar[],
  loading: true,
})

onMounted(() => {
  state.loading = true


  handleUpdates(remote)
  state.loading = false
})

function handleUpdates(remote: Remote) {
  let cal = remote.attributes.find(a => a.key === "calendar")
  if (!cal) {
    return
  }
  state.calendarAttribute = cal
  let candidate = JSON.parse(state.calendarAttribute.value) as Calendar[]
  state.calendar = candidate.filter(c => isToday(c)).sort((a, b) => new Date(a.start).getMilliseconds() > new Date(b.start).getMilliseconds() ? 1 : -1)
}

watchEffect(() => handleUpdates(remote))

function getTime(time: string): string {
  return moment(time).format("h:mm A")
}

function isToday(cal: Calendar): boolean {
  let days = cal.days.split(",")
  let lookup = ["SU", "MO", "TU", "WE", "TH", "FR", "SA"]
  let today = lookup[new Date().getDay()]
  let now = new Date()
  let current = new Date(cal.end).setHours(now.getHours(), now.getMinutes(), now.getSeconds())
  let reallyIsToday = moment(cal.start).day() == new Date().getDay() && !moment(now).isAfter(cal.end);
  return (days.includes(today) || reallyIsToday) && !moment(current).isAfter(cal.end)
}

function getTimeUntil(time: string): string {
  let now = new Date()
  let from = new Date(time).setHours(now.getHours(), now.getMinutes(), now.getSeconds())
  return moment(time, true).from(from, false)
}

</script>

<template>
  <div v-if="!state.loading" class="d-flex flex-column gap-1" style="width:16rem;">
    <div class="label-xs label-o5 label-w500">Today</div>
    <div v-for="cal in state.calendar" :key="cal.description" class="element p-2 d-flex justify-content-between">
      <div>
        <div class="label-c1 label-r label-w500 label-o5 lh-sm overflow-ellipse">{{ cal.summary }}</div>
        <div class="label-c2 label-o4">{{ cal.description }}</div>
      </div>
      <div class="d-flex flex-column justify-content-center align-items-end">
        <div class="label-c2 label-o3">{{ getTimeUntil(cal.start) }}</div>
        <div class="label-c2 label-o4">{{ getTime(cal.start) }} - {{ getTime(cal.end) }}</div>
      </div>
    </div>
  </div>
</template>


<style lang="scss" scoped>
.overflow-ellipse {
  overflow: hidden;
  white-space: nowrap;
  width: 9rem;

  text-overflow: ellipsis;
}

.context-container {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;

}

</style>