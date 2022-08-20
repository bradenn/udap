<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import moment from "moment";
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Attribute, Entity, Remote} from "@/types";
import type {CurrentWeather, Weather} from "@/weather";
import {getWeatherIcon, getWeatherState} from "@/weather"

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

let remote = inject("remote") as Remote

watchEffect(() => handleUpdates(remote))

function handleUpdates(remote: Remote) {
  if (!remote) return
  let entity = remote.entities.find(e => e.name === 'weather')
  if (!entity) return;
  state.entity = entity as Entity
  let attribute = remote.attributes.find(e => e.entity === state.entity.id && e.key === 'forecast')
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

function roundDecimal(input: number, places: number) {
  return Math.round(input * Math.pow(10, places)) / Math.pow(10, places)
}

</script>
<template>
  <div class="element p-2 pt-1">
    <div class=" d-flex flex-row align-items-center">
      <div class="flex-shrink-1 " style="min-width: 9rem; padding-left: 0.25rem">
        <h2 class="lh-md">{{ roundDecimal(state.latest.current_weather.temperature, 0) }}° F</h2>
        <div class="label-c1 label-r label-w400 label-o4 lh-1">{{
            getWeatherState(state.latest.current_weather.weathercode)
          }}
        </div>
        <div class="label-c1 label-r label-o3 label-w400 ">High {{ Math.round(state.ranges.temp.max) }}° • Low
          {{ Math.round(state.ranges.temp.min) }}°
        </div>
      </div>
      <div v-if="state.latest.hourly" class="d-flex flex-row justify-content-between align-items-end"
           style="width: 100%">
        <div
            v-for="(hour) in Array(12).keys()">
          <div v-if="state.latest.hourly.temperature_2m[hour]"
               :class="new Date().getHours()===hour?'':''"
               class=" d-flex flex-column align-items-center justify-content-center px-3">
            <div class="label-c3 label-w400 label-o3 mt-1">
              {{ moment(new Date().setHours(new Date().getHours() + hour)).format("hA") }}
            </div>
            <div class="label-sm label-o4">
              {{ getWeatherIcon(state.latest.hourly.weathercode[new Date().getHours() + hour], hour) }}
            </div>
            <div class="d-flex align-items-center justify-content-center">
              <div class="label-c3 label-w500 label-o4 mt-1">
                {{ state.latest.hourly.temperature_2m[new Date().getHours() + hour] }}
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
    </div>
  </div>


</template>

<style scoped>

.condition {
  width: 4rem;
  height: 4rem;
  align-items: center;
  justify-content: center;
  display: flex;
  font-size: 1.25rem;
}

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