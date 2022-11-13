<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {onMounted, reactive} from "vue";
import axios from "axios";
import wtf from 'wtf_wikipedia'
// @ts-ignore
import wtf_html from 'wtf-plugin-html'
import FixedScroll from "@/components/scroll/FixedScroll.vue";

wtf.extend(wtf_html)

let state = reactive({
  data: {},
  dom: ""
})


onMounted(() => {
  fetch()
})

function fetch() {
  axios.get("https://en.wikipedia.org/w/api.php?action=query&format=xml&origin=*&export=1&titles=Cat").then(res => {
    state.data = wtf(res.data).json()
    // @ts-ignore
    state.dom = wtf(res.data).html()
  }).catch(err => {
    state.dom = err
  })
}

</script>

<template>
  <div class="h-100">
    <div class="element p-2 " style="height: 100%;">
      <FixedScroll class="" style=" height: 100%; overflow-y: scroll;">
        <div class=" " v-html="state.dom"></div>
      </FixedScroll>
    </div>
  </div>
</template>

<style>


td {
  font-size: 0.75rem !important;
  width: 2rem
}

td .sentence {
  font-size: 0.5rem !important;
}

.infobox {
  border-radius: 0.5rem;
  border: 1px solid #a2a9b1;
  border-spacing: 3px;
  margin: 0.5em 0.5rem 0.5em 1em;
  padding: 0.2em;
  float: right;
  clear: right;
  line-height: 1.5em;
  width: 16em;
}

.section img {
  float: right;
  clear: right;
}

.section h1 {
  font-size: 1.4rem;
  width: 100%;
}

.section h2 {
  font-size: 1rem;
  width: 100%;
}

.text .paragraph {
  font-size: 0.75rem !important;
}

.text .paragraph {
  font-size: 0.75rem !important;
}

.text .link {
  font-size: 0.75rem !important;
}


.infobox {

}
</style>