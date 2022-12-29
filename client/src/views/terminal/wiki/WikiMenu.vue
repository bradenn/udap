<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {onMounted, reactive} from "vue";

import wtf from 'wtf_wikipedia'
// @ts-ignore
import wtf_html from 'wtf-plugin-html'
import FixedScroll from "@/components/scroll/FixedScroll.vue";


const sections = ["Culture", "Geography", "Health", "History", "Human activities", "Mathematics", "Nature", "People", "Philosophy", "Religion", "Society", "Technology"]


wtf.extend(wtf_html)

let state = reactive({
  data: {},
  dom: ""
})


onMounted(() => {
  fetch()
})

function loadArticle(doc: Document) {

}

function fetch() {


  wtf.fetch("maser").then(doc => {
    // @ts-ignore
    state.dom = doc.html()
  })
  // axios.get("https://en.wikipedia.org/w/api.php?action=query&format=xml&origin=*&export=1&titles=Cat").then(res => {
  //     // @ts-ignore
  //     state.dom = wtf(res.data).html()
  // }).catch(err => {
  //     state.dom = err
  // })
}

</script>

<template>
  <div class="h-100">
    <div v-if="false" class="">
      <div class="hex-row">
        <div class="hex">
          <div class="top"></div>
          <div class="middle">Culture</div>
          <div class="bottom"></div>
        </div>
        <div class="hex">
          <div class="top"></div>
          <div class="middle">Geography</div>
          <div class="bottom"></div>
        </div>
        <div class="hex">
          <div class="top"></div>
          <div class="middle">Health</div>
          <div class="bottom"></div>
        </div>
        <div class="hex">
          <div class="top"></div>
          <div class="middle">History</div>
          <div class="bottom"></div>
        </div>
      </div>
      <div class="hex-row even">
        <div class="hex">
          <div class="top"></div>
          <div class="middle">Activities</div>
          <div class="bottom"></div>
        </div>
        <div class="hex">
          <div class="top"></div>
          <div class="middle">Mathematics</div>
          <div class="bottom"></div>
        </div>
        <div class="hex">
          <div class="top"></div>
          <div class="middle">Nature</div>
          <div class="bottom"></div>
        </div>
        <div class="hex">
          <div class="top"></div>
          <div class="middle">People</div>
          <div class="bottom"></div>
        </div>
      </div>
    </div>
    <div v-if="false" class="menu-grid">
      <div v-for="s in sections"
           class="element d-flex justify-content-center label-c1 align-items-center">
        <div class="label-c1 label-w600 py-2">{{ s }}</div>
      </div>
    </div>
    <div v-if=true class="element p-2 " style="height: 100%;">
      <FixedScroll class="" style=" height: 100%; overflow-y: scroll;">
        <div class=" " v-html="state.dom"></div>
      </FixedScroll>
    </div>
  </div>
</template>

<style lang="scss">

$width: 5.5rem;
$width2: $width/2;
$width4: $width/4;
$height: $width*0.5769230769;
$height2: $height/2;

.hex {
  float: left;
  margin-left: 4px;
  margin-bottom: -$width4 - 0.125rem;

}

.hex .top {
  width: 0;

  border-bottom: $height2 solid rgba(255, 255, 255, 0.1);
  border-left: $width2 solid transparent;
  border-right: $width2 solid transparent;
}

.hex .middle {
  width: $width;
  height: $height;
  background: rgba(255, 255, 255, 0.1);
  font-size: 0.7rem;
  display: flex;
  font-weight: 500;
  justify-content: center;

  align-items: center;
}

.hex .bottom {
  width: 0;
  border-top: $height2 solid rgba(255, 255, 255, 0.1);
  border-left: $width2 solid transparent;
  border-right: $width2 solid transparent;
}

.hex-row {
  clear: left;
}

.hex-row.even {
  margin-left: calc($width / 2) + 0.05rem;
}

div.menu-grid {
  display: grid;
  grid-column-gap: 0.5rem;
  grid-row-gap: 0.25rem;
  grid-template-rows: repeat(8, 1fr);
  grid-template-columns: repeat(8, 1fr);
}

td {
  font-size: 0.75rem !important;
  width: 2rem
}

td .sentence {
  font-size: 0.5rem !important;
}

.infobox {
  border-radius: 0.525rem;
  background-color: rgba(255, 255, 255, 0.05);
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