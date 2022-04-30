<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import moment from "moment";
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Attribute, Entity, Remote} from "@/types";
import type {CurrentWeather, Weather} from "@/weather";
import {getWeatherIcon} from "@/weather";

interface WeatherProps {
  current: CurrentWeather
  latest: Weather
  entity: Entity
  forecast: Attribute
  loading: boolean
  ranges: any
}

let state = reactive<WeatherProps>({
  current: {} as CurrentWeather,
  latest: {} as Weather,
  entity: {} as Entity,
  forecast: {} as Attribute,
  loading: false,
  ranges: {
    temp: {
      min: 100,
      max: -100,
    },
    wind: {
      min: 100,
      max: 0,
    }, rain: {
      min: 0,
      max: 0,
    },
  },
});


onMounted(() => {
  state.loading = true
})

let remote: Remote = inject("remote") || {} as Remote

watchEffect(() => handleUpdates(remote))

function handleUpdates(remote: Remote) {
  if (!remote) return
  let entity = remote.entities.find(e => e.name === 'weather')
  if (!entity) return;
  state.entity = entity
  let attribute = remote.attributes.find(e => e.entity === (entity as Entity).id && e.key === 'forecast')
  if (!attribute) return;
  state.forecast = attribute
  parseWeather(JSON.parse(attribute.value) as Weather)
  return remote.attributes
}

function parseWeather(we: Weather) {

  state.current = we.current_weather

  if (we.hourly.temperature_2m.length <= 0) return
  for (let i = 0; i < we.hourly.temperature_2m.length; i++) {
    if (state.ranges.temp.max < we.hourly.temperature_2m[i]) {
      state.ranges.temp.max = we.hourly.temperature_2m[i]
    } else if (state.ranges.temp.min > we.hourly.temperature_2m[i]) {
      state.ranges.temp.min = we.hourly.temperature_2m[i]
    }

    if (state.ranges.rain.max < we.hourly.precipitation[i]) {
      state.ranges.rain.max = we.hourly.precipitation[i]
    } else if (state.ranges.rain.min > we.hourly.precipitation[i]) {
      state.ranges.rain.min = we.hourly.precipitation[i]
    }

  }


  state.latest = we as Weather
  state.loading = false
}

</script>
<template>
  <div class="element px-2 pt-0">
    <h3>Hourly</h3>
    <div v-if="state.latest.hourly" class="d-flex flex-row justify-content-between align-items-end">
      <div
          v-for="(hour) in Array(12).keys()">
        <div v-if="state.latest.hourly.temperature_2m[hour]"
             :class="new Date().getHours()===hour?'':''"
             class=" d-flex flex-column align-items-center justify-content-center px-3">
          <div class="label-c3 label-w400 label-o3 mt-1">
            {{ moment(new Date().setHours(new Date().getHours() + hour)).format("hA") }}
          </div>
          <div class="label-sm label-o4">
            {{ getWeatherIcon(state.latest.hourly.weathercode[hour], hour) }}
          </div>
          <div class="d-flex align-items-center justify-content-center">
            <div class="label-c3 label-w500 label-o4 mt-1">
              {{ state.latest.hourly.temperature_2m[hour] }}
            </div>

            <div class="label-c3 label-w400 label-o3 mt-1">
              {{ state.latest.hourly_units.temperature_2m }}
            </div>
          </div>
          <div class="d-flex align-items-center justify-content-center">
            <div v-if="state.latest.hourly.precipitation[hour] > 0" class="label-c3 label-w500 label-o4 mt-1 rain">
              {{ state.latest.hourly.precipitation[hour] }}
            </div>
          </div>
        </div>

      </div>

    </div>
    <div>Rain</div>
    <div v-if="state.latest.hourly" class="d-flex flex-row justify-content-between align-items-end">

      <div v-for="(hour) in Array(12).keys()">
        <div v-if="state.latest.hourly.temperature_2m[hour]"
             :class="new Date().getHours()===hour?'':''"
             class=" d-flex flex-column align-items-center justify-content-center px-3">
          <div class="label-c3 label-w400 label-o3 mt-1">
            {{ moment(new Date().setHours(new Date().getHours() + hour)).format("hA") }}
          </div>
          <div class="label-sm label-o4">
            {{ getWeatherIcon(state.latest.hourly.weathercode[hour], hour) }}
          </div>
          <div class="d-flex align-items-center justify-content-center">
            <div class="label-c3 label-w500 label-o4 mt-1">
              {{ state.latest.hourly.precipitation[hour] }}
            </div>

            <div class="label-c3 label-w400 label-o3 mt-1">
              {{ state.latest.hourly_units.precipitation }}
            </div>
          </div>
        </div>
      </div>
    </div>

  </div>


</template>

<style scoped>
.temp-chart {
  display: flex;
  align-items: end;
  gap: 2px;
}

.temp-bar {
  content: ' ';
  width: 6px;
  border-radius: 2px;
  background-color: rgba(255, 255, 255, 0.5);
}
</style>