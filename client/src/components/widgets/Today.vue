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
  state.calendar = candidate.filter(c => isToday(c)).sort((a, b) => new Date(a.start).getHours() > new Date(b.start).getHours() ? 1 : -1)
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
  <div v-if="!state.loading" class="d-flex flex-column gap-1" style="width:20rem;">

  </div>
</template>


<style lang="scss" scoped>

</style>