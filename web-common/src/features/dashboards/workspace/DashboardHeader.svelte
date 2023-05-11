<script lang="ts">
  import { goto } from "$app/navigation";
  import { Button } from "@rilldata/web-common/components/button";
  import MetricsIcon from "@rilldata/web-common/components/icons/Metrics.svelte";
  import PanelCTA from "@rilldata/web-common/components/panel/PanelCTA.svelte";
  import Tooltip from "@rilldata/web-common/components/tooltip/Tooltip.svelte";
  import TooltipContent from "@rilldata/web-common/components/tooltip/TooltipContent.svelte";
  import { projectShareStore } from "@rilldata/web-common/features/dashboards/dashboard-stores";
  import { featureFlags } from "@rilldata/web-common/features/feature-flags";
  import { behaviourEvent } from "@rilldata/web-common/metrics/initMetrics";
  import { BehaviourEventMedium } from "@rilldata/web-common/metrics/service/BehaviourEventTypes";
  import {
    MetricsEventScreenName,
    MetricsEventSpace,
  } from "@rilldata/web-common/metrics/service/MetricsTypes";
  //  import { getContext } from "svelte";
  //  import type { Tweened } from "svelte/motion";
  import { runtime } from "../../../runtime-client/runtime-store";
  import Filters from "../filters/Filters.svelte";
  import { useMetaQuery } from "../selectors";
  import TimeControls from "../time-controls/TimeControls.svelte";

  export let metricViewName: string;
  export let hasTitle: boolean;

  //  const navigationVisibilityTween = getContext(
  //    "rill:app:navigation-visibility-tween"
  //  ) as Tweened<number>;

  const viewMetrics = (metricViewName: string) => {
    goto(`/dashboard/${metricViewName}/edit`);

    behaviourEvent.fireNavigationEvent(
      metricViewName,
      BehaviourEventMedium.Button,
      MetricsEventSpace.Workspace,
      MetricsEventScreenName.Dashboard,
      MetricsEventScreenName.MetricsDefinition
    );
  };

  $: metaQuery = useMetaQuery($runtime.instanceId, metricViewName);
  $: displayName = $metaQuery.data?.label;
  $: isEditableDashboard = $featureFlags.readOnly === false;

  function deployModal() {
    projectShareStore.set(true);
  }
</script>

<section class="w-full flex flex-col" id="header">
  <!-- top row: title and call to action -->
  <!-- Rill Local includes the title, Rill Cloud does not -->

  <!-- bottom row -->
  <div class="-ml-3 p-1 py-2 space-y-2">
    <TimeControls {metricViewName} />
    <!-- {#key metricViewName}
      <Filters {metricViewName} />
    {/key} -->
  </div>
</section>
