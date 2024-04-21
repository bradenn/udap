<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import List from "udap-ui/components/List.vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";
import {since} from "udap-ui/time";
import {onMounted, reactive} from "vue";
import {Detour} from "udap-ui/types";
import DoseClock from "@/components/DoseClock.vue";

const state = reactive({
  dose: 5,
  detours: [] as Detour[],
  time: 0,
  since: "",
  timeString: "",
})


interface DetourData {
  detours: Detour[]

}

onMounted(() => {
  update()

})

function updateSince() {
  state.since = since(`${state.time}`, true)
}

function update() {
  let detour: string = localStorage.getItem("detours") as string
  if (!detour) {
    let a: DetourData = {detours: []} as DetourData
    localStorage.setItem("detours", JSON.stringify(a));
    detour = localStorage.getItem("detours") as string
  }
  let ps = JSON.parse(detour) as DetourData
  if (!ps) {
    return
  }
  state.detours = ps.detours
  state.timeString = new Date().toLocaleTimeString('it-IT', {hour: '2-digit', minute: '2-digit'})
  state.time = getDateFromTime(state.timeString)
  updateSince()
}

function addDose() {
  if (!state.detours) {
    state.detours = []
  }
  state.detours.push({
    intensity: 0,
    time: getDateFromTime(state.timeString),
    dose: state.dose
  })
  saveDetours()
}

function getDateFromTime(time: string): number {

  let now = new Date()
  now.setHours(parseInt(time.split(":")[0]), parseInt(time.split(":")[1]), 0, 0)
  return now.valueOf()

}


function saveDetours() {
  localStorage.setItem("detours", JSON.stringify({detours: state.detours} as DetourData));
}

function change(i: InputEvent) {
  // state.time = getDateFromTime(state.timeString)
  // updateSince()
}

function clearDoses() {
  state.detours = []
  saveDetours()
}

function timeSince(time: number): string {
  let delta = Date.now() - time

  let minutes = Math.floor(delta / 1000 / 60 % 60)
  let hours = Math.floor(delta / 1000 / 60 / 60)

  return `${hours} hour ${minutes} min ago`
}

function draw() {

}

function calculatePotency(dose: string, since: number): number {
  let delta = Date.now().valueOf() - since

  let x = 1 / (delta / (1000 * 60 * 60 * 3))


  return (Math.sin((x - 2) * Math.PI / 2) / 2)
}

</script>

<template>
  <List>
    <Element>
      <DoseClock :detours="state.detours" style="height: 12rem"></DoseClock>
    </Element>

    <Element>
      <ElementHeader title="Recent Detours"></ElementHeader>
      <List>
        <Element v-for="detour in state.detours" v-if="state.detours.length > 0" :key="detour.time" foreground mutable>
          <div class="d-flex justify-content-between align-items-center">
            <div class="d-flex align-items-start justify-content-center flex-column lh-1 gap-1">
              <div class="label-c3 label-w500 label-o4 mono">
                {{ new Date(detour.time).toLocaleTimeString() }}
              </div>
              <div class="label-c5 label-w500 label-o3 mono">
                {{ timeSince(detour.time) }}
              </div>

            </div>
            <div class="label-c0 label-w600 label-r px-3 align-items-start d-flex align-items-baseline gap-1 ">
              <div>
                {{ detour.dose }}
              </div>
              <div class="label-o4" style="font-size: 14px">mg</div>
            </div>
          </div>

        </Element>
        <Element v-else class="button py-2 px-3" foreground>No
          logged
          detours
        </Element>
      </List>
    </Element>
    <Element>

      <ElementHeader title="Add"></ElementHeader>
      <List>
        <List row>
          <Element :accent="state.dose==10" :cb="() => state.dose=10" class="button mono" foreground mutable>+ 10 mg
          </Element>
          <Element :accent="state.dose==5" :cb="() => state.dose=5" class="button mono" foreground mutable>+ 5 mg
          </Element>
          <Element :accent="state.dose==2.5" :cb="() => state.dose=2.5" class="button mono" foreground mutable>+ 2.5 mg
          </Element>

        </List>
        <Element foreground>
          <div class="d-flex align-items-center justify-content-between">
            <div class="label-c4 label-w500 label-o4 px-2">Ingest time</div>
            <div><input id="current-time" v-model="state.timeString" style="" type="time" @input="(a) => change"/></div>
          </div>
        </Element>

        <Element :cb="() => addDose()" class="mono button label-w600"
                 foreground>
          Record {{ state.dose }}mg @ {{ new Date(getDateFromTime(state.timeString)).toLocaleTimeString() }}
        </Element>
        <Element :cb="() => clearDoses()" class="button mono label-o3" foreground>
          Clear History
        </Element>
      </List>

    </Element>
  </List>

</template>

<style lang="scss" scoped>
input:focus-visible, input:active, input:hover {
  outline: none;
  border: none;

  //border: none;
}

input[type=time] {
  outline: none;
  background-color: transparent;

  min-width: 7rem;
  padding-inline: 0.5rem;
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 4px;
  color: rgba(255, 255, 255, 0.7);
  height: 2.4rem;
}

.button.mono {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.9rem;
  letter-spacing: -0.5px;
  min-height: 3.2rem;


}
</style>