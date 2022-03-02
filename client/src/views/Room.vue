<script>

import Hexagon from "../components/Hexagon.vue";
import Entity from "../components/entity/Entity.vue";
import Range from "../components/entity/Range.vue";

export default {
  components: {Range, Entity, Hexagon},
  data() {
    return {
      room: [[]],
      selected: Array(),
      attrMain: {key: "dim", value: '0', request: '50', type: 'range', id: 'dsddssd', entity: 'dssdsd'},
      grid: {
        width: 44,
        height: 60,
        radius: 8,
      },
      global: [
        {key: "dim", value: '0', request: '0', type: 'range', id: '', entity: ''},
        {key: "cct", value: '0', request: '0', type: 'range', id: '', entity: ''},
        {key: "hue", value: '0', request: '0', type: 'range', id: '', entity: ''}
      ],
      moods: [
        {
          name: "Dim",
          icon: "",
          attributes: {
            dim: 25,
            cct: 4000
          }
        },
      ],
      cooldown: false,
    }
  },
  created() {
    let grid = Array(this.grid.width)

    for (let y = 0; y < this.grid.width; y++) {

      grid[y] = Array(this.grid.height)

      for (let x = 0; x < this.grid.height; x++) {

        grid[y][x] = 55

      }

    }

    this.room = grid
  },
  computed: {
    live: function () {
      let res = []
      this.room.forEach(r => res.push(r))
      return res
    },
    manifestSelected: function () {
      return this.$root.entities.filter(e => this.selected.includes(e.id)).sort((a, b) => a.type - b.type)
    },
    shownGlobals: function () {
      let shown = []
      for (const {v, k} in this.groupBy(this.manifestSelected, "type")) {

      }
    },
    manifest: function () {
      function compare(a, b) {
        if (a.created < b.created)
          return -1;
        if (a.created > b.created)
          return 1;
        return 0;
      }

      return this.$root.entities.filter(a => {
        if (a.type)
          return a.type === "spectrum"
      }).map(e => {
        e.selected = this.selected.includes(e.id)
        let pos = {};

        return e
      }).sort(compare)
    }
  },
  methods: {
    set(x, y, v) {
      if (x < 0 || x > this.grid.height || y < 0 || y > this.grid.width) return
      this.room[Math.min(this.grid.width, Math.max(0, y))][Math.min(this.grid.height, Math.max(0, x))] += Math.min(255, Math.max(0, v))
    },

    get(x, y) {
      return this.room[x][y]
    },
    reset() {
      this.cooldown = false

    },
    gridClick(y, x) {


      /*    for(let a = Math.max(0, y - rad); a < Math.min(44, y + rad); a++){

          }*/
    },
    go(link) {
      this.$router.replace({path: link})
    },
    isSelected(id) {
      return function () {
        this.selected.includes(id)
      }
    },
    addSelect(id) {
      return function (e) {
        alert("YES")
        /* */
      }
    },
    selectOne: function (id) {
      if (this.selected.includes(id)) {
        this.selected = this.selected.filter(s => s !== id)
      } else {
        this.selected.push(id)
      }
    },
    selectAll: function () {
      for (let manifestElement of this.manifest) {
        this.selectOne(manifestElement.id)
      }
    },
    toggleAll: function (state) {
      let attr = {
        key: "on",
        request: `${state}`,
        value: ''
      }
      this.commitChangeAll(attr)
    },
    dimAll: function (value) {
      let attr = {
        key: "dim",
        request: `${value}`,
        value: '10'
      }
      this.commitChangeAll(attr)
    },
    commitChangeAll(attribute) {
      let filtered = this.$root.attributes.filter(a => a.key === attribute.key && this.selected.includes(a.entity))
      filtered.map(a => {
        a.request = attribute.request;
        this.$root.requestId("attribute", "request", a, a.id)
        return a
      })
    },
    groupBy(xs, key) {
      return xs.reduce(function (rv, x) {
        (rv[x[key]] = rv[x[key]] || []).push(x);
        return rv;
      }, {});
    },
  },
}


</script>

