<script lang="ts" setup>
import {reactive} from "vue";

let state = reactive({
  chartData: [
    ['T', 'X', 20],
    ['A', 'X', 1],
    ['A', 'Y', 4],
    ['A', 'Z', 8],
    ['B', 'C', 2],
    ['B', 'D', 1],
    ['B', 'Z', 8],
  ]
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
  <div class="d-flex flex-row gap justify-content-start align-items-start align-content-center h-100 pt-3">
    <div class="element p-3 pt-2 w-100">
      <h4>Current Load</h4>
      <div></div>
    </div>
    <div class="element p-3 pt-2">
      <h4>Dynamic Allocations</h4>
      <div id="piechart_3d" style="height:350px; width: 300px;"></div>
    </div>
    <div class="element p-3 pt-2">
      <h4>Dynamic Allocations</h4>
      <div id="sankey_basic" style="height:350px; width: 300px;"></div>
    </div>

  </div>
</template>

<style scoped>

</style>
