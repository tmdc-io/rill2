<script lang="ts">
  // @ts-nocheck

  import picasso from "picasso.js";
  import {
    createQueryServiceMetricsViewAggregation
  } from "@rilldata/web-common/runtime-client";
  import { getStateManagers } from "@rilldata/web-common/features/dashboards/state-managers/state-managers";

  const {
    selectors: {
      dimensions: {
        dimensionTableDimName
      },
      timeRangeSelectors: {
        timeControlsState
      },
    },
    dashboardStore,
    metricsViewName,
    runtime,
    queryClient,
  } = getStateManagers();

  let chartElement;
  let p;

  $: data = createQueryServiceMetricsViewAggregation(
    $runtime.instanceId,
    $metricsViewName,
    {
      measures: [{ name: $dashboardStore?.selectedMeasureNames[0] }],
      dimensions: [{ name: $dimensionTableDimName }],
      filter: $dashboardStore.filters,
      timeStart: $timeControlsState.adjustedStart,
      timeEnd: $timeControlsState.adjustedEnd,
      limit: "250"
    },
    {
      query: {
        enabled: true,
        queryClient
      }
    }
  )

  let headers;
  $: headers = $data.data?.schema?.fields?.map(d => d.name)
  
  let rows: any;
  $: rows = $data.data?.data?.map(d => Object.values(d))
  $: if ($data.isSuccess && rows?.length && headers?.length) {
    if (p) {
      p.update({
        settings: getSettings(),
        data: {
          type: "matrix",
          data: [headers, ...rows]
        }
      })
    } else {
      renderChart([headers, ...rows])
    }
  }
  const getSettings = () => ({
    scales: {
      y: {
        data: { field: headers[1] },
        invert: true,
        include: [0],
      },
      t: { data: { extract: { field: headers[0] } }, padding: 0.3 },
    },
    components: [
      {
        type: "axis",
        dock: "left",
        scale: "y",
      },
      {
        type: "axis",
        dock: "bottom",
        scale: "t",
      },
      {
        key: "bars",
        type: "box",
        animations: {
          enabled: true,
          trackBy: (node) => node.data.value,
        },
        data: {
          extract: {
            field: 0,
            props: {
              start: 0,
              end: { field: headers[1] },
            },
          },
        },
        settings: {
          major: { scale: "t" },
          minor: { scale: "y" },
          box: {
            fill: "rgb(219 234 254)",
          },
        },
      },
    ],
  });

  const renderChart = (data) => {
    p = picasso.chart({
      element: chartElement,
      data: [{
        type: "matrix",
        data
      }],
      settings: getSettings(),
    });
  };

</script>

<div bind:this={chartElement} class="chart h-full flex flex-col" />

<style>
  .chart {
    height: 200px;
    width: 100%;
    position: relative;
  }
</style>
