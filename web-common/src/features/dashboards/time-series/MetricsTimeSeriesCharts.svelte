<script lang="ts">
  import {
    metricsExplorerStore,
    useDashboardStore,
  } from "@rilldata/web-common/features/dashboards/dashboard-stores";
  import TestSlot from "./TestSlot.svelte";

  $: dashboardStore = useDashboardStore("ad_bids_dashboard");

  $: testDC = $dashboardStore?.showComparison;
</script>

<div>Outside of slots, showComparison is: {testDC}</div>
<TestSlot>
  Inside of slots, showComparison is: {testDC}
  <div>metrics store {$metricsExplorerStore}</div>

  <!-- <div>dash read {$dashboardStore.showComparison}</div> -->
</TestSlot>
