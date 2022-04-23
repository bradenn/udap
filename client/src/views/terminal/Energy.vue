<script lang="ts" setup>
import {reactive} from "vue";
import type {Allocation} from "@/components/AllocationBar.vue";
import AllocationBar from "@/components/AllocationBar.vue";

let s1: Allocation[] = [
  {
    name: 'MacMini',
    source: "sq2",
    allocated: 0.325,
    active: 0.21
  },
  {
    name: 'Nuc',
    source: "sq2",
    allocated: 1.5,
    active: 1.17
  }, {
    name: 'AP',
    source: "sq2",
    allocated: 0.15,
    active: 0.14
  },
]

let s2: Allocation[] = [

  {
    name: 'Macbook Pro',
    source: "sq1",
    allocated: 2,
    active: 1.56
  },
  {
    name: 'LG Monitor 1',
    source: "sq1",
    allocated: 1.8,
    active: 1.1
  },
  {
    name: 'LG Monitor 2',
    source: "sq1",
    allocated: 1.2,
    active: 0.8
  },
  {
    name: 'HP1',
    source: "sq1",
    allocated: 0.5,
    active: 0.45
  }, {
    name: 'HP2',
    source: "sq1",
    allocated: 0.5,
    active: 0.45
  }, {
    name: 'TV',
    source: "sq1",
    allocated: 1.21,
    active: 0.11
  },
  {
    name: 'AT',
    source: "sq1",
    allocated: 1,
    active: 0.12
  }
]

let ups: Allocation[] = [
  {
    name: 'Mac Mini',
    source: "sq2",
    allocated: 0.325,
    active: 0.22
  },
  {
    name: 'Nuc',
    source: "sq2",
    allocated: 1.5,
    active: 1.34
  }, {
    name: 'AP',
    source: "sq2",
    allocated: 0.15,
    active: 0.14
  },
]

let state = reactive({
  side1: s1,
  side2: s2,
  ups: ups,
})
//
// onMounted(() => {
//   GoogleCharts.charts.load('current', {'packages': ['sankey']});
//   GoogleCharts.charts.setOnLoadCallback(drawChart);
//   let colors = ['#a6cee3', '#b2df8a', '#fb9a99', '#fdbf6f',
//     '#cab2d6', '#ffff99', '#1f78b4', '#33a02c'];
//
//   function drawChart() {
//     let data = new GoogleCharts.visualization.DataTable();
//     data.addColumn('string', 'From');
//     data.addColumn('string', 'To');
//     data.addColumn('number', 'Weight');
//     data.addRows([
//       ['C1', 'A', 15],
//       ['C2', 'B', 15],
//       ['A', 'MAC', 2],
//       ['A', 'MON1', 7.37],
//       ['A', 'MON2', 1.2],
//       ['B', 'AP', 0.15],
//       ['B', 'HP1', 0.5],
//       ['B', 'HP2', 0.5],
//       ['B', 'TV', 1.21],
//       ['B', 'AT', 1],
//     ]);
//
//     // Sets chart options.
//     let options = {
//       sankey: {
//         node: {
//           colors: colors
//
//         },
//         link: {
//           colorMode: 'gradient',
//           colors: colors
//         }
//       },
//       width: 300,
//     };
//
//     // Instantiates and draws our chart, passing in some options.
//     let chart = new google.visualization.Sankey(document.getElementById('sankey_basic'));
//     chart.draw(data, options);
//
//     google.charts.load("current", {packages: ["corechart"]});
//     google.charts.setOnLoadCallback(drawChart);
//
//     function drawChart() {
//       let data = google.visualization.arrayToDataTable([
//         ['DEVICE', 'AMPS'],
//         ['MAC', 2],
//         ['MON1', 7.37],
//         ['MON2', 1.2],
//         ['AP', 0.15],
//         ['HP1', 0.5],
//         ['HP2', 0.5],
//         ['TV', 1.21],
//         ['AT', 1],
//       ]);
//
//       let options = {
//         width: 400,
//         backgroundColor: 'transparent',
//         legend: {
//           textStyle: {
//             color: "white"
//           }
//         },
//         is3D: false,
//       };
//
//       let chart = new google.visualization.PieChart(document.getElementById('piechart_3d'));
//       chart.draw(data, options);
//     }
//   }
// })
</script>

<template>

  <div class="d-flex flex-column p-2 element gap-2">
    <div>
      <h3 class="px-1 label-o5">Energy Allocations</h3>
    </div>
    <AllocationBar :allocatable="4" :allocations="state.ups" name="Battery Reserve"></AllocationBar>
    <AllocationBar :allocatable="15" :allocations="state.side1" name="Squid 1"></AllocationBar>
    <AllocationBar :allocatable="15" :allocations="state.side2" name="Squid 2"></AllocationBar>
  </div>
</template>

<style lang="scss" scoped>

</style>