<template>
  <div class="">

    <div class="grid-pane">
      <div class="pane-room element">
        <div class="room-container">
          <div class="room label-o2 m-1">
            <div
                v-for="light in manifest.filter(m => m.type === 'spectrum')"
                :key="light.id"
                :style="`position: absolute; top: ${JSON.parse(light.position).y}%; left: ${JSON.parse(light.position).x}%; height: 1rem; width: 100%;`"
                @click="() => selectOne(light.id)">
              <Entity :entity="light" :selected="light.selected" bulb></Entity>

            </div>
          </div>
        </div>
      </div>

      <div class="pane-entities d-flex flex-column gap top mx-2">
        <div class="element d-flex gap align-items-center px-2">
          <div class="label-xs label-w500 label-o4 px-1" v-on:click="toggleAll(true)">ON</div>
          <div class="label-xs label-w500 label-o4 px-1" v-on:click="toggleAll(false)">OFF</div>
          <div class="mx-2 my-1"
               style="width: 0.0255rem; height: 1rem; border-radius: 1rem; background-color: rgba(255,255,255,0.1);"></div>
          <div class="label-xs label-w500 label-o4 px-1" v-on:click="dimAll(25)">25%</div>
          <div class="label-xs label-w500 label-o4 px-1" v-on:click="dimAll(50)">50%</div>
          <div class="label-xs label-w500 label-o4 px-1" v-on:click="dimAll(50)">75%</div>
          <div class="mx-2 my-1"
               style="width: 0.0255rem; height: 1rem; border-radius: 1rem; background-color: rgba(255,255,255,0.1);"></div>
          <div class="label-xs label-w500 label-o4 px-3" v-on:click="selectAll()">Select All</div>
        </div>

        <div class="d-flex flex-column gap">
          <Range v-for="attribute in global" :key="attribute.id" :attribute="attribute" :commit="commitChangeAll"
          ></Range>
        </div>

        <div v-for="(entities, type) in groupBy(manifest, 'type')" :key="type">
          <div class="label-xxs label-o6 mb-1 label-w500 text-capitalize" style="width: 5rem">{{ type }}</div>
          <div class="cluster gap">
            <div v-for="entity in entities" :key="entity.type" @click="() => selectOne(entity.id)">
              <Entity :entity="entity" :selected="entity.selected" manage></Entity>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>


<style scoped>
.grid-pane {
  margin-top: 1rem;
  display: flex;
  justify-content: center;
  align-content: center;
  align-items: center;
  height: 100%;
  min-width: 80vw !important;
}

.pane-room {
  grid-column: 2;
}

.pane-entities {

}

.room-container {
  filter: drop-shadow(0 0 4px rgba(0, 0, 0, 0.24));
  overflow: auto;
}

.matrix {
  padding: 0;
  display: grid;
  grid-template-columns: repeat(44, 0.25rem);
  grid-template-rows: repeat(60, 0.25rem);
  gap: 0;
  overflow: hidden !important;
}

.room {
  aspect-ratio: 0.7352973 !important;
  height: 20rem;
  clip-path: polygon(0% 100%, 100% 100%, 100% 0%, 65% 0%, 65% 8%, 0% 52%);
  background-color: rgba(255, 255, 255, 0.0625);
}

.unit-inner {
  transform: scale(0.2);
}

.unit {

  position: relative;
  line-height: 0.25rem;
  font-size: 100%;
  width: 0.25rem;
  height: 0.25rem;
  aspect-ratio: 1/1 !important;
  display: flex;
  justify-content: center;
  align-items: center;
  align-content: center;
  color: rgba(255, 255, 255, 0.1);
}


.l1 {
  position: absolute;
  bottom: calc(50% - 1rem);
  left: calc(50% - 1rem);
  font-size: 1rem;
  line-height: 2rem;
  width: 2rem;
  height: 2rem;

  display: flex;
  align-items: center;
  align-content: center;
  justify-content: center;
}

.l1:before {
  position: absolute;
  bottom: calc(50% - 1rem);
  left: calc(50% - 1rem);
  content: ' ';
  border-radius: 100%;
  width: 2rem;
  height: 2rem;
  background-color: rgba(255, 255, 255, 1);
  filter: blur(16px) !important;
}

.subroutine {
  width: 10rem;
  display: flex;
  justify-content: start;
  flex-direction: row;
}
</style>
