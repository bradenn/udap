<!-- Copyright (c) 2024 Braden Nicholson -->

<script lang="ts" setup>
import Element from "udap-ui/components/Element.vue";
import List from "udap-ui/components/List.vue";
import ElementDecimal from "udap-ui/components/ElementDecimal.vue";
import {reactive} from "vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";


const state = reactive({
  radius: 0,
  mass: 0,
  initialVelocity: 0,
  finalVelocity: 0,
  elapsedTime: 0,
})


const shapes = ["cylinder"]

</script>

<template>
  <Element>
    <List>
      <ElementHeader title="Inertia"></ElementHeader>
      <ElementDecimal :change="(v) => state.radius = v" :value="state.radius" icon="􀰬" mutable title="Radius"
                      unit="mm"></ElementDecimal>
      <ElementDecimal :change="(v) => state.mass = v" :value="state.mass" icon="􀭮" mutable title="Mass"
                      unit="g"></ElementDecimal>
      <ElementHeader title="Angular Acceleration"></ElementHeader>
      <ElementDecimal :change="(v) => state.initialVelocity = v" :value="state.initialVelocity" icon="􀐳" mutable
                      title="Initial Velocity"
                      unit="rpm"></ElementDecimal>
      <ElementDecimal :change="(v) => state.finalVelocity = v" :value="state.finalVelocity" icon="􀐳" mutable
                      title="Final Velocity"
                      unit="rpm"></ElementDecimal>
      <ElementDecimal :change="(v) => state.elapsedTime = v" :value="state.elapsedTime" icon="􀣔" mutable
                      title="Elapsed Time"
                      unit="s"></ElementDecimal>
      <ElementHeader title="Torque"></ElementHeader>
      <ElementDecimal
          :value="(Math.pow(state.radius/1000, 2) * (state.mass/1000) * ((state.finalVelocity * (Math.PI*2)/60) - (state.initialVelocity * (Math.PI*2)/60)) / (state.elapsedTime))*980.665"
          icon="􀎖" title="Torque"
          unit="N•cm"
      ></ElementDecimal>
    </List>

  </Element>
</template>

<style lang="scss" scoped>

</style>